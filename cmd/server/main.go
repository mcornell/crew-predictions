package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	firebase "firebase.google.com/go/v4"
	"github.com/joho/godotenv"
	"github.com/mcornell/crew-predictions/internal/espn"
	"github.com/mcornell/crew-predictions/internal/handlers"
	"github.com/mcornell/crew-predictions/internal/models"
	"github.com/mcornell/crew-predictions/internal/repository"
	"google.golang.org/api/option"
)

func main() {
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	ctx := context.Background()

	var predStore repository.PredictionStore
	if project := os.Getenv("GOOGLE_CLOUD_PROJECT"); project != "" {
		fs, err := repository.NewFirestorePredictionStore(ctx, project)
		if err != nil {
			log.Fatalf("failed to connect to Firestore: %v", err)
		}
		predStore = fs
		log.Printf("using Firestore (project: %s)", project)
	} else {
		predStore = repository.NewMemoryPredictionStore()
		log.Printf("using in-memory store (set GOOGLE_CLOUD_PROJECT to use Firestore)")
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
		log.Printf("Firebase unavailable: %v", err)
	} else if authClient, err := app.Auth(ctx); err != nil {
		log.Printf("Firebase Auth unavailable: %v", err)
	} else {
		verifier = handlers.NewFirebaseTokenVerifier(authClient)
		log.Printf("Firebase Auth initialized (project: %s)", projectID)
	}
	sh := handlers.NewSessionHandler(verifier)

	mux := http.NewServeMux()

	// Match store: source of truth for cached matches
	matchStore := repository.NewMemoryMatchStore()

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
	mux.HandleFunc("GET /api/me", handlers.Me)
	ph := handlers.NewPredictionsHandler(predStore, matchFetcher)
	mux.HandleFunc("POST /api/predictions", ph.Submit)
	lh := handlers.NewLeaderboardHandler(predStore, resultStore, "Columbus Crew")
	mux.HandleFunc("GET /api/leaderboard", lh.APIList)
	rh := handlers.NewResultsHandler(resultStore)
	mux.HandleFunc("POST /admin/results", rh.Submit)
	rmh := handlers.NewRefreshMatchesHandler(matchStore, refreshFetcher)
	mux.HandleFunc("POST /admin/refresh-matches", rmh.Refresh)

	// Auth endpoints
	mux.HandleFunc("POST /auth/session", sh.Create)
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
		seedMH := handlers.NewSeedMatchHandler(matchStore)
		mux.HandleFunc("POST /admin/seed-match", seedMH.Submit)
		log.Printf("test seed endpoint registered at POST /admin/seed-match")
	}

	log.Printf("listening on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
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
