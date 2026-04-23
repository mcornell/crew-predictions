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

type leaderboardBody struct {
	Entries []struct {
		UserID          string `json:"userID"`
		Handle          string `json:"handle"`
		AcesRadioPoints int    `json:"acesRadioPoints"`
		Upper90Points   int    `json:"upper90ClubPoints"`
		HasProfile      bool   `json:"hasProfile"`
	} `json:"entries"`
}

func decodeLeaderboard(t *testing.T, w *httptest.ResponseRecorder) leaderboardBody {
	t.Helper()
	var body leaderboardBody
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	return body
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
	body := decodeLeaderboard(t, w)
	if len(body.Entries) == 0 || body.Entries[0].Handle != "BlackAndGold@bsky.mock" || body.Entries[0].AcesRadioPoints != 15 {
		t.Errorf("unexpected entries: %+v", body.Entries)
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

	body := decodeLeaderboard(t, w)
	if len(body.Entries) == 0 || body.Entries[0].Handle != "ColumbusNordecke@bsky.mock" || body.Entries[0].Upper90Points != 2 {
		t.Errorf("unexpected entries: %+v", body.Entries)
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

	body := decodeLeaderboard(t, w)
	if len(body.Entries) == 0 || body.Entries[0].Handle != "CrewForever" {
		t.Errorf("expected handle CrewForever from UserStore, got %+v", body.Entries)
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

	body := decodeLeaderboard(t, w)
	if len(body.Entries) == 0 || body.Entries[0].UserID != "firebase:abc" {
		t.Errorf("expected userID in leaderboard response, got %+v", body.Entries)
	}
}

func TestLeaderboardAPIHandler_ShowsUsersWithUnscoredPredictionsAtZero(t *testing.T) {
	predictions := repository.NewMemoryPredictionStore()
	results := repository.NewMemoryResultStore()
	users := repository.NewMemoryUserStore()
	ctx := context.Background()

	predictions.Save(ctx, repository.Prediction{MatchID: "m1", UserID: "u1", Handle: "EarlyFan", HomeGoals: 2, AwayGoals: 0})
	predictions.Save(ctx, repository.Prediction{MatchID: "m2", UserID: "u2", Handle: "ScoredFan", HomeGoals: 2, AwayGoals: 0})
	results.SaveResult(ctx, repository.Result{MatchID: "m2", HomeTeam: "Columbus Crew", AwayTeam: "FC Dallas", HomeGoals: 1, AwayGoals: 0})
	users.Upsert(ctx, repository.User{UserID: "u1", Handle: "EarlyFan"})
	users.Upsert(ctx, repository.User{UserID: "u2", Handle: "ScoredFan"})

	lh := handlers.NewLeaderboardHandler(predictions, results, users, "Columbus Crew")
	req := httptest.NewRequest(http.MethodGet, "/api/leaderboard", nil)
	w := httptest.NewRecorder()
	lh.APIList(w, req)

	body := decodeLeaderboard(t, w)

	found := map[string]int{}
	for _, e := range body.Entries {
		found[e.UserID] = e.AcesRadioPoints
	}
	if _, ok := found["u1"]; !ok {
		t.Errorf("expected u1 (unscored prediction) to appear in leaderboard, got %+v", body.Entries)
	}
	if found["u1"] != 0 {
		t.Errorf("expected u1 points=0, got %d", found["u1"])
	}
	if found["u2"] != 10 {
		t.Errorf("expected u2 acesRadioPoints=10 (correct winner), got %d", found["u2"])
	}
}

func TestLeaderboardAPIHandler_HasProfileTrueWhenUserInStore(t *testing.T) {
	predictions := repository.NewMemoryPredictionStore()
	results := repository.NewMemoryResultStore()
	users := repository.NewMemoryUserStore()
	ctx := context.Background()

	predictions.Save(ctx, repository.Prediction{MatchID: "m1", UserID: "firebase:abc", Handle: "known", HomeGoals: 1, AwayGoals: 0})
	users.Upsert(ctx, repository.User{UserID: "firebase:abc", Handle: "known"})
	predictions.Save(ctx, repository.Prediction{MatchID: "m2", UserID: "", Handle: "legacyfan", HomeGoals: 1, AwayGoals: 0})

	lh := handlers.NewLeaderboardHandler(predictions, results, users, "Columbus Crew")
	req := httptest.NewRequest(http.MethodGet, "/api/leaderboard", nil)
	w := httptest.NewRecorder()
	lh.APIList(w, req)

	body := decodeLeaderboard(t, w)

	byHandle := map[string]bool{}
	for _, e := range body.Entries {
		byHandle[e.Handle] = e.HasProfile
	}
	if !byHandle["known"] {
		t.Errorf("expected hasProfile=true for user in UserStore, got %+v", body.Entries)
	}
	if byHandle["legacyfan"] {
		t.Errorf("expected hasProfile=false for legacy handle-only user, got %+v", body.Entries)
	}
}

func TestLeaderboardAPIHandler_FallsBackToPredictionHandleWhenNoUserRecord(t *testing.T) {
	predictions := repository.NewMemoryPredictionStore()
	results := repository.NewMemoryResultStore()
	ctx := context.Background()

	predictions.Save(ctx, repository.Prediction{MatchID: "m1", UserID: "firebase:abc", Handle: "legacyfan", HomeGoals: 2, AwayGoals: 0})
	results.SaveResult(ctx, repository.Result{MatchID: "m1", HomeTeam: "Portland Timbers", AwayTeam: "Columbus Crew", HomeGoals: 2, AwayGoals: 0})

	lh := newLeaderboard(predictions, results)
	req := httptest.NewRequest(http.MethodGet, "/api/leaderboard", nil)
	w := httptest.NewRecorder()
	lh.APIList(w, req)

	body := decodeLeaderboard(t, w)
	if len(body.Entries) == 0 || body.Entries[0].Handle != "legacyfan" {
		t.Errorf("expected fallback to prediction handle legacyfan, got %+v", body.Entries)
	}
}
