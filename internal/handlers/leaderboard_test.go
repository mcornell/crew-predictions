package handlers_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mcornell/crew-predictions/internal/handlers"
	"github.com/mcornell/crew-predictions/internal/repository"
)

func newLeaderboardHandler() *handlers.LeaderboardHandler {
	return handlers.NewLeaderboardHandler(
		repository.NewMemoryPredictionStore(),
		repository.NewMemoryResultStore(),
		"Columbus Crew",
	)
}

func TestLeaderboardHandler_ReturnsOK(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/leaderboard", nil)
	w := httptest.NewRecorder()

	newLeaderboardHandler().List(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestLeaderboardHandler_ShowsUpper90ClubPoints(t *testing.T) {
	predictions := repository.NewMemoryPredictionStore()
	results := repository.NewMemoryResultStore()
	ctx := context.Background()

	// Correct winner, Columbus is away — Upper90Club: +1
	predictions.Save(ctx, repository.Prediction{MatchID: "m1", Handle: "ColumbusNordecke@bsky.mock", HomeGoals: 1, AwayGoals: 0})
	results.SaveResult(ctx, repository.Result{MatchID: "m1", HomeTeam: "Portland", AwayTeam: "Columbus Crew", HomeGoals: 3, AwayGoals: 0})

	lh := handlers.NewLeaderboardHandler(predictions, results, "Columbus Crew")
	req := httptest.NewRequest(http.MethodGet, "/leaderboard", nil)
	w := httptest.NewRecorder()
	lh.List(w, req)

	body := w.Body.String()
	if !strings.Contains(body, "ColumbusNordecke@bsky.mock") {
		t.Errorf("expected handle in leaderboard, got: %s", body)
	}
	if !strings.Contains(body, "1") {
		t.Errorf("expected 1 Upper90Club point in leaderboard, got: %s", body)
	}
}

func TestLeaderboardHandler_ShowsPointsForExactScore(t *testing.T) {
	predictions := repository.NewMemoryPredictionStore()
	results := repository.NewMemoryResultStore()
	ctx := context.Background()

	predictions.Save(ctx, repository.Prediction{MatchID: "m1", Handle: "BlackAndGold@bsky.mock", HomeGoals: 2, AwayGoals: 0})
	results.SaveResult(ctx, repository.Result{MatchID: "m1", HomeTeam: "Portland Timbers", AwayTeam: "Columbus Crew", HomeGoals: 2, AwayGoals: 0})

	lh := handlers.NewLeaderboardHandler(predictions, results, "Columbus Crew")
	req := httptest.NewRequest(http.MethodGet, "/leaderboard", nil)
	w := httptest.NewRecorder()
	lh.List(w, req)

	body := w.Body.String()
	if !strings.Contains(body, "BlackAndGold@bsky.mock") {
		t.Errorf("expected handle in leaderboard, got: %s", body)
	}
	if !strings.Contains(body, "15") {
		t.Errorf("expected 15 points in leaderboard, got: %s", body)
	}
}
