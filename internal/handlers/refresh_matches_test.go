package handlers_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mcornell/crew-predictions/internal/handlers"
	"github.com/mcornell/crew-predictions/internal/models"
	"github.com/mcornell/crew-predictions/internal/repository"
)

func TestRefreshMatchesHandler_SavesFetcherResultsToStore(t *testing.T) {
	store := repository.NewMemoryMatchStore()
	matches := []models.Match{
		{ID: "m-espn", HomeTeam: "Columbus Crew", AwayTeam: "FC Dallas", Kickoff: time.Now()},
	}
	fetcher := func() ([]models.Match, error) { return matches, nil }

	h := handlers.NewRefreshMatchesHandler(store, fetcher)
	req := httptest.NewRequest(http.MethodPost, "/admin/refresh-matches", nil)
	w := httptest.NewRecorder()
	h.Refresh(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", w.Code)
	}
	got, _ := store.GetAll()
	if len(got) != 1 || got[0].ID != "m-espn" {
		t.Errorf("expected 1 match with id m-espn, got %+v", got)
	}
}

type errSaveMatchStore struct{ repository.MatchStore }

func (e *errSaveMatchStore) SaveAll(_ []models.Match) error { return fmt.Errorf("store write failed") }

func TestRefreshMatchesHandler_Returns500WhenStoreSaveFails(t *testing.T) {
	store := &errSaveMatchStore{MatchStore: repository.NewMemoryMatchStore()}
	fetcher := func() ([]models.Match, error) {
		return []models.Match{{ID: "m1", Kickoff: time.Now()}}, nil
	}

	h := handlers.NewRefreshMatchesHandler(store, fetcher)
	req := httptest.NewRequest(http.MethodPost, "/admin/refresh-matches", nil)
	w := httptest.NewRecorder()
	h.Refresh(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestRefreshMatchesHandler_Returns500WhenFetcherFails(t *testing.T) {
	store := repository.NewMemoryMatchStore()
	fetcher := func() ([]models.Match, error) { return nil, fmt.Errorf("espn down") }

	h := handlers.NewRefreshMatchesHandler(store, fetcher)
	req := httptest.NewRequest(http.MethodPost, "/admin/refresh-matches", nil)
	w := httptest.NewRecorder()
	h.Refresh(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}
