package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/mcornell/crew-predictions/internal/models"
	"github.com/mcornell/crew-predictions/internal/poll"
	"github.com/mcornell/crew-predictions/internal/repository"
	"github.com/mcornell/crew-predictions/internal/tasks"
)

// pollChainInterval is how far in the future the next chain task is scheduled
// when a match remains in a non-terminal state. Kept short so the live clock
// stays current.
const pollChainInterval = 2 * time.Minute

// terminalStatuses mirrors internal/poll.terminalStatuses plus the "won't
// resume today" states. Chain ends on any of these; the next 4am/12pm/6pm
// refresh picks up rescheduled kickoffs and seeds a fresh chain if needed.
var terminalStatuses = map[string]bool{
	"STATUS_FULL_TIME":  true,
	"STATUS_FINAL_AET":  true,
	"STATUS_FINAL_PEN":  true,
	"STATUS_POSTPONED":  true,
	"STATUS_CANCELED":   true,
	"STATUS_CANCELLED":  true, // ESPN inconsistency between US/UK spelling
	"STATUS_ABANDONED":  true,
	"STATUS_FORFEIT":    true,
}

type PollScoresHandler struct {
	matchStore  repository.MatchStore
	resultStore repository.ResultStore
	fetcher     func() ([]models.Match, error)
	recalcFn    func(context.Context)
	enqueuer    tasks.Enqueuer // nil = no chain continuation (legacy callers, tests)
}

func NewPollScoresHandler(matchStore repository.MatchStore, resultStore repository.ResultStore, fetcher func() ([]models.Match, error), recalcFn func(context.Context)) *PollScoresHandler {
	return &PollScoresHandler{matchStore: matchStore, resultStore: resultStore, fetcher: fetcher, recalcFn: recalcFn}
}

// WithEnqueuer attaches a chain-task enqueuer. When set, Poll consults the
// just-polled match for the matchID query param and schedules a follow-up
// task ~2 min out if the match is still in a non-terminal state. Without an
// enqueuer the handler retains its legacy "poll all matches, no chain"
// behavior (suitable for ad-hoc admin triggers and tests).
func (h *PollScoresHandler) WithEnqueuer(e tasks.Enqueuer) *PollScoresHandler {
	h.enqueuer = e
	return h
}

func (h *PollScoresHandler) Poll(w http.ResponseWriter, r *http.Request) {
	if err := poll.PollOnce(r.Context(), h.matchStore, h.resultStore, h.fetcher); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.recalcFn(r.Context())

	if matchID := r.URL.Query().Get("matchID"); matchID != "" && h.enqueuer != nil {
		h.maybeEnqueueNext(r.Context(), matchID)
	}

	w.WriteHeader(http.StatusNoContent)
}

// maybeEnqueueNext consults the freshly-polled store for the named match and
// schedules a follow-up task at now + pollChainInterval if the match is
// still in a non-terminal state. A missing match or a terminal status ends
// the chain (no enqueue).
func (h *PollScoresHandler) maybeEnqueueNext(ctx context.Context, matchID string) {
	matches, err := h.matchStore.GetAll()
	if err != nil {
		slog.Error("poll_scores: chain continuation read failed", "matchID", matchID, "error", err)
		return
	}
	for _, m := range matches {
		if m.ID != matchID {
			continue
		}
		if terminalStatuses[m.Status] {
			return
		}
		runAt := time.Now().UTC().Add(pollChainInterval)
		if err := h.enqueuer.EnqueuePoll(ctx, matchID, runAt); err != nil {
			slog.Error("poll_scores: chain enqueue failed", "matchID", matchID, "error", err)
		}
		return
	}
}
