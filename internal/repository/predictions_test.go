package repository_test

import (
	"context"
	"testing"

	"github.com/mcornell/crew-predictions/internal/repository"
)

func TestMemoryPredictionStore_ReturnsNilWhenNotFound(t *testing.T) {
	store := repository.NewMemoryPredictionStore()
	got, err := store.GetByMatchAndHandle(context.Background(), "no-match", "nobody")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != nil {
		t.Errorf("expected nil for missing prediction, got %+v", got)
	}
}

func TestMemoryPredictionStore_GetAllReturnsAllSaved(t *testing.T) {
	store := repository.NewMemoryPredictionStore()
	ctx := context.Background()

	store.Save(ctx, repository.Prediction{MatchID: "m1", Handle: "a@bsky.mock", HomeGoals: 1, AwayGoals: 0})
	store.Save(ctx, repository.Prediction{MatchID: "m2", Handle: "b@bsky.mock", HomeGoals: 2, AwayGoals: 1})

	all, err := store.GetAll(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(all) != 2 {
		t.Errorf("expected 2 predictions, got %d", len(all))
	}
}

func TestMemoryPredictionStore_SaveAndRetrieve(t *testing.T) {
	store := repository.NewMemoryPredictionStore()
	ctx := context.Background()

	pred := repository.Prediction{
		MatchID:   "match1",
		Handle:    "BlackAndGold@bsky.mock",
		HomeGoals: 3,
		AwayGoals: 1,
	}

	if err := store.Save(ctx, pred); err != nil {
		t.Fatalf("unexpected error saving: %v", err)
	}

	got, err := store.GetByMatchAndHandle(ctx, "match1", "BlackAndGold@bsky.mock")
	if err != nil {
		t.Fatalf("unexpected error retrieving: %v", err)
	}
	if got.HomeGoals != 3 || got.AwayGoals != 1 {
		t.Errorf("expected 3-1, got %d-%d", got.HomeGoals, got.AwayGoals)
	}
}
