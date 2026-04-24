package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	firebase "firebase.google.com/go/v4"
	"github.com/joho/godotenv"
	"github.com/mcornell/crew-predictions/internal/bot"
	"github.com/mcornell/crew-predictions/internal/espn"
	"github.com/mcornell/crew-predictions/internal/handlers"
	"github.com/mcornell/crew-predictions/internal/models"
	"github.com/mcornell/crew-predictions/internal/poll"
	"github.com/mcornell/crew-predictions/internal/recalculator"
	"github.com/mcornell/crew-predictions/internal/repository"
	"google.golang.org/api/option"
)

func main() {
	godotenv.Load()

	// Cloud Logging parses JSON from stdout and maps "severity"/"message" fields.
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == slog.LevelKey {
				a.Key = "severity"
				if a.Value.String() == "WARN" {
					a.Value = slog.StringValue("WARNING")
				}
			}
			if a.Key == slog.MessageKey {
				a.Key = "message"
			}
			return a
		},
	})))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	ctx := context.Background()

	var userStore repository.UserStore
	if project := os.Getenv("GOOGLE_CLOUD_PROJECT"); project != "" {
		fs, err := repository.NewFirestoreUserStore(ctx, project)
		if err != nil {
			log.Fatalf("failed to connect to Firestore users: %v", err)
		}
		userStore = fs
	} else {
		userStore = repository.NewMemoryUserStore()
	}

	var predStore repository.PredictionStore
	if project := os.Getenv("GOOGLE_CLOUD_PROJECT"); project != "" {
		fs, err := repository.NewFirestorePredictionStore(ctx, project)
		if err != nil {
			log.Fatalf("failed to connect to Firestore: %v", err)
		}
		predStore = fs
		slog.Info("store: Firestore", "project", project)
	} else {
		predStore = repository.NewMemoryPredictionStore()
		slog.Info("store: in-memory (set GOOGLE_CLOUD_PROJECT to use Firestore)")
	}

	var resultStore repository.ResultStore
	if project := os.Getenv("GOOGLE_CLOUD_PROJECT"); project != "" {
		fs, err := repository.NewFirestoreResultStore(ctx, project)
		if err != nil {
			log.Fatalf("failed to connect to Firestore results: %v", err)
		}
		resultStore = fs
	} else {
		resultStore = repository.NewMemoryResultStore()
	}

	var verifier handlers.TokenVerifier = handlers.NoopTokenVerifier{}
	projectID := os.Getenv("FIREBASE_PROJECT_ID")
	if projectID == "" {
		projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
	}
	var appOpts []option.ClientOption
	if os.Getenv("FIREBASE_AUTH_EMULATOR_HOST") != "" {
		appOpts = append(appOpts, option.WithoutAuthentication())
	}
	fbConfig := &firebase.Config{ProjectID: projectID}
	app, err := firebase.NewApp(ctx, fbConfig, appOpts...)
	if err != nil {
		slog.Error("Firebase unavailable", "error", err)
	} else if authClient, err := app.Auth(ctx); err != nil {
		slog.Error("Firebase Auth unavailable", "error", err)
	} else {
		verifier = handlers.NewFirebaseTokenVerifier(authClient)
		slog.Info("Firebase Auth initialized", "project", projectID)
	}
	sh := handlers.NewSessionHandler(verifier, userStore)
	hh := handlers.NewHandleHandler(userStore)

	mux := http.NewServeMux()

	// Match store: memory is primary (fast reads); Firestore secondary persists across restarts
	memMatchStore := repository.NewMemoryMatchStore()
	var matchStore repository.MatchStore = memMatchStore
	if project := os.Getenv("GOOGLE_CLOUD_PROJECT"); project != "" {
		fsMatches, err := repository.NewFirestoreMatchStore(project)
		if err != nil {
			log.Fatalf("failed to connect to Firestore matches: %v", err)
		}
		matchStore = repository.NewWriteThroughMatchStore(memMatchStore, fsMatches)
		// Pre-populate memory from Firestore so match data survives restarts
		if stored, err := fsMatches.GetAll(); err != nil {
			slog.Warn("could not load matches from Firestore", "error", err)
		} else if len(stored) > 0 {
			memMatchStore.SaveAll(stored)
			slog.Info("startup: matches loaded from Firestore", "count", len(stored))
		}
	}

	// ESPN fetcher used by the refresh endpoint; in TEST_MODE reads from seeded store
	var refreshFetcher func() ([]models.Match, error)
	if os.Getenv("TEST_MODE") == "1" {
		refreshFetcher = matchStore.GetAll
	} else {
		refreshFetcher = espn.FetchCrewMatches
	}

	// Fetcher for PredictionsHandler kickoff validation — reads from the same store
	matchFetcher := func() ([]models.Match, error) { return matchStore.GetAll() }

	// API endpoints (JSON)
	mh := handlers.NewMatchesHandler(predStore, matchStore)
	mux.HandleFunc("GET /api/matches", mh.APIList)
	meh := handlers.NewMeHandler(userStore)
	mux.HandleFunc("GET /api/me", meh.Get)
	ph := handlers.NewPredictionsHandler(predStore, matchFetcher)
	mux.HandleFunc("POST /api/predictions", ph.Submit)
	lh := handlers.NewLeaderboardHandler(predStore, resultStore, userStore, "Columbus Crew")
	mux.HandleFunc("GET /api/leaderboard", lh.APIList)
	prh := handlers.NewProfileHandler(predStore, resultStore, userStore, "Columbus Crew")
	mux.HandleFunc("GET /api/profile/{userID}", prh.Get)
	mdh := handlers.NewMatchDetailHandler(predStore, resultStore, matchStore, userStore, "Columbus Crew")
	mux.HandleFunc("GET /api/matches/{matchId}", mdh.Get)
	if os.Getenv("TEST_MODE") != "1" && os.Getenv("ADMIN_KEY") == "" {
		log.Fatal("ADMIN_KEY env var must be set in production")
	}

	rh := handlers.NewResultsHandler(resultStore)
	mux.HandleFunc("POST /admin/results", handlers.AdminAuth(rh.Submit))
	var matchPoller *poll.MatchPoller
	if os.Getenv("TEST_MODE") != "1" {
		matchPoller = poll.NewMatchPoller(
			matchStore, resultStore, refreshFetcher,
			func(d time.Duration, f func()) { time.AfterFunc(d, f) },
		)
	}

	twoOneBot := bot.New(predStore, userStore, "Columbus Crew")
	onRefresh := func(matches []models.Match) {
		if matchPoller != nil {
			matchPoller.Reset(matches)
		}
		twoOneBot.Predict(context.Background(), matches)
	}
	rmh := handlers.NewRefreshMatchesHandler(matchStore, refreshFetcher, onRefresh)
	mux.HandleFunc("POST /admin/refresh-matches", handlers.AdminAuth(rmh.Refresh))
	bfh := handlers.NewBackfillUsersHandler(predStore, userStore)
	mux.HandleFunc("POST /admin/backfill-users", handlers.AdminAuth(bfh.Backfill))
	psh := handlers.NewPollScoresHandler(matchStore, resultStore, refreshFetcher)
	mux.HandleFunc("POST /admin/poll-scores", handlers.AdminAuth(psh.Poll))

	// Auth endpoints
	mux.HandleFunc("POST /auth/session", sh.Create)
	mux.HandleFunc("POST /auth/handle", hh.Update)
	mux.HandleFunc("GET /auth/logout", handlers.Logout)
	mux.HandleFunc("GET /auth/config.js", serveFirebaseConfig)

	// Vite build assets
	mux.Handle("GET /assets/", assetsHandler("dist"))

	// SPA shell — all other routes serve index.html
	mux.Handle("GET /", spaHandler("dist"))

	if os.Getenv("TEST_MODE") == "1" {
		if memPred, ok := predStore.(*repository.MemoryPredictionStore); ok {
			mux.HandleFunc("DELETE /admin/reset", func(w http.ResponseWriter, r *http.Request) {
				memPred.Reset()
				if memResult, ok := resultStore.(*repository.MemoryResultStore); ok {
					memResult.Reset()
				}
				matchStore.Reset()
				w.WriteHeader(http.StatusNoContent)
			})
			log.Printf("test reset endpoint registered at DELETE /admin/reset")
			seedH := handlers.NewSeedPredictionHandler(predStore)
			mux.HandleFunc("POST /admin/seed-prediction", seedH.Submit)
			log.Printf("test seed endpoint registered at POST /admin/seed-prediction")
		}
		seedMH := handlers.NewSeedMatchHandler(memMatchStore)
		mux.HandleFunc("POST /admin/seed-match", seedMH.Submit)
		log.Printf("test seed endpoint registered at POST /admin/seed-match")
	}

	if os.Getenv("TEST_MODE") != "1" {
		etLoc, err := time.LoadLocation("America/New_York")
		if err != nil {
			log.Fatalf("failed to load ET timezone: %v", err)
		}
		recalcFn := func(ctx context.Context) {
			if err := recalculator.Recalculate(ctx, predStore, resultStore, userStore, "Columbus Crew"); err != nil {
				slog.Error("recalculate failed", "error", err)
			}
		}
		matchPoller.SetOnResultSaved(recalcFn)
		stop := startDailyRefresh(matchStore, refreshFetcher, matchPoller, twoOneBot, recalcFn, etLoc)
		defer close(stop)

		pollerCtx, cancelPoller := context.WithCancel(context.Background())
		defer cancelPoller()
		go matchPoller.Run(pollerCtx, 2*time.Minute)
	}

	slog.Info("server listening", "port", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}

