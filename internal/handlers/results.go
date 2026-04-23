package handlers

import (
	"net/http"
	"strconv"

	"github.com/mcornell/crew-predictions/internal/repository"
)

type ResultsHandler struct {
	store repository.ResultStore
}

func NewResultsHandler(store repository.ResultStore) *ResultsHandler {
	return &ResultsHandler{store: store}
}

func (h *ResultsHandler) Submit(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	matchID := r.FormValue("match_id")
	home, err := strconv.Atoi(r.FormValue("home_goals"))
	if err != nil {
		http.Error(w, "home_goals must be an integer", http.StatusBadRequest)
		return
	}
	away, err := strconv.Atoi(r.FormValue("away_goals"))
	if err != nil {
		http.Error(w, "away_goals must be an integer", http.StatusBadRequest)
		return
	}
	if home < 0 || home > 99 || away < 0 || away > 99 {
		http.Error(w, "goals must be between 0 and 99", http.StatusBadRequest)
		return
	}
	h.store.SaveResult(r.Context(), repository.Result{
		MatchID:   matchID,
		HomeTeam:  r.FormValue("home_team"),
		AwayTeam:  r.FormValue("away_team"),
		HomeGoals: home,
		AwayGoals: away,
	})
	w.WriteHeader(http.StatusNoContent)
}
