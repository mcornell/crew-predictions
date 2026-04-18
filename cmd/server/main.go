package main

import (
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

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/matches", http.StatusFound)
	})
	store := repository.NewMemoryPredictionStore()
	mh := handlers.NewMatchesHandler(store, espn.FetchCrewMatches)
	mux.HandleFunc("GET /matches", mh.List)
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	ph := handlers.NewPredictionsHandler(store)
	mux.HandleFunc("POST /predictions", ph.Submit)
	mux.HandleFunc("GET /auth/login", handlers.Login)
	mux.HandleFunc("GET /auth/callback", handlers.Callback)
	mux.HandleFunc("GET /auth/logout", handlers.Logout)

	log.Printf("listening on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
