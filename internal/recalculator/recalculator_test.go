package recalculator_test

import (
	"context"
	"testing"

	"github.com/mcornell/crew-predictions/internal/recalculator"
	"github.com/mcornell/crew-predictions/internal/repository"
)

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
