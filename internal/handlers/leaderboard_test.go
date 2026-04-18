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

func TestLeaderboardHandler_ShowsPointsForExactScore(t *testing.T) {
	predictions := repository.NewMemoryPredictionStore()
	results := repository.NewMemoryResultStore()
	ctx := context.Background()

	predictions.Save(ctx, repository.Prediction{MatchID: "m1", Handle: "BlackAndGold@bsky.mock", HomeGoals: 2, AwayGoals: 0})
	results.SaveResult(ctx, repository.Result{MatchID: "m1", HomeGoals: 2, AwayGoals: 0})

	lh := handlers.NewLeaderboardHandler(predictions, results)
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
