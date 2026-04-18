package repository_test

import (
	"context"
	"testing"

	"github.com/mcornell/crew-predictions/internal/repository"
)

func TestMemoryUserStore_ReturnsNilWhenNotFound(t *testing.T) {
	store := repository.NewMemoryUserStore()
	got, err := store.GetUser(context.Background(), "google:nobody")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != nil {
		t.Errorf("expected nil for missing user, got %+v", got)
	}
}

func TestMemoryUserStore_UpsertAndRetrieve(t *testing.T) {
	store := repository.NewMemoryUserStore()
	ctx := context.Background()

	user := repository.User{
		UserID:   "google:110048215615",
		Handle:   "BlackAndGold@bsky.mock",
		Provider: "google",
	}
	if err := store.Upsert(ctx, user); err != nil {
		t.Fatalf("unexpected error upserting: %v", err)
	}

	got, err := store.GetUser(ctx, "google:110048215615")
	if err != nil {
		t.Fatalf("unexpected error retrieving: %v", err)
	}
	if got == nil {
		t.Fatal("expected user, got nil")
	}
	if got.Handle != "BlackAndGold@bsky.mock" {
		t.Errorf("expected handle BlackAndGold@bsky.mock, got %s", got.Handle)
	}
}

func TestMemoryUserStore_UpsertUpdatesHandle(t *testing.T) {
	store := repository.NewMemoryUserStore()
	ctx := context.Background()

	store.Upsert(ctx, repository.User{UserID: "google:123", Handle: "old@bsky.mock", Provider: "google"})
	store.Upsert(ctx, repository.User{UserID: "google:123", Handle: "new@bsky.mock", Provider: "google"})

	got, _ := store.GetUser(ctx, "google:123")
	if got.Handle != "new@bsky.mock" {
		t.Errorf("expected handle to be updated, got %s", got.Handle)
	}
}
