package handlers_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mcornell/crew-predictions/internal/handlers"
	"github.com/mcornell/crew-predictions/internal/repository"
)

type errGetAllPredictionStore struct{ repository.PredictionStore }

func (e *errGetAllPredictionStore) GetAll(_ context.Context) ([]repository.Prediction, error) {
	return nil, fmt.Errorf("store down")
}

func TestLeaderboardAPIHandler_ReturnsJSON(t *testing.T) {
	predictions := repository.NewMemoryPredictionStore()
	results := repository.NewMemoryResultStore()
	ctx := context.Background()

	predictions.Save(ctx, repository.Prediction{MatchID: "m1", Handle: "BlackAndGold@bsky.mock", HomeGoals: 2, AwayGoals: 0})
	results.SaveResult(ctx, repository.Result{MatchID: "m1", HomeTeam: "Portland Timbers", AwayTeam: "Columbus Crew", HomeGoals: 2, AwayGoals: 0})

	lh := handlers.NewLeaderboardHandler(predictions, results, "Columbus Crew")
	req := httptest.NewRequest(http.MethodGet, "/api/leaderboard", nil)
	w := httptest.NewRecorder()

	lh.APIList(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var body struct {
		AcesRadio []struct {
			Handle string `json:"handle"`
			Points int    `json:"points"`
		} `json:"acesRadio"`
	}
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(body.AcesRadio) == 0 || body.AcesRadio[0].Handle != "BlackAndGold@bsky.mock" || body.AcesRadio[0].Points != 15 {
		t.Errorf("unexpected acesRadio: %+v", body.AcesRadio)
	}
}

func TestLeaderboardAPIHandler_Returns500WhenGetAllFails(t *testing.T) {
	store := &errGetAllPredictionStore{PredictionStore: repository.NewMemoryPredictionStore()}
	lh := handlers.NewLeaderboardHandler(store, repository.NewMemoryResultStore(), "Columbus Crew")
	req := httptest.NewRequest(http.MethodGet, "/api/leaderboard", nil)
	w := httptest.NewRecorder()

	lh.APIList(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestLeaderboardAPIHandler_IncludesUpper90Club(t *testing.T) {
	predictions := repository.NewMemoryPredictionStore()
	results := repository.NewMemoryResultStore()
	ctx := context.Background()

	predictions.Save(ctx, repository.Prediction{MatchID: "m1", Handle: "ColumbusNordecke@bsky.mock", HomeGoals: 1, AwayGoals: 0})
	results.SaveResult(ctx, repository.Result{MatchID: "m1", HomeTeam: "Portland Timbers", AwayTeam: "Columbus Crew", HomeGoals: 3, AwayGoals: 0})

	lh := handlers.NewLeaderboardHandler(predictions, results, "Columbus Crew")
	req := httptest.NewRequest(http.MethodGet, "/api/leaderboard", nil)
	w := httptest.NewRecorder()

	lh.APIList(w, req)

	var body struct {
		Upper90Club []struct {
			Handle string `json:"handle"`
			Points int    `json:"points"`
		} `json:"upper90Club"`
	}
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(body.Upper90Club) == 0 || body.Upper90Club[0].Handle != "ColumbusNordecke@bsky.mock" || body.Upper90Club[0].Points != 2 {
		t.Errorf("unexpected upper90Club: %+v", body.Upper90Club)
	}
}
