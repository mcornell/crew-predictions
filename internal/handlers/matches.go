package handlers

import (
	"net/http"

	"github.com/mcornell/crew-predictions/internal/espn"
	"github.com/mcornell/crew-predictions/templates"
)

func Matches(w http.ResponseWriter, r *http.Request) {
	matches, err := espn.FetchCrewMatches()
	if err != nil {
		http.Error(w, "couldn't fetch matches, try again", http.StatusInternalServerError)
		return
	}
	templates.MatchList(matches).Render(r.Context(), w)
}
