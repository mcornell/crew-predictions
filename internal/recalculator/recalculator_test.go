package recalculator_test

import (
	"context"
	"testing"

	"github.com/mcornell/crew-predictions/internal/recalculator"
	"github.com/mcornell/crew-predictions/internal/repository"
)

func TestRecalculate_SetsAllScoringSystemPoints(t *testing.T) {
	ctx := context.Background()
	users := repository.NewMemoryUserStore()
	preds := repository.NewMemoryPredictionStore()
	results := repository.NewMemoryResultStore()

	users.Upsert(ctx, repository.User{UserID: "u1", Handle: "CrewFan"})
	// Columbus home, exact 2-1 → AcesRadio=15, Upper90=3 (correct outcome + both scores), Grouchy=1 (win by 1)
	preds.Save(ctx, repository.Prediction{MatchID: "m1", UserID: "u1", Handle: "CrewFan", HomeGoals: 2, AwayGoals: 1})
	results.SaveResult(ctx, repository.Result{MatchID: "m1", HomeTeam: "Columbus Crew", AwayTeam: "FC Dallas", HomeGoals: 2, AwayGoals: 1})

	recalculator.Recalculate(ctx, preds, results, users, "Columbus Crew")

	u, _ := users.GetByID(ctx, "u1")
	if u.AcesRadioPoints != 15 {
		t.Errorf("expected 15 AcesRadio points, got %d", u.AcesRadioPoints)
	}
	if u.Upper90Points != 3 {
		t.Errorf("expected 3 Upper90 points, got %d", u.Upper90Points)
	}
	if u.GrouchyPoints != 1 {
		t.Errorf("expected 1 Grouchy point, got %d", u.GrouchyPoints)
	}
}

func TestRecalculate_ZerosPointsForUserWithNoPredictions(t *testing.T) {
	ctx := context.Background()
	users := repository.NewMemoryUserStore()
	// User exists in store but has made no predictions
	users.Upsert(ctx, repository.User{UserID: "u1", Handle: "Lurker", AcesRadioPoints: 99})

	recalculator.Recalculate(ctx, repository.NewMemoryPredictionStore(), repository.NewMemoryResultStore(), users, "Columbus Crew")

	u, _ := users.GetByID(ctx, "u1")
	if u.AcesRadioPoints != 0 {
		t.Errorf("expected 0 points after recalculate with no predictions, got %d", u.AcesRadioPoints)
	}
}

func TestRecalculate_SetsPredictionCount(t *testing.T) {
	ctx := context.Background()
	users := repository.NewMemoryUserStore()
	preds := repository.NewMemoryPredictionStore()

	users.Upsert(ctx, repository.User{UserID: "u1", Handle: "CrewFan"})
	preds.Save(ctx, repository.Prediction{MatchID: "m1", UserID: "u1", HomeGoals: 2, AwayGoals: 1})
	preds.Save(ctx, repository.Prediction{MatchID: "m2", UserID: "u1", HomeGoals: 1, AwayGoals: 0})

	recalculator.Recalculate(ctx, preds, repository.NewMemoryResultStore(), users, "Columbus Crew")

	u, _ := users.GetByID(ctx, "u1")
	if u.PredictionCount != 2 {
		t.Errorf("expected PredictionCount 2, got %d", u.PredictionCount)
	}
}

func TestRecalculate_ReturnsErrorWhenPredictionStoreFails(t *testing.T) {
	ctx := context.Background()
	err := recalculator.Recalculate(ctx, repository.NewErrorGetAllPredictionStore(), repository.NewMemoryResultStore(), repository.NewMemoryUserStore(), "Columbus Crew")
	if err == nil {
		t.Error("expected error when prediction store fails, got nil")
	}
}

func TestRecalculate_ReturnsErrorWhenUserStoreFails(t *testing.T) {
	ctx := context.Background()
	err := recalculator.Recalculate(ctx, repository.NewMemoryPredictionStore(), repository.NewMemoryResultStore(), repository.NewErrorGetAllUserStore(), "Columbus Crew")
	if err == nil {
		t.Error("expected error when user store GetAll fails, got nil")
	}
}

func TestRecalculate_ReturnsErrorWhenUpsertFails(t *testing.T) {
	ctx := context.Background()
	users := repository.NewErrorUpsertUserStore()
	// Seed one prediction so the upsert path is reached — but ErrorUpsertUserStore.GetAll returns empty,
	// so there are no users to upsert. Use a store that has a user but fails on Upsert.
	err := recalculator.Recalculate(ctx, repository.NewMemoryPredictionStore(), repository.NewMemoryResultStore(), repository.NewErrorUpsertWithUserStore(), "Columbus Crew")
	if err == nil {
		t.Errorf("expected error when upsert fails, got nil (users=%v)", users)
	}
}

func TestRecalculate_SetsAcesRadioPoints(t *testing.T) {
	ctx := context.Background()
	users := repository.NewMemoryUserStore()
	preds := repository.NewMemoryPredictionStore()
	results := repository.NewMemoryResultStore()

	users.Upsert(ctx, repository.User{UserID: "u1", Handle: "CrewFan"})
	preds.Save(ctx, repository.Prediction{MatchID: "m1", UserID: "u1", Handle: "CrewFan", HomeGoals: 2, AwayGoals: 1})
	results.SaveResult(ctx, repository.Result{MatchID: "m1", HomeTeam: "Columbus Crew", AwayTeam: "FC Dallas", HomeGoals: 2, AwayGoals: 1})

	if err := recalculator.Recalculate(ctx, preds, results, users, "Columbus Crew"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	u, _ := users.GetByID(ctx, "u1")
	if u.AcesRadioPoints != 15 {
		t.Errorf("expected 15 AcesRadio points (exact score), got %d", u.AcesRadioPoints)
	}
}
