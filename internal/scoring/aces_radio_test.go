package scoring_test

import (
	"testing"

	"github.com/mcornell/crew-predictions/internal/scoring"
)

func TestAcesRadio_ExactScore(t *testing.T) {
	got := scoring.AcesRadio(scoring.Result{Home: 2, Away: 0}, scoring.Prediction{Home: 2, Away: 0})
	if got != 15 {
		t.Errorf("expected 15 for exact score, got %d", got)
	}
}
