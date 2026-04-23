package handlers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mcornell/crew-predictions/internal/handlers"
	"github.com/mcornell/crew-predictions/internal/models"
	"github.com/mcornell/crew-predictions/internal/repository"
)

func TestMatchDetailHandler_ReturnsPredictionsWithScores(t *testing.T) {
	predStore := repository.NewMemoryPredictionStore()
	resultStore := repository.NewMemoryResultStore()
	matchStore := repository.NewMemoryMatchStore()
	ctx := context.Background()

	match := models.Match{
		ID:        "m-test",
		HomeTeam:  "Columbus Crew",
		AwayTeam:  "FC Dallas",
		Kickoff:   time.Now().Add(-24 * time.Hour),
		Status:    "STATUS_FULL_TIME",
		HomeScore: "2",
		AwayScore: "1",
	}
	matchStore.Seed([]models.Match{match})

	predStore.Save(ctx, repository.Prediction{MatchID: "m-test", UserID: "google:u1", Handle: "fan1@bsky.mock", HomeGoals: 2, AwayGoals: 1})
	predStore.Save(ctx, repository.Prediction{MatchID: "m-test", UserID: "google:u2", Handle: "fan2@bsky.mock", HomeGoals: 0, AwayGoals: 0})

	resultStore.SaveResult(ctx, repository.Result{
		MatchID:   "m-test",
		HomeTeam:  "Columbus Crew",
		AwayTeam:  "FC Dallas",
		HomeGoals: 2,
		AwayGoals: 1,
	})

	h := handlers.NewMatchDetailHandler(predStore, resultStore, matchStore, repository.NewMemoryUserStore(), "Columbus Crew")

	req := httptest.NewRequest("GET", "/api/matches/m-test", nil)
	req.SetPathValue("matchId", "m-test")
	w := httptest.NewRecorder()
	h.Get(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp map[string]any
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	predictions, ok := resp["predictions"].([]any)
	if !ok {
		t.Fatal("expected predictions array in response")
	}
	if len(predictions) != 2 {
		t.Errorf("expected 2 predictions, got %d", len(predictions))
	}

	// fan1 predicted exactly right (2-1) — should have more Aces Radio points
	var fan1Points, fan2Points float64
	for _, p := range predictions {
		entry := p.(map[string]any)
		if entry["handle"] == "fan1@bsky.mock" {
			fan1Points = entry["acesRadioPoints"].(float64)
		}
		if entry["handle"] == "fan2@bsky.mock" {
			fan2Points = entry["acesRadioPoints"].(float64)
		}
	}
	if fan1Points <= fan2Points {
		t.Errorf("fan1 (exact) should have more points than fan2, got %v vs %v", fan1Points, fan2Points)
	}
}

func TestMatchDetailHandler_IncludesGrouchyPoints(t *testing.T) {
	predStore := repository.NewMemoryPredictionStore()
	resultStore := repository.NewMemoryResultStore()
	matchStore := repository.NewMemoryMatchStore()
	ctx := context.Background()

	matchStore.Seed([]models.Match{{
		ID: "m-grouchy", HomeTeam: "Columbus Crew", AwayTeam: "FC Dallas",
		Kickoff: time.Now().Add(-24 * time.Hour), Status: "STATUS_FULL_TIME",
		HomeScore: "2", AwayScore: "0",
	}})
	// Columbus home, predicted 3-0 (Win by 2+), actual 2-0 (Win by 2+) → 1 pt
	predStore.Save(ctx, repository.Prediction{MatchID: "m-grouchy", UserID: "u1", Handle: "GrouchyFan", HomeGoals: 3, AwayGoals: 0})
	resultStore.SaveResult(ctx, repository.Result{
		MatchID: "m-grouchy", HomeTeam: "Columbus Crew", AwayTeam: "FC Dallas", HomeGoals: 2, AwayGoals: 0,
	})

	h := handlers.NewMatchDetailHandler(predStore, resultStore, matchStore, repository.NewMemoryUserStore(), "Columbus Crew")
	req := httptest.NewRequest("GET", "/api/matches/m-grouchy", nil)
	req.SetPathValue("matchId", "m-grouchy")
	w := httptest.NewRecorder()
	h.Get(w, req)

	var resp map[string]any
	json.NewDecoder(w.Body).Decode(&resp)
	predictions := resp["predictions"].([]any)
	grouchy := predictions[0].(map[string]any)["grouchyPoints"].(float64)
	if grouchy != 1 {
		t.Errorf("expected grouchyPoints=1, got %v", grouchy)
	}
}

func TestMatchDetailHandler_ScoringFormatsIncludesGrouchy(t *testing.T) {
	predStore := repository.NewMemoryPredictionStore()
	resultStore := repository.NewMemoryResultStore()
	matchStore := repository.NewMemoryMatchStore()
	matchStore.Seed([]models.Match{{
		ID: "m-fmt", HomeTeam: "Columbus Crew", AwayTeam: "FC Dallas",
		Kickoff: time.Now().Add(-24 * time.Hour), Status: "STATUS_FULL_TIME",
		HomeScore: "1", AwayScore: "0",
	}})

	h := handlers.NewMatchDetailHandler(predStore, resultStore, matchStore, repository.NewMemoryUserStore(), "Columbus Crew")
	req := httptest.NewRequest("GET", "/api/matches/m-fmt", nil)
	req.SetPathValue("matchId", "m-fmt")
	w := httptest.NewRecorder()
	h.Get(w, req)

	var resp map[string]any
	json.NewDecoder(w.Body).Decode(&resp)
	formats := resp["scoringFormats"].([]any)
	keys := make([]string, 0, len(formats))
	for _, f := range formats {
		keys = append(keys, f.(map[string]any)["key"].(string))
	}
	found := false
	for _, k := range keys {
		if k == "grouchy" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected grouchy in scoringFormats, got %v", keys)
	}
}

func TestMatchDetailHandler_UsesCurrentHandleFromUserStore(t *testing.T) {
	predStore := repository.NewMemoryPredictionStore()
	resultStore := repository.NewMemoryResultStore()
	matchStore := repository.NewMemoryMatchStore()
	userStore := repository.NewMemoryUserStore()
	ctx := context.Background()

	matchStore.Seed([]models.Match{{
		ID: "m-handle", HomeTeam: "Columbus Crew", AwayTeam: "FC Dallas",
		Kickoff: time.Now().Add(-24 * time.Hour), Status: "STATUS_FULL_TIME",
		HomeScore: "1", AwayScore: "0",
	}})
	predStore.Save(ctx, repository.Prediction{MatchID: "m-handle", UserID: "google:u1", Handle: "My Full Name", HomeGoals: 1, AwayGoals: 0})
	userStore.Upsert(ctx, repository.User{UserID: "google:u1", Handle: "BlackAndGold"})

	h := handlers.NewMatchDetailHandler(predStore, resultStore, matchStore, userStore, "Columbus Crew")

	req := httptest.NewRequest("GET", "/api/matches/m-handle", nil)
	req.SetPathValue("matchId", "m-handle")
	w := httptest.NewRecorder()
	h.Get(w, req)

	var resp map[string]any
	json.NewDecoder(w.Body).Decode(&resp)
	predictions := resp["predictions"].([]any)
	handle := predictions[0].(map[string]any)["handle"].(string)
	if handle != "BlackAndGold" {
		t.Errorf("expected handle from UserStore %q, got %q", "BlackAndGold", handle)
	}
}

func TestMatchDetailHandler_FallsBackToStoredHandleWhenNoUserStoreEntry(t *testing.T) {
	predStore := repository.NewMemoryPredictionStore()
	resultStore := repository.NewMemoryResultStore()
	matchStore := repository.NewMemoryMatchStore()
	userStore := repository.NewMemoryUserStore()
	ctx := context.Background()

	matchStore.Seed([]models.Match{{
		ID: "m-legacy", HomeTeam: "Columbus Crew", AwayTeam: "FC Dallas",
		Kickoff: time.Now().Add(-24 * time.Hour), Status: "STATUS_FULL_TIME",
		HomeScore: "1", AwayScore: "0",
	}})
	predStore.Save(ctx, repository.Prediction{MatchID: "m-legacy", UserID: "google:legacy", Handle: "OldHandle", HomeGoals: 1, AwayGoals: 0})

	h := handlers.NewMatchDetailHandler(predStore, resultStore, matchStore, userStore, "Columbus Crew")

	req := httptest.NewRequest("GET", "/api/matches/m-legacy", nil)
	req.SetPathValue("matchId", "m-legacy")
	w := httptest.NewRecorder()
	h.Get(w, req)

	var resp map[string]any
	json.NewDecoder(w.Body).Decode(&resp)
	predictions := resp["predictions"].([]any)
	handle := predictions[0].(map[string]any)["handle"].(string)
	if handle != "OldHandle" {
		t.Errorf("expected fallback handle %q, got %q", "OldHandle", handle)
	}
}

func TestMatchDetailHandler_Returns404ForUnknownMatch(t *testing.T) {
	predStore := repository.NewMemoryPredictionStore()
	resultStore := repository.NewMemoryResultStore()
	matchStore := repository.NewMemoryMatchStore()
	userStore := repository.NewMemoryUserStore()

	h := handlers.NewMatchDetailHandler(predStore, resultStore, matchStore, userStore, "Columbus Crew")

	req := httptest.NewRequest("GET", "/api/matches/no-such-match", nil)
	req.SetPathValue("matchId", "no-such-match")
	w := httptest.NewRecorder()
	h.Get(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestMatchDetailHandler_ReturnsEmptyPredictionsWhenNone(t *testing.T) {
	predStore := repository.NewMemoryPredictionStore()
	resultStore := repository.NewMemoryResultStore()
	matchStore := repository.NewMemoryMatchStore()

	matchStore.Seed([]models.Match{{
		ID:        "m-empty",
		HomeTeam:  "Columbus Crew",
		AwayTeam:  "FC Dallas",
		Kickoff:   time.Now().Add(-24 * time.Hour),
		Status:    "STATUS_FULL_TIME",
		HomeScore: "1",
		AwayScore: "0",
	}})

	h := handlers.NewMatchDetailHandler(predStore, resultStore, matchStore, repository.NewMemoryUserStore(), "Columbus Crew")

	req := httptest.NewRequest("GET", "/api/matches/m-empty", nil)
	req.SetPathValue("matchId", "m-empty")
	w := httptest.NewRecorder()
	h.Get(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp map[string]any
	json.NewDecoder(w.Body).Decode(&resp)
	predictions := resp["predictions"].([]any)
	if len(predictions) != 0 {
		t.Errorf("expected 0 predictions, got %d", len(predictions))
	}
}
