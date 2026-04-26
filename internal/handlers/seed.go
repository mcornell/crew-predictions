package handlers

import (
	"net/http"
	"strconv"

	"github.com/mcornell/crew-predictions/internal/repository"
)

type SeedPredictionHandler struct {
	store repository.PredictionStore
}

func NewSeedPredictionHandler(store repository.PredictionStore) *SeedPredictionHandler {
	return &SeedPredictionHandler{store: store}
}

func (h *SeedPredictionHandler) Submit(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
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
	h.store.Save(r.Context(), repository.Prediction{
		MatchID:   r.FormValue("match_id"),
		UserID:    r.FormValue("user_id"),
		HomeGoals: home,
		AwayGoals: away,
	})
	w.WriteHeader(http.StatusNoContent)
}

type SeedUserHandler struct {
	store repository.UserStore
}

func NewSeedUserHandler(store repository.UserStore) *SeedUserHandler {
	return &SeedUserHandler{store: store}
}

func (h *SeedUserHandler) Submit(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userID := r.FormValue("user_id")
	handle := r.FormValue("handle")
	if userID == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}
	h.store.Upsert(r.Context(), repository.User{UserID: userID, Handle: handle})
	w.WriteHeader(http.StatusNoContent)
}
