package handlers

import (
	"net/http"

	"github.com/mcornell/crew-predictions/internal/models"
	"github.com/mcornell/crew-predictions/internal/poll"
	"github.com/mcornell/crew-predictions/internal/repository"
)

type PollScoresHandler struct {
	matchStore  repository.MatchStore
	resultStore repository.ResultStore
	fetcher     func() ([]models.Match, error)
}

func NewPollScoresHandler(matchStore repository.MatchStore, resultStore repository.ResultStore, fetcher func() ([]models.Match, error)) *PollScoresHandler {
	return &PollScoresHandler{matchStore: matchStore, resultStore: resultStore, fetcher: fetcher}
}

func (h *PollScoresHandler) Poll(w http.ResponseWriter, r *http.Request) {
	if err := poll.PollOnce(r.Context(), h.matchStore, h.resultStore, h.fetcher); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
