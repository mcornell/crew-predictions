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
)

func TestPollScoresHandler_Returns500WhenFetcherFails(t *testing.T) {
	matchStore := repository.NewMemoryMatchStore()
	resultStore := repository.NewMemoryResultStore()
	fetcher := func() ([]models.Match, error) { return nil, fmt.Errorf("espn down") }

	h := handlers.NewPollScoresHandler(matchStore, resultStore, fetcher)
	req := httptest.NewRequest(http.MethodPost, "/admin/poll-scores", nil)
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

	h := handlers.NewPollScoresHandler(matchStore, resultStore, fetcher)
	req := httptest.NewRequest(http.MethodPost, "/admin/poll-scores", nil)
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
