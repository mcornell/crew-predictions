package handlers

import (
	"fmt"
	"net/http"

	"github.com/mcornell/crew-predictions/internal/espn"
	"github.com/mcornell/crew-predictions/internal/models"
)

func Matches(w http.ResponseWriter, r *http.Request) {
	matches, err := espn.FetchCrewMatches()
	if err != nil {
		http.Error(w, "couldn't fetch matches, try again", http.StatusInternalServerError)
		return
	}
	renderMatches(w, matches)
}

func renderMatches(w http.ResponseWriter, matches []models.Match) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "<h1>Upcoming Matches</h1>")
	for _, m := range matches {
		fmt.Fprintf(w, `<div data-testid="match-card">%s vs %s</div>`, m.HomeTeam, m.AwayTeam)
	}
}
