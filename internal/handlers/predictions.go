package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/mcornell/crew-predictions/internal/repository"
)

type PredictionsHandler struct {
	store repository.PredictionStore
}

func NewPredictionsHandler(store repository.PredictionStore) *PredictionsHandler {
	return &PredictionsHandler{store: store}
}

func (h *PredictionsHandler) Submit(w http.ResponseWriter, r *http.Request) {
	handle := userFromSession(r)
	if handle == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	r.ParseForm()
	if r.FormValue("match_id") == "" || r.FormValue("home_goals") == "" || r.FormValue("away_goals") == "" {
		http.Error(w, "missing fields", http.StatusBadRequest)
		return
	}
	home, err := strconv.Atoi(r.FormValue("home_goals"))
	if err != nil {
		http.Error(w, "goals must be integers", http.StatusBadRequest)
		return
	}
	away, err := strconv.Atoi(r.FormValue("away_goals"))
	if err != nil {
		http.Error(w, "goals must be integers", http.StatusBadRequest)
		return
	}
	h.store.Save(r.Context(), repository.Prediction{
		MatchID:   r.FormValue("match_id"),
		Handle:    handle,
		HomeGoals: home,
		AwayGoals: away,
	})
	if r.Header.Get("HX-Request") == "true" {
		fmt.Fprintf(w, `<div data-testid="match-card"><div class="saved-score">%d – %d</div><div class="saved-label">Your Pick</div></div>`, home, away)
		return
	}
	http.Redirect(w, r, "/matches", http.StatusFound)
}

func SubmitPrediction(w http.ResponseWriter, r *http.Request) {
	if userFromSession(r) == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	r.ParseForm()
	if r.FormValue("match_id") == "" || r.FormValue("home_goals") == "" || r.FormValue("away_goals") == "" {
		http.Error(w, "missing fields", http.StatusBadRequest)
		return
	}
	home, err := strconv.Atoi(r.FormValue("home_goals"))
	if err != nil {
		http.Error(w, "goals must be integers", http.StatusBadRequest)
		return
	}
	away, err := strconv.Atoi(r.FormValue("away_goals"))
	if err != nil {
		http.Error(w, "goals must be integers", http.StatusBadRequest)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		fmt.Fprintf(w, `<div data-testid="match-card"><div class="saved-score">%d – %d</div><div class="saved-label">Your Pick</div></div>`, home, away)
		return
	}
	http.Redirect(w, r, "/matches", http.StatusFound)
}
