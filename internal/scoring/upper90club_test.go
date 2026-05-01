package scoring_test

import (
	"testing"

	"github.com/mcornell/crew-predictions/internal/scoring"
)

// Columbus is away in all synthetic tests below (columbusIsHome=false).
// Away goals = Columbus goals.

func TestUpper90Club_ExactScore(t *testing.T) {
	// Exact score: correct outcome + correct Crew goals + correct opponent goals → +3
	got := scoring.Upper90Club(scoring.Result{Home: 2, Away: 0}, scoring.Prediction{Home: 2, Away: 0}, false)
	if got != 3 {
		t.Errorf("expected 3 for exact score, got %d", got)
	}
}

func TestUpper90Club_CorrectOutcomeAndOpponentGoals(t *testing.T) {
	// Portland 3-2 Crew, pick Portland 3-0: correct outcome + correct opponent goals, wrong Crew goals → +2
	got := scoring.Upper90Club(scoring.Result{Home: 3, Away: 2}, scoring.Prediction{Home: 3, Away: 0}, false)
	if got != 2 {
		t.Errorf("expected 2, got %d", got)
	}
}

func TestUpper90Club_UserExamples_Portland3Crew2(t *testing.T) {
	result := scoring.Result{Home: 3, Away: 2} // Portland 3-2 Crew, Crew is away
	cases := []struct {
		pred scoring.Prediction
		want int
		desc string
	}{
		{scoring.Prediction{Home: 1, Away: 0}, 1, "Portland 1-0: correct outcome only"},
		{scoring.Prediction{Home: 1, Away: 2}, 1, "Portland 1-2: correct Crew goals only"},
		{scoring.Prediction{Home: 4, Away: 2}, 2, "Portland 4-2: correct outcome + Crew goals"},
		{scoring.Prediction{Home: 3, Away: 2}, 3, "Portland 3-2: exact"},
	}
	for _, c := range cases {
		got := scoring.Upper90Club(result, c.pred, false)
		if got != c.want {
			t.Errorf("%s: expected %d, got %d", c.desc, c.want, got)
		}
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

func TestUpper90Club_ColumbusIsHome_ExactScore(t *testing.T) {
	// Columbus is home, exact 2-1 prediction matches actual 2-1 → +3
	got := scoring.Upper90Club(scoring.Result{Home: 2, Away: 1}, scoring.Prediction{Home: 2, Away: 1}, true)
	if got != 3 {
		t.Errorf("expected 3 for exact score with Columbus at home, got %d", got)
	}
}

func TestUpper90Club_ColumbusIsHome_CorrectOutcomeAndOpponentGoals(t *testing.T) {
	// Columbus is home and wins 3-1; pick 4-1: correct outcome + correct
	// opponent (away) goals, wrong Crew home goals → +2.
	got := scoring.Upper90Club(scoring.Result{Home: 3, Away: 1}, scoring.Prediction{Home: 4, Away: 1}, true)
	if got != 2 {
		t.Errorf("expected 2 for correct outcome + opponent goals (Columbus home), got %d", got)
	}
}

func TestUpper90Club_ColumbusIsHome_CorrectWinnerOnly(t *testing.T) {
	// Columbus is home and wins 3-1; pick 2-0: correct outcome only → +1.
	got := scoring.Upper90Club(scoring.Result{Home: 3, Away: 1}, scoring.Prediction{Home: 2, Away: 0}, true)
	if got != 1 {
		t.Errorf("expected 1 for correct winner only (Columbus home), got %d", got)
	}
}
