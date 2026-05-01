package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/mcornell/crew-predictions/internal/models"
	"github.com/mcornell/crew-predictions/internal/repository"
)

type SeedMatchEventsHandler struct {
	store *repository.MemoryMatchStore
}

func NewSeedMatchEventsHandler(store *repository.MemoryMatchStore) *SeedMatchEventsHandler {
	return &SeedMatchEventsHandler{store: store}
}

type seedMatchEventsBody struct {
	MatchID string              `json:"matchID"`
	Events  []models.MatchEvent `json:"events"`
}

func (h *SeedMatchEventsHandler) Submit(w http.ResponseWriter, r *http.Request) {
	var body seedMatchEventsBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	matches, err := h.store.GetAll()
	if err != nil {
		http.Error(w, "could not load matches", http.StatusInternalServerError)
		return
	}
	found := false
	for i, m := range matches {
		if m.ID == body.MatchID {
			matches[i].Events = body.Events
			found = true
			break
		}
	}
	if !found {
		http.Error(w, "match not found", http.StatusNotFound)
		return
	}
	if err := h.store.SaveAll(matches); err != nil {
		http.Error(w, "could not save matches", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
