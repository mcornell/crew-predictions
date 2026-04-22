package repository_test

import (
	"context"
	"testing"

	"github.com/mcornell/crew-predictions/internal/repository"
)

func TestMemoryUserStore_UpsertAndGetByID(t *testing.T) {
	s := repository.NewMemoryUserStore()
	ctx := context.Background()

	u := repository.User{UserID: "firebase:abc", Handle: "crewfan", Provider: "google"}
	if err := s.Upsert(ctx, u); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := s.GetByID(ctx, "firebase:abc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil || got.Handle != "crewfan" {
		t.Errorf("expected handle crewfan, got %+v", got)
	}
}

func TestMemoryUserStore_UpsertOverwritesHandle(t *testing.T) {
	s := repository.NewMemoryUserStore()
	ctx := context.Background()

	s.Upsert(ctx, repository.User{UserID: "firebase:abc", Handle: "old"})
	s.Upsert(ctx, repository.User{UserID: "firebase:abc", Handle: "new"})

	got, _ := s.GetByID(ctx, "firebase:abc")
	if got == nil || got.Handle != "new" {
		t.Errorf("expected handle new, got %+v", got)
	}
}

func TestMemoryUserStore_GetAll(t *testing.T) {
	s := repository.NewMemoryUserStore()
	ctx := context.Background()

	s.Upsert(ctx, repository.User{UserID: "u1", Handle: "fan1"})
	s.Upsert(ctx, repository.User{UserID: "u2", Handle: "fan2"})

	all, err := s.GetAll(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(all) != 2 {
		t.Errorf("expected 2 users, got %d", len(all))
	}
}

func TestMemoryUserStore_UpsertPreservesLocationWhenEmpty(t *testing.T) {
	s := repository.NewMemoryUserStore()
	ctx := context.Background()

	s.Upsert(ctx, repository.User{UserID: "u1", Handle: "fan", Location: "Columbus, OH"})
	s.Upsert(ctx, repository.User{UserID: "u1", Handle: "fan"}) // no location

	got, _ := s.GetByID(ctx, "u1")
	if got == nil || got.Location != "Columbus, OH" {
		t.Errorf("expected location Columbus, OH preserved, got %+v", got)
	}
}

func TestMemoryUserStore_UpsertPreservesProviderWhenEmpty(t *testing.T) {
	s := repository.NewMemoryUserStore()
	ctx := context.Background()

	s.Upsert(ctx, repository.User{UserID: "u1", Handle: "fan", Provider: "google.com"})
	s.Upsert(ctx, repository.User{UserID: "u1", Handle: "fan"}) // no provider

	got, _ := s.GetByID(ctx, "u1")
	if got == nil || got.Provider != "google.com" {
		t.Errorf("expected provider google.com preserved, got %+v", got)
	}
}

func TestMemoryUserStore_GetByID_NotFound(t *testing.T) {
	s := repository.NewMemoryUserStore()
	got, err := repository.NewMemoryUserStore().GetByID(context.Background(), "nope")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != nil {
		t.Errorf("expected nil for unknown user, got %+v", got)
	}
	_ = s
}
