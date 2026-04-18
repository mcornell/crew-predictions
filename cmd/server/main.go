package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/mcornell/crew-predictions/internal/espn"
	"github.com/mcornell/crew-predictions/internal/handlers"
	"github.com/mcornell/crew-predictions/internal/repository"
)

func main() {
	// load .env in development; ignored if file doesn't exist (e.g. production)
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	var predStore repository.PredictionStore
	if project := os.Getenv("GOOGLE_CLOUD_PROJECT"); project != "" {
		fs, err := repository.NewFirestorePredictionStore(context.Background(), project)
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
	mux.HandleFunc("GET /auth/login", handlers.Login)
	mux.HandleFunc("GET /auth/callback", handlers.Callback)
	mux.HandleFunc("GET /auth/logout", handlers.Logout)

	log.Printf("listening on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