func next4amET(now time.Time, loc *time.Location) time.Time {
	t := now.In(loc)
	candidate := time.Date(t.Year(), t.Month(), t.Day(), 4, 0, 0, 0, loc)
	if !candidate.After(now) {
		candidate = candidate.Add(24 * time.Hour)
	}
	return candidate
}

func startDailyRefresh(store repository.MatchStore, fetcher func() ([]models.Match, error), poller *poll.MatchPoller, twoOneBot *bot.TwoOneBot, recalcFn func(context.Context), etLoc *time.Location) chan struct{} {
	stop := make(chan struct{})
	go func() {
		refresh := func() {
			matches, err := fetcher()
			if err != nil {
				slog.Error("daily match refresh: ESPN fetch failed", "error", err)
				return
			}
			if err := store.SaveAll(matches); err != nil {
				slog.Error("daily match refresh: store save failed", "error", err)
				return
			}
			// Backfill results for matches that finished while no poller was active
			// (e.g., after a Cloud Run recycle, or matches completed before this deploy).
			poller.Backfill(context.Background(), matches)
			recalcFn(context.Background())
			twoOneBot.Predict(context.Background(), matches)
			slog.Info("daily match refresh complete", "matchCount", len(matches))
			poller.Reset(matches)
		}
		refresh()
		for {
			delay := time.Until(next4amET(time.Now(), etLoc))
			timer := time.NewTimer(delay)
			select {
			case <-timer.C:
				refresh()
			case <-stop:
				timer.Stop()
				return
			}
		}
	}()
	return stop
}

func serveFirebaseConfig(w http.ResponseWriter, r *http.Request) {
	cfg, _ := json.Marshal(map[string]string{
		"apiKey":           os.Getenv("FIREBASE_API_KEY"),
		"authDomain":       os.Getenv("FIREBASE_AUTH_DOMAIN"),
		"projectId":        os.Getenv("FIREBASE_PROJECT_ID"),
		"authEmulatorHost": os.Getenv("FIREBASE_AUTH_EMULATOR_HOST"),
	})
	w.Header().Set("Content-Type", "application/javascript")
	fmt.Fprintf(w, "window.__firebaseConfig = %s;", cfg)
}

func spaHandler(distDir string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		http.ServeFile(w, r, distDir+"/index.html")
	})
}

func assetsHandler(distDir string) http.Handler {
	fs := http.FileServer(http.Dir(distDir))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		fs.ServeHTTP(w, r)
	})
}
