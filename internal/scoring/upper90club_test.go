package scoring_test

import (
	"testing"

	"github.com/mcornell/crew-predictions/internal/scoring"
)

// Columbus is away in all synthetic tests below (columbusIsHome=false).
// Away goals = Columbus goals.

func TestUpper90Club_ExactScore(t *testing.T) {
	// Exact score implies both correct winner and correct Columbus goals → +2
	got := scoring.Upper90Club(scoring.Result{Home: 2, Away: 0}, scoring.Prediction{Home: 2, Away: 0}, false)
	if got != 2 {
		t.Errorf("expected 2 for exact score, got %d", got)
	}
}

func TestUpper90Club_CorrectWinner(t *testing.T) {
	// Correct winner, wrong Columbus goals → +1
	got := scoring.Upper90Club(scoring.Result{Home: 2, Away: 0}, scoring.Prediction{Home: 3, Away: 1}, false)
	if got != 1 {
		t.Errorf("expected 1 for correct winner only, got %d", got)
	}
}

func TestUpper90Club_CorrectColumbusGoals(t *testing.T) {
	// Wrong winner, correct Columbus goals → +1
	got := scoring.Upper90Club(scoring.Result{Home: 2, Away: 1}, scoring.Prediction{Home: 0, Away: 1}, false)
	if got != 1 {
		t.Errorf("expected 1 for correct Columbus goals only, got %d", got)
	}
}

func TestUpper90Club_CorrectWinnerAndColumbusGoals(t *testing.T) {
	// Correct winner (home) + correct Columbus away goals (2) but not exact score (4≠3) → +2
	got := scoring.Upper90Club(scoring.Result{Home: 3, Away: 2}, scoring.Prediction{Home: 4, Away: 2}, false)
	if got != 2 {
		t.Errorf("expected 2 for correct winner and Columbus goals, got %d", got)
	}
}

func TestUpper90Club_ZeroForAnythingElse(t *testing.T) {
	// Wrong winner, wrong Columbus goals → 0
	got := scoring.Upper90Club(scoring.Result{Home: 2, Away: 0}, scoring.Prediction{Home: 0, Away: 1}, false)
	if got != 0 {
		t.Errorf("expected 0 for wrong winner and wrong Columbus goals, got %d", got)
	}
}

func TestUpper90Club_ColumbusIsHome(t *testing.T) {
	// Columbus is home, correct home goals → +1 (wrong winner)
	got := scoring.Upper90Club(scoring.Result{Home: 2, Away: 3}, scoring.Prediction{Home: 2, Away: 0}, true)
	if got != 1 {
		t.Errorf("expected 1 for correct Columbus home goals, got %d", got)
	}
}
