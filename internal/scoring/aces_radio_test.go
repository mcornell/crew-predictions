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

func TestAcesRadio_CorrectWinner(t *testing.T) {
	got := scoring.AcesRadio(scoring.Result{Home: 2, Away: 0}, scoring.Prediction{Home: 3, Away: 1})
	if got != 10 {
		t.Errorf("expected 10 for correct winner, got %d", got)
	}
}

func TestAcesRadio_FlippedSameScore(t *testing.T) {
	// Predict Crew 3-2 Portland, actual Portland 3-2 Crew
	got := scoring.AcesRadio(scoring.Result{Home: 3, Away: 2}, scoring.Prediction{Home: 2, Away: 3})
	if got != -15 {
		t.Errorf("expected -15 for flipped same scoreline, got %d", got)
	}
}

func TestAcesRadio_ZeroForOtherPredictions(t *testing.T) {
	got := scoring.AcesRadio(scoring.Result{Home: 2, Away: 0}, scoring.Prediction{Home: 0, Away: 1})
	if got != 0 {
		t.Errorf("expected 0 for wrong winner wrong score, got %d", got)
	}
}
