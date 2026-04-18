package handlers_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/mcornell/crew-predictions/internal/handlers"
	"github.com/mcornell/crew-predictions/internal/models"
	"github.com/mcornell/crew-predictions/internal/repository"
)

func stubFetcher(matches []models.Match) func() ([]models.Match, error) {
	return func() ([]models.Match, error) { return matches, nil }
}

func oneMatch() []models.Match {
	return []models.Match{{ID: "match-99", HomeTeam: "Columbus Crew", AwayTeam: "FC Cincinnati", Kickoff: time.Now()}}
}

func TestMatchesHandler_ReturnsOK(t *testing.T) {
	mh := handlers.NewMatchesHandler(repository.NewMemoryPredictionStore(), stubFetcher(oneMatch()))
	req := httptest.NewRequest(http.MethodGet, "/matches", nil)
	w := httptest.NewRecorder()

	mh.List(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestMatchesHandler_ReturnsHTML(t *testing.T) {
	mh := handlers.NewMatchesHandler(repository.NewMemoryPredictionStore(), stubFetcher(oneMatch()))
	req := httptest.NewRequest(http.MethodGet, "/matches", nil)
	w := httptest.NewRecorder()

	mh.List(w, req)

	ct := w.Header().Get("Content-Type")
	if !strings.Contains(ct, "text/html") {
		t.Errorf("expected text/html content-type, got %s", ct)
	}
}

func TestMatchesHandler_ReturnsInternalErrorWhenFetchFails(t *testing.T) {
	mh := handlers.NewMatchesHandler(repository.NewMemoryPredictionStore(), func() ([]models.Match, error) {
		return nil, fmt.Errorf("espn is down")
	})
	req := httptest.NewRequest(http.MethodGet, "/matches", nil)
	w := httptest.NewRecorder()

	mh.List(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestMatchesHandler_ShowsSavedPrediction(t *testing.T) {
	store := repository.NewMemoryPredictionStore()
	store.Save(context.Background(), repository.Prediction{
		MatchID: "match-99", UserID: "google:abc123", Handle: "BlackAndGold@bsky.mock", HomeGoals: 3, AwayGoals: 1,
	})
	mh := handlers.NewMatchesHandler(store, stubFetcher(oneMatch()))
	req := httptest.NewRequest(http.MethodGet, "/matches", nil)
	req.AddCookie(sessionCookie("google:abc123", "BlackAndGold@bsky.mock"))
	w := httptest.NewRecorder()

	mh.List(w, req)

	if !strings.Contains(w.Body.String(), "3") || !strings.Contains(w.Body.String(), "Your Pick") {
		t.Errorf("expected saved score in response, got: %s", w.Body.String())
	}
}
