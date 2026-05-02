package espn

import (
	"testing"
)

func TestFixtureFetcher_ReturnsParsedSummaryForKnownMatch(t *testing.T) {
	fetch := FixtureFetcher("testdata")
	summary, err := fetch("761573")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if summary.Attendance != 19903 {
		t.Errorf("expected attendance 19903 (Philadelphia fixture), got %d", summary.Attendance)
	}
	if summary.HomeLogo == "" || summary.AwayLogo == "" {
		t.Errorf("expected both logos populated, got home=%q away=%q", summary.HomeLogo, summary.AwayLogo)
	}
	if len(summary.Events) == 0 {
		t.Error("expected events from the fixture, got empty slice")
	}
}

func TestFixtureFetcher_ReturnsRefereeWhenFixtureHasOne(t *testing.T) {
	fetch := FixtureFetcher("testdata")
	summary, err := fetch("761499")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if summary.Referee != "Pierre-Luc Lauziere" {
		t.Errorf("expected referee Pierre-Luc Lauziere, got %q", summary.Referee)
	}
}

func TestFixtureFetcher_ReturnsEmptyForUnknownMatchID(t *testing.T) {
	fetch := FixtureFetcher("testdata")
	summary, err := fetch("does-not-exist")
	if err != nil {
		t.Fatalf("expected nil error for missing fixture, got %v", err)
	}
	if summary.Attendance != 0 || len(summary.Events) != 0 {
		t.Errorf("expected empty MatchSummary for unknown match, got %+v", summary)
	}
}

func TestFixtureFetcher_ReturnsErrorOnInvalidDirectory(t *testing.T) {
	fetch := FixtureFetcher("/nonexistent-fixture-dir")
	// Missing files inside a missing directory are still "not exist" — graceful.
	summary, err := fetch("761573")
	if err != nil {
		t.Fatalf("expected nil error when fixture missing, got %v", err)
	}
	if summary.Attendance != 0 {
		t.Errorf("expected empty summary, got %+v", summary)
	}
}
