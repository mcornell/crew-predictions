package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mcornell/crew-predictions/internal/repository"
)

func TestProfileHandler_ReturnsHandleAndLocation(t *testing.T) {
	users := repository.NewMemoryUserStore()
	ctx := context.Background()
	users.Upsert(ctx, repository.User{UserID: "u1", Handle: "CrewFan", Location: "Columbus, OH"})

	h := NewProfileHandler(repository.NewMemoryPredictionStore(), repository.NewMemoryResultStore(), users, "Columbus Crew")
	req := httptest.NewRequest(http.MethodGet, "/api/profile/u1", nil)
	req.SetPathValue("userID", "u1")
	w := httptest.NewRecorder()
	h.Get(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var body struct {
		Handle   string `json:"handle"`
		Location string `json:"location"`
	}
	json.NewDecoder(w.Body).Decode(&body)
	if body.Handle != "CrewFan" {
		t.Errorf("expected handle CrewFan, got %q", body.Handle)
	}
	if body.Location != "Columbus, OH" {
		t.Errorf("expected location Columbus, OH, got %q", body.Location)
	}
}

func TestProfileHandler_ReturnsPredictionCount(t *testing.T) {
	users := repository.NewMemoryUserStore()
	preds := repository.NewMemoryPredictionStore()
	ctx := context.Background()
	users.Upsert(ctx, repository.User{UserID: "u1", Handle: "CrewFan"})
	preds.Save(ctx, repository.Prediction{MatchID: "m1", UserID: "u1", Handle: "CrewFan", HomeGoals: 2, AwayGoals: 0})
	preds.Save(ctx, repository.Prediction{MatchID: "m2", UserID: "u1", Handle: "CrewFan", HomeGoals: 1, AwayGoals: 1})
	preds.Save(ctx, repository.Prediction{MatchID: "m1", UserID: "u2", Handle: "Other", HomeGoals: 0, AwayGoals: 0})

	h := NewProfileHandler(preds, repository.NewMemoryResultStore(), users, "Columbus Crew")
	req := httptest.NewRequest(http.MethodGet, "/api/profile/u1", nil)
	req.SetPathValue("userID", "u1")
	w := httptest.NewRecorder()
	h.Get(w, req)

	var body struct {
		PredictionCount int `json:"predictionCount"`
	}
	json.NewDecoder(w.Body).Decode(&body)
	if body.PredictionCount != 2 {
		t.Errorf("expected predictionCount 2, got %d", body.PredictionCount)
	}
}

func TestProfileHandler_ReturnsLeaderboardStanding(t *testing.T) {
	users := repository.NewMemoryUserStore()
	preds := repository.NewMemoryPredictionStore()
	results := repository.NewMemoryResultStore()
	ctx := context.Background()

	users.Upsert(ctx, repository.User{UserID: "u1", Handle: "TopFan"})
	users.Upsert(ctx, repository.User{UserID: "u2", Handle: "OtherFan"})
	preds.Save(ctx, repository.Prediction{MatchID: "m1", UserID: "u1", Handle: "TopFan", HomeGoals: 2, AwayGoals: 1})
	preds.Save(ctx, repository.Prediction{MatchID: "m1", UserID: "u2", Handle: "OtherFan", HomeGoals: 0, AwayGoals: 0})
	results.SaveResult(ctx, repository.Result{MatchID: "m1", HomeTeam: "Columbus Crew", AwayTeam: "FC Dallas", HomeGoals: 2, AwayGoals: 1})

	h := NewProfileHandler(preds, results, users, "Columbus Crew")
	req := httptest.NewRequest(http.MethodGet, "/api/profile/u1", nil)
	req.SetPathValue("userID", "u1")
	w := httptest.NewRecorder()
	h.Get(w, req)

	var body struct {
		AcesRadio struct {
			Points int `json:"points"`
			Rank   int `json:"rank"`
		} `json:"acesRadio"`
	}
	json.NewDecoder(w.Body).Decode(&body)
	if body.AcesRadio.Points != 15 {
		t.Errorf("expected 15 Aces Radio points, got %d", body.AcesRadio.Points)
	}
	if body.AcesRadio.Rank != 1 {
		t.Errorf("expected rank 1, got %d", body.AcesRadio.Rank)
	}
}

func TestProfileHandler_ReturnsGrouchyStanding(t *testing.T) {
	users := repository.NewMemoryUserStore()
	preds := repository.NewMemoryPredictionStore()
	results := repository.NewMemoryResultStore()
	ctx := context.Background()

	users.Upsert(ctx, repository.User{UserID: "u1", Handle: "GrouchyFan"})
	users.Upsert(ctx, repository.User{UserID: "u2", Handle: "OtherFan"})
	// Columbus home, 2-0 win (Win by 2+). u1 predicts 3-0 (Win by 2+) → 1 pt. u2 predicts 0-1 (Lose by 1) → 0 pt.
	preds.Save(ctx, repository.Prediction{MatchID: "m1", UserID: "u1", HomeGoals: 3, AwayGoals: 0})
	preds.Save(ctx, repository.Prediction{MatchID: "m1", UserID: "u2", HomeGoals: 0, AwayGoals: 1})
	results.SaveResult(ctx, repository.Result{MatchID: "m1", HomeTeam: "Columbus Crew", AwayTeam: "FC Dallas", HomeGoals: 2, AwayGoals: 0})

	h := NewProfileHandler(preds, results, users, "Columbus Crew")
	req := httptest.NewRequest(http.MethodGet, "/api/profile/u1", nil)
	req.SetPathValue("userID", "u1")
	w := httptest.NewRecorder()
	h.Get(w, req)

	var body struct {
		Grouchy struct {
			Points int `json:"points"`
			Rank   int `json:"rank"`
		} `json:"grouchy"`
	}
	json.NewDecoder(w.Body).Decode(&body)
	if body.Grouchy.Points != 1 {
		t.Errorf("expected 1 Grouchy point, got %d", body.Grouchy.Points)
	}
	if body.Grouchy.Rank != 1 {
		t.Errorf("expected rank 1, got %d", body.Grouchy.Rank)
	}
}

func TestProfileHandler_Returns500WhenUserStoreFails(t *testing.T) {
	h := NewProfileHandler(repository.NewMemoryPredictionStore(), repository.NewMemoryResultStore(), repository.NewErrorGetByIDUserStore(), "Columbus Crew")
	req := httptest.NewRequest(http.MethodGet, "/api/profile/u1", nil)
	req.SetPathValue("userID", "u1")
	w := httptest.NewRecorder()
	h.Get(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestProfileHandler_Returns404ForUnknownUser(t *testing.T) {
	h := NewProfileHandler(repository.NewMemoryPredictionStore(), repository.NewMemoryResultStore(), repository.NewMemoryUserStore(), "Columbus Crew")
	req := httptest.NewRequest(http.MethodGet, "/api/profile/nobody", nil)
	req.SetPathValue("userID", "nobody")
	w := httptest.NewRecorder()
	h.Get(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestProfileHandler_ReturnsZeroRankWhenNoScoredPredictions(t *testing.T) {
	users := repository.NewMemoryUserStore()
	ctx := context.Background()
	users.Upsert(ctx, repository.User{UserID: "u1", Handle: "NewFan"})

	h := NewProfileHandler(repository.NewMemoryPredictionStore(), repository.NewMemoryResultStore(), users, "Columbus Crew")
	req := httptest.NewRequest(http.MethodGet, "/api/profile/u1", nil)
	req.SetPathValue("userID", "u1")
	w := httptest.NewRecorder()
	h.Get(w, req)

	var body struct {
		AcesRadio struct {
			Points int `json:"points"`
			Rank   int `json:"rank"`
		} `json:"acesRadio"`
	}
	json.NewDecoder(w.Body).Decode(&body)
	if body.AcesRadio.Rank != 0 {
		t.Errorf("expected rank 0 when no scored predictions, got %d", body.AcesRadio.Rank)
	}
}

