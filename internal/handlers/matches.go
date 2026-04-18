package handlers

import (
	"net/http"

	"github.com/mcornell/crew-predictions/internal/models"
	"github.com/mcornell/crew-predictions/internal/repository"
	"github.com/mcornell/crew-predictions/templates"
)

type MatchesHandler struct {
	store   repository.PredictionStore
	fetcher func() ([]models.Match, error)
}

func NewMatchesHandler(store repository.PredictionStore, fetcher func() ([]models.Match, error)) *MatchesHandler {
	return &MatchesHandler{store: store, fetcher: fetcher}
}

func (h *MatchesHandler) List(w http.ResponseWriter, r *http.Request) {
	matches, err := h.fetcher()
	if err != nil {
		http.Error(w, "couldn't fetch matches, try again", http.StatusInternalServerError)
		return
	}
	user := userFromSession(r)
	handle := ""
	predictions := map[string]*repository.Prediction{}
	if user != nil {
		handle = user.Handle
		for _, m := range matches {
			p, _ := h.store.GetByMatchAndUser(r.Context(), m.ID, user.UserID)
			if p != nil {
				predictions[m.ID] = p
			}
		}
	}
	templates.MatchList(matches, handle, predictions).Render(r.Context(), w)
}
