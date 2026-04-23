package scoring_test

import (
	"testing"

	"github.com/mcornell/crew-predictions/internal/scoring"
)

func TestGrouchy_ExactWinBy2Plus(t *testing.T) {
	got := scoring.Grouchy(scoring.Result{Home: 2, Away: 0}, scoring.Prediction{Home: 3, Away: 0}, true)
	if got != 1 {
		t.Errorf("expected 1, got %d", got)
	}
}

func TestGrouchy_WrongCategory_WinBy2VsWinBy1(t *testing.T) {
	got := scoring.Grouchy(scoring.Result{Home: 2, Away: 0}, scoring.Prediction{Home: 1, Away: 0}, true)
	if got != 0 {
		t.Errorf("expected 0, got %d", got)
	}
}

func TestGrouchy_WinBy3PlusStillWinBy2Category(t *testing.T) {
	got := scoring.Grouchy(scoring.Result{Home: 3, Away: 0}, scoring.Prediction{Home: 4, Away: 0}, true)
	if got != 1 {
		t.Errorf("expected 1, got %d", got)
	}
}

func TestGrouchy_Draw(t *testing.T) {
	got := scoring.Grouchy(scoring.Result{Home: 1, Away: 1}, scoring.Prediction{Home: 2, Away: 2}, true)
	if got != 1 {
		t.Errorf("expected 1, got %d", got)
	}
}

func TestGrouchy_LoseBy1(t *testing.T) {
	got := scoring.Grouchy(scoring.Result{Home: 0, Away: 1}, scoring.Prediction{Home: 1, Away: 2}, true)
	if got != 1 {
		t.Errorf("expected 1, got %d", got)
	}
}

func TestGrouchy_LoseBy2PlusBoundary(t *testing.T) {
	got := scoring.Grouchy(scoring.Result{Home: 0, Away: 2}, scoring.Prediction{Home: 0, Away: 5}, true)
	if got != 1 {
		t.Errorf("expected 1, got %d", got)
	}
}

func TestGrouchy_AwayTeam_CorrectCategory(t *testing.T) {
	// Columbus away: home=2 away=0 → Columbus margin = away-home = -2 → Lose by 2+
	// predicted home=3 away=1 → Columbus margin = 1-3 = -2 → Lose by 2+
	got := scoring.Grouchy(scoring.Result{Home: 2, Away: 0}, scoring.Prediction{Home: 3, Away: 1}, false)
	if got != 1 {
		t.Errorf("expected 1, got %d", got)
	}
}

func TestGrouchy_AwayTeam_WrongCategory(t *testing.T) {
	// Columbus away: actual home=1 away=0 → margin = 0-1 = -1 → Lose by 1
	// predicted home=0 away=2 → margin = 2-0 = +2 → Win by 2+
	got := scoring.Grouchy(scoring.Result{Home: 1, Away: 0}, scoring.Prediction{Home: 0, Away: 2}, false)
	if got != 0 {
		t.Errorf("expected 0, got %d", got)
	}
}

func TestGrouchy_ExactScore_CorrectCategory(t *testing.T) {
	got := scoring.Grouchy(scoring.Result{Home: 2, Away: 1}, scoring.Prediction{Home: 2, Away: 1}, true)
	if got != 1 {
		t.Errorf("expected 1, got %d", got)
	}
}
