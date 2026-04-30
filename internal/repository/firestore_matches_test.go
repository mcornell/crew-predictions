//go:build integration

package repository_test

import (
	"os"
	"testing"
	"time"

	"github.com/mcornell/crew-predictions/internal/models"
	"github.com/mcornell/crew-predictions/internal/repository"
)

func firestoreMatchStoreOrSkip(t *testing.T) *repository.FirestoreMatchStore {
	t.Helper()
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		t.Skip("FIRESTORE_EMULATOR_HOST not set; skipping Firestore integration tests")
	}
	store, err := repository.NewFirestoreMatchStore("crew-predictions-test")
	if err != nil {
		t.Fatalf("failed to create FirestoreMatchStore: %v", err)
	}
	return store
}

func TestFirestoreMatchStore_SaveAllAndGetAll(t *testing.T) {
	store := firestoreMatchStoreOrSkip(t)

	matches := []models.Match{
		{ID: "fs-match-1", HomeTeam: "Columbus Crew", AwayTeam: "FC Dallas", Kickoff: time.Date(2025, 4, 1, 20, 0, 0, 0, time.UTC)},
		{ID: "fs-match-2", HomeTeam: "Columbus Crew", AwayTeam: "LA Galaxy", Kickoff: time.Date(2025, 4, 8, 19, 0, 0, 0, time.UTC)},
	}
	if err := store.SaveAll(matches); err != nil {
		t.Fatalf("unexpected error saving: %v", err)
	}

	got, err := store.GetAll()
	if err != nil {
		t.Fatalf("unexpected error retrieving: %v", err)
	}
	found := map[string]bool{}
	for _, m := range got {
		found[m.ID] = true
	}
	if !found["fs-match-1"] || !found["fs-match-2"] {
		t.Errorf("expected both matches, got IDs: %v", found)
	}
}

func TestFirestoreMatchStore_VenueRoundTrips(t *testing.T) {
	store := firestoreMatchStoreOrSkip(t)

	if err := store.SaveAll([]models.Match{
		{ID: "fs-match-venue", HomeTeam: "Columbus Crew", AwayTeam: "FC Dallas",
			Kickoff: time.Date(2026, 5, 1, 20, 0, 0, 0, time.UTC),
			Venue:   "ScottsMiracle-Gro Field"},
	}); err != nil {
		t.Fatalf("unexpected error saving: %v", err)
	}

	got, err := store.GetAll()
	if err != nil {
		t.Fatalf("unexpected error retrieving: %v", err)
	}
	for _, m := range got {
		if m.ID == "fs-match-venue" {
			if m.Venue != "ScottsMiracle-Gro Field" {
				t.Errorf("expected Venue 'ScottsMiracle-Gro Field', got %q", m.Venue)
			}
			return
		}
	}
	t.Error("match fs-match-venue not found in GetAll results")
}

func TestFirestoreMatchStore_RecordsAndFormRoundTrip(t *testing.T) {
	store := firestoreMatchStoreOrSkip(t)

	if err := store.SaveAll([]models.Match{
		{ID: "fs-match-form", HomeTeam: "Columbus Crew", AwayTeam: "FC Dallas",
			Kickoff:    time.Date(2026, 5, 1, 20, 0, 0, 0, time.UTC),
			HomeRecord: "5-3-2", AwayRecord: "4-4-2",
			HomeForm: "WWWLL", AwayForm: "LWDWL"},
	}); err != nil {
		t.Fatalf("unexpected error saving: %v", err)
	}

	got, err := store.GetAll()
	if err != nil {
		t.Fatalf("unexpected error retrieving: %v", err)
	}
	for _, m := range got {
		if m.ID == "fs-match-form" {
			if m.HomeRecord != "5-3-2" {
				t.Errorf("HomeRecord: got %q, want %q", m.HomeRecord, "5-3-2")
			}
			if m.AwayRecord != "4-4-2" {
				t.Errorf("AwayRecord: got %q, want %q", m.AwayRecord, "4-4-2")
			}
			if m.HomeForm != "WWWLL" {
				t.Errorf("HomeForm: got %q, want %q", m.HomeForm, "WWWLL")
			}
			if m.AwayForm != "LWDWL" {
				t.Errorf("AwayForm: got %q, want %q", m.AwayForm, "LWDWL")
			}
			return
		}
	}
	t.Error("match fs-match-form not found in GetAll results")
}

func TestFirestoreMatchStore_SaveAllOverwritesExisting(t *testing.T) {
	store := firestoreMatchStoreOrSkip(t)

	_ = store.SaveAll([]models.Match{
		{ID: "fs-match-overwrite", HomeTeam: "Columbus Crew", AwayTeam: "Old Opponent", Kickoff: time.Now()},
	})
	_ = store.SaveAll([]models.Match{
		{ID: "fs-match-overwrite", HomeTeam: "Columbus Crew", AwayTeam: "New Opponent", Kickoff: time.Now()},
	})

	got, err := store.GetAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, m := range got {
		if m.ID == "fs-match-overwrite" && m.AwayTeam != "New Opponent" {
			t.Errorf("expected overwritten record, got AwayTeam=%q", m.AwayTeam)
		}
	}
}
