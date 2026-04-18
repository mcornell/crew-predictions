package scoring_test

import (
	"testing"

	"github.com/mcornell/crew-predictions/internal/scoring"
)

func TestUpper90Club_ExactScore(t *testing.T) {
	got := scoring.Upper90Club(scoring.Result{Home: 2, Away: 0}, scoring.Prediction{Home: 2, Away: 0})
	if got != 2 {
		t.Errorf("expected 2 for exact score, got %d", got)
	}
}

func TestUpper90Club_ZeroForAnythingElse(t *testing.T) {
	got := scoring.Upper90Club(scoring.Result{Home: 2, Away: 0}, scoring.Prediction{Home: 3, Away: 1})
	if got != 0 {
		t.Errorf("expected 0 for non-exact prediction, got %d", got)
	}
}
