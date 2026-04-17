package espn_test

import (
	"testing"

	"github.com/mcornell/crew-predictions/internal/espn"
)

func TestFetchCrewMatches_ReturnsMatches(t *testing.T) {
	matches, err := espn.FetchCrewMatches()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(matches) == 0 {
		t.Error("expected at least one match, got none")
	}
}

func TestFetchCrewMatches_MatchesIncludeCrew(t *testing.T) {
	matches, err := espn.FetchCrewMatches()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, m := range matches {
		if m.HomeTeam != "Columbus Crew" && m.AwayTeam != "Columbus Crew" {
			t.Errorf("expected Columbus Crew in match, got %s vs %s", m.HomeTeam, m.AwayTeam)
		}
	}
}
