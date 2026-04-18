package scoring_test

import (
	"testing"

	"github.com/mcornell/crew-predictions/internal/scoring"
)

// Real match data from the 2026 Columbus Crew season.
// Match 1: Portland (H) vs Columbus (A) — final 3-2 Portland
// Match 2: Kansas City (H) vs Columbus (A) — final 2-2 draw
// Nobody predicted either score exactly, so all Upper90Club totals are 0.

func TestUpper90Club_RealData_Darby(t *testing.T) {
	pts := scoring.Upper90Club(match1Result, scoring.Prediction{Home: 3, Away: 1})
	pts += scoring.Upper90Club(match2Result, scoring.Prediction{Home: 1, Away: 2})
	if pts != 0 {
		t.Errorf("Darby: expected 0 points, got %d", pts)
	}
}

func TestUpper90Club_RealData_Alex(t *testing.T) {
	pts := scoring.Upper90Club(match1Result, scoring.Prediction{Home: 3, Away: 0})
	pts += scoring.Upper90Club(match2Result, scoring.Prediction{Home: 0, Away: 2})
	if pts != 0 {
		t.Errorf("Alex: expected 0 points, got %d", pts)
	}
}

func TestUpper90Club_RealData_Zidar(t *testing.T) {
	pts := scoring.Upper90Club(match1Result, scoring.Prediction{Home: 2, Away: 3})
	pts += scoring.Upper90Club(match2Result, scoring.Prediction{Home: 1, Away: 3})
	if pts != 0 {
		t.Errorf("Zidar: expected 0 points, got %d", pts)
	}
}

func TestUpper90Club_RealData_Lumbus(t *testing.T) {
	pts := scoring.Upper90Club(match1Result, scoring.Prediction{Home: 2, Away: 4})
	pts += scoring.Upper90Club(match2Result, scoring.Prediction{Home: 1, Away: 3})
	if pts != 0 {
		t.Errorf("Lumbus: expected 0 points, got %d", pts)
	}
}

func TestUpper90Club_RealData_Morgan(t *testing.T) {
	pts := scoring.Upper90Club(match1Result, scoring.Prediction{Home: 2, Away: 0})
	pts += scoring.Upper90Club(match2Result, scoring.Prediction{Home: 0, Away: 3})
	if pts != 0 {
		t.Errorf("Morgan: expected 0 points, got %d", pts)
	}
}
