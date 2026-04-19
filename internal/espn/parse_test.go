package espn

import (
	"strings"
	"testing"
	"time"
)

func TestParseKickoff_RFC3339WithSeconds(t *testing.T) {
	k, err := parseKickoff("2026-05-01T20:00:00Z")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if k.Year() != 2026 || k.Month() != time.May || k.Day() != 1 {
		t.Errorf("unexpected date: %v", k)
	}
}

func TestParseKickoff_ESPNFormatNoSeconds(t *testing.T) {
	k, err := parseKickoff("2026-04-19T23:00Z")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if k.Year() != 2026 || k.Month() != time.April || k.Day() != 19 {
		t.Errorf("unexpected date: %v", k)
	}
}

func TestParseKickoff_InvalidReturnsError(t *testing.T) {
	_, err := parseKickoff("not-a-date")
	if err == nil {
		t.Error("expected error for invalid date")
	}
}

func TestUpcomingURL_ContainsStartDate(t *testing.T) {
	from := time.Date(2026, 4, 19, 0, 0, 0, 0, time.UTC)
	url := upcomingURL(from)
	if !strings.Contains(url, "20260419") {
		t.Errorf("upcomingURL %q missing start date 20260419", url)
	}
}

func TestUpcomingURL_ContainsScoreboard(t *testing.T) {
	from := time.Date(2026, 4, 19, 0, 0, 0, 0, time.UTC)
	url := upcomingURL(from)
	if !strings.Contains(url, "scoreboard") {
		t.Errorf("upcomingURL %q not pointing at scoreboard endpoint", url)
	}
}

func TestDedupeByID_RemovesDuplicate(t *testing.T) {
	from := time.Date(2026, 4, 19, 0, 0, 0, 0, time.UTC)
	matches := dedupeByID([]matchRecord{
		{id: "1", kickoff: from, home: "A", away: "B", status: "STATUS_FULL_TIME"},
		{id: "1", kickoff: from, home: "A", away: "B", status: "STATUS_FULL_TIME"},
		{id: "2", kickoff: from, home: "C", away: "D", status: "STATUS_SCHEDULED"},
	})
	if len(matches) != 2 {
		t.Errorf("expected 2 matches after dedupe, got %d", len(matches))
	}
}
