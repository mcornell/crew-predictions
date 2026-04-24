package handlers

import (
	"context"
	"net/http"

	"github.com/mcornell/crew-predictions/internal/models"
	"github.com/mcornell/crew-predictions/internal/poll"
	"github.com/mcornell/crew-predictions/internal/repository"
)

type PollScoresHandler struct {
	matchStore  repository.MatchStore
	resultStore repository.ResultStore
	fetcher     func() ([]models.Match, error)
	recalcFn    func(context.Context)
}

func NewPollScoresHandler(matchStore repository.MatchStore, resultStore repository.ResultStore, fetcher func() ([]models.Match, error), recalcFn func(context.Context)) *PollScoresHandler {
	return &PollScoresHandler{matchStore: matchStore, resultStore: resultStore, fetcher: fetcher, recalcFn: recalcFn}
}

func (h *PollScoresHandler) Poll(w http.ResponseWriter, r *http.Request) {
	if err := poll.PollOnce(r.Context(), h.matchStore, h.resultStore, h.fetcher); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.recalcFn(r.Context())
	w.WriteHeader(http.StatusNoContent)
}
