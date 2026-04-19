package handlers

import (
	"net/http"
	"time"

	"github.com/mcornell/crew-predictions/internal/models"
	"github.com/mcornell/crew-predictions/internal/repository"
)

type SeedMatchHandler struct {
	store *repository.MemoryMatchStore
}

func NewSeedMatchHandler(store *repository.MemoryMatchStore) *SeedMatchHandler {
	return &SeedMatchHandler{store: store}
}

func (h *SeedMatchHandler) Submit(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	kickoff, err := time.Parse(time.RFC3339, r.FormValue("kickoff"))
	if err != nil {
		http.Error(w, "kickoff must be RFC3339", http.StatusBadRequest)
		return
	}
	h.store.Seed([]models.Match{{
		ID:       r.FormValue("id"),
		HomeTeam: r.FormValue("home_team"),
		AwayTeam: r.FormValue("away_team"),
		Kickoff:  kickoff,
		Status:   r.FormValue("status"),
	}})
	w.WriteHeader(http.StatusNoContent)
}
