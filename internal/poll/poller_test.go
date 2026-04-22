package poll_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/mcornell/crew-predictions/internal/models"
	"github.com/mcornell/crew-predictions/internal/poll"
	"github.com/mcornell/crew-predictions/internal/repository"
)

func TestPollOnce_WritesResultForTerminalMatch(t *testing.T) {
	matchStore := repository.NewMemoryMatchStore()
	resultStore := repository.NewMemoryResultStore()

	matches := []models.Match{
		{
			ID:        "m-done",
			HomeTeam:  "Columbus Crew",
			AwayTeam:  "FC Dallas",
			Status:    "STATUS_FULL_TIME",
			State:     "post",
			HomeScore: "3",
			AwayScore: "1",
		},
	}
	fetcher := func() ([]models.Match, error) { return matches, nil }

	if err := poll.PollOnce(context.Background(), matchStore, resultStore, fetcher); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err := resultStore.GetResult(context.Background(), "m-done")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected result to be saved, got nil")
	}
	if result.HomeGoals != 3 || result.AwayGoals != 1 {
		t.Errorf("expected 3-1, got %d-%d", result.HomeGoals, result.AwayGoals)
	}
}

func TestPollOnce_SkipsNonTerminalMatch(t *testing.T) {
	matchStore := repository.NewMemoryMatchStore()
	resultStore := repository.NewMemoryResultStore()

	matches := []models.Match{
		{ID: "m-live", HomeTeam: "Columbus Crew", AwayTeam: "FC Dallas", Status: "STATUS_IN_PROGRESS", State: "in", HomeScore: "1", AwayScore: "0"},
	}
	fetcher := func() ([]models.Match, error) { return matches, nil }

	if err := poll.PollOnce(context.Background(), matchStore, resultStore, fetcher); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, _ := resultStore.GetResult(context.Background(), "m-live")
	if result != nil {
		t.Errorf("expected no result saved for in-progress match, got %+v", result)
	}
}

func TestPollOnce_UpdatesMatchStoreWithFetchedData(t *testing.T) {
	matchStore := repository.NewMemoryMatchStore()
	resultStore := repository.NewMemoryResultStore()

	matches := []models.Match{
		{ID: "m-live", HomeTeam: "Columbus Crew", AwayTeam: "FC Dallas", Status: "STATUS_IN_PROGRESS", State: "in", HomeScore: "2", AwayScore: "0"},
	}
	fetcher := func() ([]models.Match, error) { return matches, nil }

	if err := poll.PollOnce(context.Background(), matchStore, resultStore, fetcher); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	stored, _ := matchStore.GetAll()
	if len(stored) != 1 || stored[0].HomeScore != "2" {
		t.Errorf("expected matchStore updated with live score, got %+v", stored)
	}
}

func TestPollOnce_WritesResultForFinalAET(t *testing.T) {
	matchStore := repository.NewMemoryMatchStore()
	resultStore := repository.NewMemoryResultStore()

	matches := []models.Match{
		{ID: "m-aet", HomeTeam: "Columbus Crew", AwayTeam: "FC Dallas", Status: "STATUS_FINAL_AET", State: "post", HomeScore: "2", AwayScore: "2"},
	}
	fetcher := func() ([]models.Match, error) { return matches, nil }

	poll.PollOnce(context.Background(), matchStore, resultStore, fetcher)

	result, _ := resultStore.GetResult(context.Background(), "m-aet")
	if result == nil || result.HomeGoals != 2 || result.AwayGoals != 2 {
		t.Errorf("expected result for STATUS_FINAL_AET, got %+v", result)
	}
}

func TestPollOnce_WritesResultForFinalPEN(t *testing.T) {
	matchStore := repository.NewMemoryMatchStore()
	resultStore := repository.NewMemoryResultStore()

	matches := []models.Match{
		{ID: "m-pen", HomeTeam: "Columbus Crew", AwayTeam: "FC Dallas", Status: "STATUS_FINAL_PEN", State: "post", HomeScore: "1", AwayScore: "1"},
	}
	fetcher := func() ([]models.Match, error) { return matches, nil }

	poll.PollOnce(context.Background(), matchStore, resultStore, fetcher)

	result, _ := resultStore.GetResult(context.Background(), "m-pen")
	if result == nil {
		t.Errorf("expected result for STATUS_FINAL_PEN, got nil")
	}
}

func TestPollOnce_ReturnsErrorWhenSaveAllFails(t *testing.T) {
	resultStore := repository.NewMemoryResultStore()
	matches := []models.Match{
		{ID: "m1", HomeTeam: "Columbus Crew", AwayTeam: "FC Dallas", Status: "STATUS_FULL_TIME", HomeScore: "1", AwayScore: "0"},
	}
	fetcher := func() ([]models.Match, error) { return matches, nil }

	errStore := &errMatchStore{}
	err := poll.PollOnce(context.Background(), errStore, resultStore, fetcher)
	if err == nil {
		t.Error("expected error when SaveAll fails, got nil")
	}
}

type errMatchStore struct{}

func (e *errMatchStore) GetAll() ([]models.Match, error)    { return nil, nil }
func (e *errMatchStore) SaveAll(_ []models.Match) error     { return fmt.Errorf("store failed") }
func (e *errMatchStore) Reset()                             {}

func TestPollOnce_IsIdempotent(t *testing.T) {
	matchStore := repository.NewMemoryMatchStore()
	resultStore := repository.NewMemoryResultStore()

	matches := []models.Match{
		{ID: "m-done", HomeTeam: "Columbus Crew", AwayTeam: "FC Dallas", Status: "STATUS_FULL_TIME", State: "post", HomeScore: "1", AwayScore: "0"},
	}
	fetcher := func() ([]models.Match, error) { return matches, nil }

	poll.PollOnce(context.Background(), matchStore, resultStore, fetcher)
	if err := poll.PollOnce(context.Background(), matchStore, resultStore, fetcher); err != nil {
		t.Fatalf("second poll failed: %v", err)
	}

	result, _ := resultStore.GetResult(context.Background(), "m-done")
	if result == nil || result.HomeGoals != 1 || result.AwayGoals != 0 {
		t.Errorf("expected idempotent result, got %+v", result)
	}
}
