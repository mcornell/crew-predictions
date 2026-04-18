package scoring_test

import (
	"testing"

	"github.com/mcornell/crew-predictions/internal/scoring"
)

func TestUpper90_ExactScore(t *testing.T) {
	got := scoring.Upper90(scoring.Result{Home: 2, Away: 0}, scoring.Prediction{Home: 2, Away: 0})
	if got != 1 {
		t.Errorf("expected 1 for exact score, got %d", got)
	}
}
