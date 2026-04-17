package handlers

import (
	"fmt"
	"net/http"

	"github.com/mcornell/crew-predictions/internal/models"
)

func Matches(w http.ResponseWriter, r *http.Request) {
	MatchesWithData(w, r, []models.Match{})
}

func MatchesWithData(w http.ResponseWriter, r *http.Request, matches []models.Match) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "<h1>Upcoming Matches</h1>")
	for _, m := range matches {
		fmt.Fprintf(w, `<div data-testid="match-card">%s vs %s</div>`, m.HomeTeam, m.AwayTeam)
	}
}
