package handlers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mcornell/crew-predictions/internal/handlers"
	"github.com/mcornell/crew-predictions/internal/repository"
)

func newLeaderboard(users repository.UserStore) *handlers.LeaderboardHandler {
	return handlers.NewLeaderboardHandler(repository.NewMemoryPredictionStore(), repository.NewMemoryResultStore(), users, "Columbus Crew")
}

type leaderboardBody struct {
	Entries []struct {
		UserID          string `json:"userID"`
		Handle          string `json:"handle"`
		AcesRadioPoints int    `json:"acesRadioPoints"`
		Upper90Points   int    `json:"upper90ClubPoints"`
		GrouchyPoints   int    `json:"grouchyPoints"`
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

func TestLeaderboardAPIHandler_UsesPrecomputedPointsFromUserDoc(t *testing.T) {
	users := repository.NewMemoryUserStore()
	ctx := context.Background()
	users.Upsert(ctx, repository.User{UserID: "u1", Handle: "CrewFan"})
	users.UpdateScores(ctx, "u1", 5, 42, 15, 3)

	lh := newLeaderboard(users)
	req := httptest.NewRequest(http.MethodGet, "/api/leaderboard", nil)
	w := httptest.NewRecorder()
	lh.APIList(w, req)

	body := decodeLeaderboard(t, w)
	if len(body.Entries) != 1 || body.Entries[0].AcesRadioPoints != 42 {
		t.Errorf("expected precomputed 42 AcesRadio points from user doc, got %+v", body.Entries)
	}
}

func TestLeaderboardAPIHandler_ReturnsJSON(t *testing.T) {
	users := repository.NewMemoryUserStore()
	ctx := context.Background()
	users.Upsert(ctx, repository.User{UserID: "u1", Handle: "BlackAndGold@bsky.mock"})
	users.UpdateScores(ctx, "u1", 1, 15, 0, 0)

	lh := newLeaderboard(users)
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
	lh := handlers.NewLeaderboardHandler(repository.NewMemoryPredictionStore(), repository.NewMemoryResultStore(), repository.NewErrorGetAllUserStore(), "Columbus Crew")
	req := httptest.NewRequest(http.MethodGet, "/api/leaderboard", nil)
	w := httptest.NewRecorder()

	lh.APIList(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestLeaderboardAPIHandler_IncludesUpper90Club(t *testing.T) {
	users := repository.NewMemoryUserStore()
	ctx := context.Background()
	users.Upsert(ctx, repository.User{UserID: "u1", Handle: "ColumbusNordecke@bsky.mock"})
	users.UpdateScores(ctx, "u1", 1, 0, 2, 0)

	lh := newLeaderboard(users)
	req := httptest.NewRequest(http.MethodGet, "/api/leaderboard", nil)
	w := httptest.NewRecorder()
	lh.APIList(w, req)

	body := decodeLeaderboard(t, w)
	if len(body.Entries) == 0 || body.Entries[0].Handle != "ColumbusNordecke@bsky.mock" || body.Entries[0].Upper90Points != 2 {
		t.Errorf("unexpected entries: %+v", body.Entries)
	}
}

func TestLeaderboardAPIHandler_UsesHandleFromUserDoc(t *testing.T) {
	users := repository.NewMemoryUserStore()
	ctx := context.Background()
	users.Upsert(ctx, repository.User{UserID: "firebase:abc", Handle: "CrewForever"})
	users.UpdateScores(ctx, "firebase:abc", 1, 15, 0, 0)

	lh := newLeaderboard(users)
	req := httptest.NewRequest(http.MethodGet, "/api/leaderboard", nil)
	w := httptest.NewRecorder()
	lh.APIList(w, req)

	body := decodeLeaderboard(t, w)
	if len(body.Entries) == 0 || body.Entries[0].Handle != "CrewForever" {
		t.Errorf("expected handle CrewForever from UserStore, got %+v", body.Entries)
	}
}

func TestLeaderboardAPIHandler_IncludesUserID(t *testing.T) {
	users := repository.NewMemoryUserStore()
	ctx := context.Background()
	users.Upsert(ctx, repository.User{UserID: "firebase:abc", Handle: "crewfan"})
	users.UpdateScores(ctx, "firebase:abc", 1, 15, 0, 0)

	lh := newLeaderboard(users)
	req := httptest.NewRequest(http.MethodGet, "/api/leaderboard", nil)
	w := httptest.NewRecorder()
	lh.APIList(w, req)

	body := decodeLeaderboard(t, w)
	if len(body.Entries) == 0 || body.Entries[0].UserID != "firebase:abc" {
		t.Errorf("expected userID in leaderboard response, got %+v", body.Entries)
	}
}

func TestLeaderboardAPIHandler_ShowsUsersWithUnscoredPredictionsAtZero(t *testing.T) {
	users := repository.NewMemoryUserStore()
	ctx := context.Background()
	users.Upsert(ctx, repository.User{UserID: "u1", Handle: "EarlyFan"})
	users.UpdateScores(ctx, "u1", 1, 0, 0, 0)
	users.Upsert(ctx, repository.User{UserID: "u2", Handle: "ScoredFan"})
	users.UpdateScores(ctx, "u2", 1, 10, 0, 0)

	lh := newLeaderboard(users)
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
		t.Errorf("expected u2 acesRadioPoints=10, got %d", found["u2"])
	}
}

func TestLeaderboardAPIHandler_HasProfileTrue(t *testing.T) {
	users := repository.NewMemoryUserStore()
	ctx := context.Background()
	users.Upsert(ctx, repository.User{UserID: "firebase:abc", Handle: "known"})
	users.UpdateScores(ctx, "firebase:abc", 1, 0, 0, 0)

	lh := newLeaderboard(users)
	req := httptest.NewRequest(http.MethodGet, "/api/leaderboard", nil)
	w := httptest.NewRecorder()
	lh.APIList(w, req)

	body := decodeLeaderboard(t, w)
	if len(body.Entries) == 0 || !body.Entries[0].HasProfile {
		t.Errorf("expected hasProfile=true for user in UserStore, got %+v", body.Entries)
	}
}

func TestLeaderboardAPIHandler_ExcludesUsersWithNoPredictions(t *testing.T) {
	users := repository.NewMemoryUserStore()
	ctx := context.Background()
	users.Upsert(ctx, repository.User{UserID: "u1", Handle: "NoPredsFan", PredictionCount: 0})

	lh := newLeaderboard(users)
	req := httptest.NewRequest(http.MethodGet, "/api/leaderboard", nil)
	w := httptest.NewRecorder()
	lh.APIList(w, req)

	body := decodeLeaderboard(t, w)
	if len(body.Entries) != 0 {
		t.Errorf("expected no entries for user with PredictionCount=0, got %+v", body.Entries)
	}
}

func TestLeaderboardAPIHandler_IncludesGrouchyPoints(t *testing.T) {
	users := repository.NewMemoryUserStore()
	ctx := context.Background()
	users.Upsert(ctx, repository.User{UserID: "u1", Handle: "GrouchyFan"})
	users.UpdateScores(ctx, "u1", 1, 0, 0, 1)

	lh := newLeaderboard(users)
	req := httptest.NewRequest(http.MethodGet, "/api/leaderboard", nil)
	w := httptest.NewRecorder()
	lh.APIList(w, req)

	body := decodeLeaderboard(t, w)
	if len(body.Entries) == 0 || body.Entries[0].GrouchyPoints != 1 {
		t.Errorf("expected grouchyPoints=1, got %+v", body.Entries)
	}
}
