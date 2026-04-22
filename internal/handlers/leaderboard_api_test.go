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

func newLeaderboard(predictions repository.PredictionStore, results repository.ResultStore) *handlers.LeaderboardHandler {
	return handlers.NewLeaderboardHandler(predictions, results, repository.NewMemoryUserStore(), "Columbus Crew")
}

func TestLeaderboardAPIHandler_ReturnsJSON(t *testing.T) {
	predictions := repository.NewMemoryPredictionStore()
	results := repository.NewMemoryResultStore()
	ctx := context.Background()

	predictions.Save(ctx, repository.Prediction{MatchID: "m1", UserID: "u1", Handle: "BlackAndGold@bsky.mock", HomeGoals: 2, AwayGoals: 0})
	results.SaveResult(ctx, repository.Result{MatchID: "m1", HomeTeam: "Portland Timbers", AwayTeam: "Columbus Crew", HomeGoals: 2, AwayGoals: 0})

	lh := newLeaderboard(predictions, results)
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
	lh := handlers.NewLeaderboardHandler(store, repository.NewMemoryResultStore(), repository.NewMemoryUserStore(), "Columbus Crew")
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

	predictions.Save(ctx, repository.Prediction{MatchID: "m1", UserID: "u1", Handle: "ColumbusNordecke@bsky.mock", HomeGoals: 1, AwayGoals: 0})
	results.SaveResult(ctx, repository.Result{MatchID: "m1", HomeTeam: "Portland Timbers", AwayTeam: "Columbus Crew", HomeGoals: 3, AwayGoals: 0})

	lh := newLeaderboard(predictions, results)
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

func TestLeaderboardAPIHandler_UsesUserStoreHandleOverPredictionHandle(t *testing.T) {
	predictions := repository.NewMemoryPredictionStore()
	results := repository.NewMemoryResultStore()
	users := repository.NewMemoryUserStore()
	ctx := context.Background()

	predictions.Save(ctx, repository.Prediction{MatchID: "m1", UserID: "firebase:abc", Handle: "oldfan", HomeGoals: 2, AwayGoals: 0})
	results.SaveResult(ctx, repository.Result{MatchID: "m1", HomeTeam: "Portland Timbers", AwayTeam: "Columbus Crew", HomeGoals: 2, AwayGoals: 0})
	users.Upsert(ctx, repository.User{UserID: "firebase:abc", Handle: "CrewForever"})

	lh := handlers.NewLeaderboardHandler(predictions, results, users, "Columbus Crew")
	req := httptest.NewRequest(http.MethodGet, "/api/leaderboard", nil)
	w := httptest.NewRecorder()
	lh.APIList(w, req)

	var body struct {
		AcesRadio []struct {
			Handle string `json:"handle"`
			Points int    `json:"points"`
		} `json:"acesRadio"`
	}
	json.NewDecoder(w.Body).Decode(&body)
	if len(body.AcesRadio) == 0 || body.AcesRadio[0].Handle != "CrewForever" {
		t.Errorf("expected handle CrewForever from UserStore, got %+v", body.AcesRadio)
	}
}

func TestLeaderboardAPIHandler_IncludesUserID(t *testing.T) {
	predictions := repository.NewMemoryPredictionStore()
	results := repository.NewMemoryResultStore()
	users := repository.NewMemoryUserStore()
	ctx := context.Background()

	predictions.Save(ctx, repository.Prediction{MatchID: "m1", UserID: "firebase:abc", Handle: "crewfan", HomeGoals: 1, AwayGoals: 0})
	results.SaveResult(ctx, repository.Result{MatchID: "m1", HomeTeam: "Columbus Crew", AwayTeam: "FC Dallas", HomeGoals: 1, AwayGoals: 0})
	users.Upsert(ctx, repository.User{UserID: "firebase:abc", Handle: "crewfan"})

	lh := handlers.NewLeaderboardHandler(predictions, results, users, "Columbus Crew")
	req := httptest.NewRequest(http.MethodGet, "/api/leaderboard", nil)
	w := httptest.NewRecorder()
	lh.APIList(w, req)

	var body struct {
		AcesRadio []struct {
			UserID string `json:"userID"`
			Handle string `json:"handle"`
		} `json:"acesRadio"`
	}
	json.NewDecoder(w.Body).Decode(&body)
	if len(body.AcesRadio) == 0 || body.AcesRadio[0].UserID != "firebase:abc" {
		t.Errorf("expected userID in leaderboard response, got %+v", body.AcesRadio)
	}
}

func TestLeaderboardAPIHandler_ShowsUsersWithUnscoredPredictionsAtZero(t *testing.T) {
	predictions := repository.NewMemoryPredictionStore()
	results := repository.NewMemoryResultStore()
	users := repository.NewMemoryUserStore()
	ctx := context.Background()

	// u1 has a prediction but no result yet
	predictions.Save(ctx, repository.Prediction{MatchID: "m1", UserID: "u1", Handle: "EarlyFan", HomeGoals: 2, AwayGoals: 0})
	// u2 has a scored prediction
	predictions.Save(ctx, repository.Prediction{MatchID: "m2", UserID: "u2", Handle: "ScoredFan", HomeGoals: 2, AwayGoals: 0})
	results.SaveResult(ctx, repository.Result{MatchID: "m2", HomeTeam: "Columbus Crew", AwayTeam: "FC Dallas", HomeGoals: 1, AwayGoals: 0})
	users.Upsert(ctx, repository.User{UserID: "u1", Handle: "EarlyFan"})
	users.Upsert(ctx, repository.User{UserID: "u2", Handle: "ScoredFan"})

	lh := handlers.NewLeaderboardHandler(predictions, results, users, "Columbus Crew")
	req := httptest.NewRequest(http.MethodGet, "/api/leaderboard", nil)
	w := httptest.NewRecorder()
	lh.APIList(w, req)

	var body struct {
		AcesRadio []struct {
			UserID string `json:"userID"`
			Points int    `json:"points"`
		} `json:"acesRadio"`
	}
	json.NewDecoder(w.Body).Decode(&body)

	found := map[string]int{}
	for _, e := range body.AcesRadio {
		found[e.UserID] = e.Points
	}
	if _, ok := found["u1"]; !ok {
		t.Errorf("expected u1 (unscored prediction) to appear in leaderboard, got %+v", body.AcesRadio)
	}
	if found["u1"] != 0 {
		t.Errorf("expected u1 points=0, got %d", found["u1"])
	}
	if found["u2"] != 10 {
		t.Errorf("expected u2 points=10 (correct winner), got %d", found["u2"])
	}
}

func TestLeaderboardAPIHandler_HasProfileTrueWhenUserInStore(t *testing.T) {
	predictions := repository.NewMemoryPredictionStore()
	results := repository.NewMemoryResultStore()
	users := repository.NewMemoryUserStore()
	ctx := context.Background()

	// known: has a UserStore entry
	predictions.Save(ctx, repository.Prediction{MatchID: "m1", UserID: "firebase:abc", Handle: "known", HomeGoals: 1, AwayGoals: 0})
	users.Upsert(ctx, repository.User{UserID: "firebase:abc", Handle: "known"})
	// legacy: no UserStore entry, key falls back to handle
	predictions.Save(ctx, repository.Prediction{MatchID: "m2", UserID: "", Handle: "legacyfan", HomeGoals: 1, AwayGoals: 0})

	lh := handlers.NewLeaderboardHandler(predictions, results, users, "Columbus Crew")
	req := httptest.NewRequest(http.MethodGet, "/api/leaderboard", nil)
	w := httptest.NewRecorder()
	lh.APIList(w, req)

	var body struct {
		AcesRadio []struct {
			UserID     string `json:"userID"`
			Handle     string `json:"handle"`
			HasProfile bool   `json:"hasProfile"`
		} `json:"acesRadio"`
	}
	json.NewDecoder(w.Body).Decode(&body)

	byHandle := map[string]bool{}
	for _, e := range body.AcesRadio {
		byHandle[e.Handle] = e.HasProfile
	}
	if !byHandle["known"] {
		t.Errorf("expected hasProfile=true for user in UserStore, got %+v", body.AcesRadio)
	}
	if byHandle["legacyfan"] {
		t.Errorf("expected hasProfile=false for legacy handle-only user, got %+v", body.AcesRadio)
	}
}

func TestLeaderboardAPIHandler_FallsBackToPredictionHandleWhenNoUserRecord(t *testing.T) {
	predictions := repository.NewMemoryPredictionStore()
	results := repository.NewMemoryResultStore()
	ctx := context.Background()

	predictions.Save(ctx, repository.Prediction{MatchID: "m1", UserID: "firebase:abc", Handle: "legacyfan", HomeGoals: 2, AwayGoals: 0})
	results.SaveResult(ctx, repository.Result{MatchID: "m1", HomeTeam: "Portland Timbers", AwayTeam: "Columbus Crew", HomeGoals: 2, AwayGoals: 0})

	lh := newLeaderboard(predictions, results) // empty UserStore
	req := httptest.NewRequest(http.MethodGet, "/api/leaderboard", nil)
	w := httptest.NewRecorder()
	lh.APIList(w, req)

	var body struct {
		AcesRadio []struct {
			Handle string `json:"handle"`
		} `json:"acesRadio"`
	}
	json.NewDecoder(w.Body).Decode(&body)
	if len(body.AcesRadio) == 0 || body.AcesRadio[0].Handle != "legacyfan" {
		t.Errorf("expected fallback to prediction handle legacyfan, got %+v", body.AcesRadio)
	}
}
