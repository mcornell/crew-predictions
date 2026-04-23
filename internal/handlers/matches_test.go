package handlers_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mcornell/crew-predictions/internal/handlers"
	"github.com/mcornell/crew-predictions/internal/models"
	"github.com/mcornell/crew-predictions/internal/repository"
)

func oneMatch() []models.Match {
	return []models.Match{{ID: "match-99", HomeTeam: "Portland Timbers", AwayTeam: "Columbus Crew", Kickoff: time.Now()}}
}

func matchStoreWith(matches []models.Match) *repository.MemoryMatchStore {
	ms := repository.NewMemoryMatchStore()
	ms.SaveAll(matches) //nolint
	return ms
}

func TestAPIMatchesHandler_ReadsFromMatchStore(t *testing.T) {
	mh := handlers.NewMatchesHandler(repository.NewMemoryPredictionStore(), matchStoreWith(oneMatch()))
	req := httptest.NewRequest(http.MethodGet, "/api/matches", nil)
	w := httptest.NewRecorder()

	mh.APIList(w, req)

	var body struct {
		Matches []struct{ ID string `json:"id"` } `json:"matches"`
	}
	json.NewDecoder(w.Body).Decode(&body)
	if len(body.Matches) != 1 || body.Matches[0].ID != "match-99" {
		t.Errorf("expected match-99, got %+v", body.Matches)
	}
}

func TestAPIMatchesHandler_ReturnsJSON(t *testing.T) {
	mh := handlers.NewMatchesHandler(repository.NewMemoryPredictionStore(), matchStoreWith(oneMatch()))
	req := httptest.NewRequest(http.MethodGet, "/api/matches", nil)
	w := httptest.NewRecorder()

	mh.APIList(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected application/json, got %s", ct)
	}
	var body struct {
		Matches []struct {
			ID       string `json:"id"`
			HomeTeam string `json:"homeTeam"`
			AwayTeam string `json:"awayTeam"`
		} `json:"matches"`
	}
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(body.Matches) != 1 || body.Matches[0].AwayTeam != "Columbus Crew" {
		t.Errorf("unexpected matches: %+v", body.Matches)
	}
}

type errMatchStore struct{}

func (e *errMatchStore) GetAll() ([]models.Match, error) { return nil, fmt.Errorf("store error") }
func (e *errMatchStore) SaveAll(_ []models.Match) error  { return nil }
func (e *errMatchStore) Reset()                          {}

func TestAPIMatchesHandler_ReturnsErrorWhenFetchFails(t *testing.T) {
	mh := handlers.NewMatchesHandler(repository.NewMemoryPredictionStore(), &errMatchStore{})
	req := httptest.NewRequest(http.MethodGet, "/api/matches", nil)
	w := httptest.NewRecorder()

	mh.APIList(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestAPIMatchesHandler_IncludesStateInResponse(t *testing.T) {
	ms := repository.NewMemoryMatchStore()
	ms.SaveAll([]models.Match{{ID: "m-live", HomeTeam: "Columbus Crew", AwayTeam: "FC Dallas", Kickoff: time.Now(), State: "in"}})
	mh := handlers.NewMatchesHandler(repository.NewMemoryPredictionStore(), ms)
	req := httptest.NewRequest(http.MethodGet, "/api/matches", nil)
	w := httptest.NewRecorder()

	mh.APIList(w, req)

	var body struct {
		Matches []struct {
			ID    string `json:"id"`
			State string `json:"state"`
		} `json:"matches"`
	}
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(body.Matches) != 1 || body.Matches[0].State != "in" {
		t.Errorf("expected state=in, got %+v", body.Matches)
	}
}

func TestAPIMatchesHandler_ReturnsSortedByKickoffAscending(t *testing.T) {
	later := time.Now().Add(48 * time.Hour)
	earlier := time.Now().Add(24 * time.Hour)
	ms := repository.NewMemoryMatchStore()
	ms.SaveAll([]models.Match{
		{ID: "m-late", HomeTeam: "Columbus Crew", AwayTeam: "FC Dallas", Kickoff: later},
		{ID: "m-early", HomeTeam: "Columbus Crew", AwayTeam: "LA Galaxy", Kickoff: earlier},
	})
	mh := handlers.NewMatchesHandler(repository.NewMemoryPredictionStore(), ms)
	req := httptest.NewRequest(http.MethodGet, "/api/matches", nil)
	w := httptest.NewRecorder()

	mh.APIList(w, req)

	var body struct {
		Matches []struct{ ID string `json:"id"` } `json:"matches"`
	}
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(body.Matches) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(body.Matches))
	}
	if body.Matches[0].ID != "m-early" || body.Matches[1].ID != "m-late" {
		t.Errorf("expected m-early before m-late, got %v then %v", body.Matches[0].ID, body.Matches[1].ID)
	}
}

func TestAPIMatchesHandler_IncludesPredictionForLoggedInUser(t *testing.T) {
	store := repository.NewMemoryPredictionStore()
	store.Save(context.Background(), repository.Prediction{
		MatchID: "match-99", UserID: "google:abc123", HomeGoals: 2, AwayGoals: 1,
	})
	mh := handlers.NewMatchesHandler(store, matchStoreWith(oneMatch()))
	req := httptest.NewRequest(http.MethodGet, "/api/matches", nil)
	req.AddCookie(sessionCookie("google:abc123", "BlackAndGold@bsky.mock"))
	w := httptest.NewRecorder()

	mh.APIList(w, req)

	var body struct {
		Predictions map[string]struct {
			HomeGoals int `json:"homeGoals"`
			AwayGoals int `json:"awayGoals"`
		} `json:"predictions"`
	}
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	p, ok := body.Predictions["match-99"]
	if !ok {
		t.Fatal("expected prediction for match-99 in response")
	}
	if p.HomeGoals != 2 || p.AwayGoals != 1 {
		t.Errorf("expected 2-1, got %d-%d", p.HomeGoals, p.AwayGoals)
	}
}
