package repository_test

import (
	"context"
	"os"
	"testing"

	"github.com/mcornell/crew-predictions/internal/repository"
)

func firestoreOrSkip(t *testing.T) *repository.FirestorePredictionStore {
	t.Helper()
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		t.Skip("FIRESTORE_EMULATOR_HOST not set; skipping Firestore integration tests")
	}
	store, err := repository.NewFirestorePredictionStore(context.Background(), "crew-predictions-test")
	if err != nil {
		t.Fatalf("failed to create FirestorePredictionStore: %v", err)
	}
	return store
}

func TestFirestorePredictionStore_SaveAndRetrieve(t *testing.T) {
	store := firestoreOrSkip(t)
	ctx := context.Background()

	pred := repository.Prediction{
		MatchID:   "firestore-match-1",
		Handle:    "ColumbusNordecke@bsky.mock",
		HomeGoals: 2,
		AwayGoals: 0,
	}

	if err := store.Save(ctx, pred); err != nil {
		t.Fatalf("unexpected error saving: %v", err)
	}

	got, err := store.GetByMatchAndHandle(ctx, "firestore-match-1", "ColumbusNordecke@bsky.mock")
	if err != nil {
		t.Fatalf("unexpected error retrieving: %v", err)
	}
	if got == nil {
		t.Fatal("expected prediction, got nil")
	}
	if got.HomeGoals != 2 || got.AwayGoals != 0 {
		t.Errorf("expected 2-0, got %d-%d", got.HomeGoals, got.AwayGoals)
	}
}

func TestFirestorePredictionStore_ReturnsNilWhenNotFound(t *testing.T) {
	store := firestoreOrSkip(t)
	got, err := store.GetByMatchAndHandle(context.Background(), "no-such-match", "nobody@bsky.mock")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != nil {
		t.Errorf("expected nil for missing prediction, got %+v", got)
	}
}
