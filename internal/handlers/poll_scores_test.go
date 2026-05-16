package handlers_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mcornell/crew-predictions/internal/handlers"
	"github.com/mcornell/crew-predictions/internal/models"
	"github.com/mcornell/crew-predictions/internal/repository"
	"github.com/mcornell/crew-predictions/internal/tasks"
)

func TestPollScoresHandler_CallsRecalcFnAfterSuccessfulPoll(t *testing.T) {
	matchStore := repository.NewMemoryMatchStore()
	matchStore.Seed([]models.Match{{
		ID: "m1", HomeTeam: "Columbus Crew", AwayTeam: "FC Dallas",
		Kickoff: time.Now().Add(-2 * time.Hour), Status: "STATUS_FULL_TIME", State: "post",
		HomeScore: "2", AwayScore: "0",
	}})
	fetcher := func() ([]models.Match, error) { return matchStore.GetAll() }

	called := 0
	recalcFn := func(_ context.Context) { called++ }
	h := handlers.NewPollScoresHandler(matchStore, repository.NewMemoryResultStore(), fetcher, recalcFn)
	req := httptest.NewRequest(http.MethodPost, "/admin/poll-scores", http.NoBody)
	w := httptest.NewRecorder()
	h.Poll(w, req)

	if called != 1 {
		t.Errorf("expected recalcFn called once, got %d", called)
	}
}

func TestPollScoresHandler_DoesNotCallRecalcFnOnFetcherFailure(t *testing.T) {
	called := 0
	recalcFn := func(_ context.Context) { called++ }
	fetcher := func() ([]models.Match, error) { return nil, fmt.Errorf("espn down") }
	h := handlers.NewPollScoresHandler(repository.NewMemoryMatchStore(), repository.NewMemoryResultStore(), fetcher, recalcFn)
	req := httptest.NewRequest(http.MethodPost, "/admin/poll-scores", http.NoBody)
	w := httptest.NewRecorder()
	h.Poll(w, req)
	if called != 0 {
		t.Errorf("expected recalcFn not called on fetch failure, got %d", called)
	}
}

func TestPollScoresHandler_Returns500WhenFetcherFails(t *testing.T) {
	matchStore := repository.NewMemoryMatchStore()
	resultStore := repository.NewMemoryResultStore()
	fetcher := func() ([]models.Match, error) { return nil, fmt.Errorf("espn down") }

	h := handlers.NewPollScoresHandler(matchStore, resultStore, fetcher, func(_ context.Context) {})
	req := httptest.NewRequest(http.MethodPost, "/admin/poll-scores", http.NoBody)
	w := httptest.NewRecorder()
	h.Poll(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestPollScoresHandler_CallsFetcherAndWritesResult(t *testing.T) {
	matchStore := repository.NewMemoryMatchStore()
	resultStore := repository.NewMemoryResultStore()

	matchStore.Seed([]models.Match{{
		ID: "m-done", HomeTeam: "Columbus Crew", AwayTeam: "FC Dallas",
		Kickoff: time.Now().Add(-2 * time.Hour), Status: "STATUS_FULL_TIME", State: "post",
		HomeScore: "2", AwayScore: "0",
	}})

	fetcher := func() ([]models.Match, error) { return matchStore.GetAll() }

	h := handlers.NewPollScoresHandler(matchStore, resultStore, fetcher, func(_ context.Context) {})
	req := httptest.NewRequest(http.MethodPost, "/admin/poll-scores", http.NoBody)
	w := httptest.NewRecorder()

	h.Poll(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", w.Code)
	}
	result, err := resultStore.GetResult(context.Background(), "m-done")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || result.HomeGoals != 2 || result.AwayGoals != 0 {
		t.Errorf("expected result 2-0, got %+v", result)
	}
}

func TestPollScoresHandler_EnqueuesNextTaskForInProgressMatch(t *testing.T) {
	matchStore := repository.NewMemoryMatchStore()
	matchStore.Seed([]models.Match{{
		ID: "m-live", HomeTeam: "Columbus Crew", AwayTeam: "LAFC",
		Kickoff: time.Now().Add(-20 * time.Minute), Status: "STATUS_IN_PROGRESS", State: "in",
	}})
	fetcher := func() ([]models.Match, error) { return matchStore.GetAll() }
	enqueuer := tasks.NewFakeEnqueuer()

	h := handlers.NewPollScoresHandler(matchStore, repository.NewMemoryResultStore(), fetcher, func(_ context.Context) {}).
		WithEnqueuer(enqueuer)

	req := httptest.NewRequest(http.MethodPost, "/admin/poll-scores?matchID=m-live", http.NoBody)
	w := httptest.NewRecorder()
	before := time.Now()
	h.Poll(w, req)
	after := time.Now()

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d: %s", w.Code, w.Body.String())
	}
	calls := enqueuer.Calls()
	if len(calls) != 1 {
		t.Fatalf("expected exactly one enqueued task for in-progress match, got %d", len(calls))
	}
	if calls[0].MatchID != "m-live" {
		t.Errorf("enqueued task matchID: got %q, want %q", calls[0].MatchID, "m-live")
	}
	minRunAt := before.Add(2 * time.Minute)
	maxRunAt := after.Add(2 * time.Minute)
	if calls[0].RunAt.Before(minRunAt) || calls[0].RunAt.After(maxRunAt) {
		t.Errorf("enqueued RunAt %v outside expected window [%v, %v]", calls[0].RunAt, minRunAt, maxRunAt)
	}
}

func TestPollScoresHandler_DoesNotEnqueueWhenMatchIsTerminal(t *testing.T) {
	matchStore := repository.NewMemoryMatchStore()
	matchStore.Seed([]models.Match{{
		ID: "m-done", HomeTeam: "Columbus Crew", AwayTeam: "LAFC",
		Kickoff: time.Now().Add(-2 * time.Hour), Status: "STATUS_FULL_TIME", State: "post",
		HomeScore: "2", AwayScore: "1",
	}})
	fetcher := func() ([]models.Match, error) { return matchStore.GetAll() }
	enqueuer := tasks.NewFakeEnqueuer()

	h := handlers.NewPollScoresHandler(matchStore, repository.NewMemoryResultStore(), fetcher, func(_ context.Context) {}).
		WithEnqueuer(enqueuer)

	req := httptest.NewRequest(http.MethodPost, "/admin/poll-scores?matchID=m-done", http.NoBody)
	w := httptest.NewRecorder()
	h.Poll(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
	if got := len(enqueuer.Calls()); got != 0 {
		t.Errorf("expected no enqueue for terminal match, got %d call(s)", got)
	}
}
