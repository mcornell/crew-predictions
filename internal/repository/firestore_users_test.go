package repository_test

import (
	"context"
	"os"
	"testing"

	"github.com/mcornell/crew-predictions/internal/repository"
)

func firestoreUserStoreOrSkip(t *testing.T) *repository.FirestoreUserStore {
	t.Helper()
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		t.Skip("FIRESTORE_EMULATOR_HOST not set; skipping Firestore integration tests")
	}
	store, err := repository.NewFirestoreUserStore(context.Background(), "crew-predictions-test")
	if err != nil {
		t.Fatalf("failed to create FirestoreUserStore: %v", err)
	}
	return store
}

func TestFirestoreUserStore_UpsertPreservesProviderWhenEmpty(t *testing.T) {
	store := firestoreUserStoreOrSkip(t)
	ctx := context.Background()

	_ = store.Upsert(ctx, repository.User{UserID: "fs-user-provider", Handle: "fan", Provider: "google.com"})
	_ = store.Upsert(ctx, repository.User{UserID: "fs-user-provider", Handle: "fan"}) // empty provider

	got, err := store.GetByID(ctx, "fs-user-provider")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil || got.Provider != "google.com" {
		t.Errorf("expected provider google.com preserved, got %+v", got)
	}
}

func TestFirestoreUserStore_UpsertPreservesLocationWhenEmpty(t *testing.T) {
	store := firestoreUserStoreOrSkip(t)
	ctx := context.Background()

	_ = store.Upsert(ctx, repository.User{UserID: "fs-user-location", Handle: "fan", Location: "Columbus, OH"})
	_ = store.Upsert(ctx, repository.User{UserID: "fs-user-location", Handle: "fan"}) // no location

	got, err := store.GetByID(ctx, "fs-user-location")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil || got.Location != "Columbus, OH" {
		t.Errorf("expected location Columbus, OH preserved, got %+v", got)
	}
}

func TestFirestoreUserStore_UpsertUpdatesHandle(t *testing.T) {
	store := firestoreUserStoreOrSkip(t)
	ctx := context.Background()

	_ = store.Upsert(ctx, repository.User{UserID: "fs-user-handle", Handle: "old", Provider: "google.com"})
	_ = store.Upsert(ctx, repository.User{UserID: "fs-user-handle", Handle: "new", Provider: "google.com"})

	got, err := store.GetByID(ctx, "fs-user-handle")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil || got.Handle != "new" {
		t.Errorf("expected handle new, got %+v", got)
	}
}
