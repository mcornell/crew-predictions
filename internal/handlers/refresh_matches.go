package handlers

import (
	"net/http"

	"github.com/mcornell/crew-predictions/internal/models"
	"github.com/mcornell/crew-predictions/internal/repository"
)

type RefreshMatchesHandler struct {
	store     repository.MatchStore
	fetcher   func() ([]models.Match, error)
	onRefresh func([]models.Match)
}

func NewRefreshMatchesHandler(store repository.MatchStore, fetcher func() ([]models.Match, error), onRefresh func([]models.Match)) *RefreshMatchesHandler {
	return &RefreshMatchesHandler{store: store, fetcher: fetcher, onRefresh: onRefresh}
}

func (h *RefreshMatchesHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	matches, err := h.fetcher()
	if err != nil {
		http.Error(w, "couldn't fetch matches", http.StatusInternalServerError)
		return
	}
	if err := h.store.SaveAll(matches); err != nil {
		http.Error(w, "couldn't save matches", http.StatusInternalServerError)
		return
	}
	if h.onRefresh != nil {
		h.onRefresh(matches)
	}
	w.WriteHeader(http.StatusNoContent)
}
