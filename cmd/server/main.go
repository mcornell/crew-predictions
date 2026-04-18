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
	"github.com/mcornell/crew-predictions/internal/repository"
	"github.com/mcornell/crew-predictions/templates"
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

	resultStore := repository.NewMemoryResultStore()

	var verifier handlers.TokenVerifier = handlers.NoopTokenVerifier{}
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		log.Printf("Firebase unavailable (no credentials): %v", err)
	} else if authClient, err := app.Auth(ctx); err != nil {
		log.Printf("Firebase Auth unavailable: %v", err)
	} else {
		verifier = handlers.NewFirebaseTokenVerifier(authClient)
		log.Printf("Firebase Auth initialized")
	}
	sh := handlers.NewSessionHandler(verifier)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/matches", http.StatusFound)
	})
	mh := handlers.NewMatchesHandler(predStore, espn.FetchCrewMatches)
	mux.HandleFunc("GET /matches", mh.List)
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	ph := handlers.NewPredictionsHandler(predStore)
	mux.HandleFunc("POST /predictions", ph.Submit)
	lh := handlers.NewLeaderboardHandler(predStore, resultStore, "Columbus Crew")
	mux.HandleFunc("GET /leaderboard", lh.List)
	rh := handlers.NewResultsHandler(resultStore)
	mux.HandleFunc("POST /admin/results", rh.Submit)
	mux.HandleFunc("POST /auth/session", sh.Create)
	mux.HandleFunc("GET /auth/logout", handlers.Logout)
	mux.HandleFunc("GET /auth/config.js", serveFirebaseConfig)
	mux.HandleFunc("GET /login", func(w http.ResponseWriter, r *http.Request) {
		templates.Login().Render(r.Context(), w)
	})

	log.Printf("listening on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}

func serveFirebaseConfig(w http.ResponseWriter, r *http.Request) {
	cfg, _ := json.Marshal(map[string]string{
		"apiKey":    os.Getenv("FIREBASE_API_KEY"),
		"authDomain": os.Getenv("FIREBASE_AUTH_DOMAIN"),
		"projectId": os.Getenv("FIREBASE_PROJECT_ID"),
	})
	w.Header().Set("Content-Type", "application/javascript")
	fmt.Fprintf(w, "window.__firebaseConfig = %s;", cfg)
}
