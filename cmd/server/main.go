package main

import (
	"log"
	"net/http"
	"os"

	"github.com/mcornell/crew-predictions/internal/handlers"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/matches", http.StatusFound)
	})
	mux.HandleFunc("GET /matches", handlers.Matches)

	log.Printf("listening on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
