package espn

import (
	"encoding/json"
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
	url := upcomingURL(espnBase, "usa.1", from)
	if !strings.Contains(url, "20260419") {
		t.Errorf("upcomingURL %q missing start date 20260419", url)
	}
}

func TestUpcomingURL_ContainsScoreboard(t *testing.T) {
	from := time.Date(2026, 4, 19, 0, 0, 0, 0, time.UTC)
	url := upcomingURL(espnBase, "usa.1", from)
	if !strings.Contains(url, "scoreboard") {
		t.Errorf("upcomingURL %q not pointing at scoreboard endpoint", url)
	}
}

func TestUpcomingURL_EndDate8DaysAhead(t *testing.T) {
	from := time.Date(2026, 4, 19, 0, 0, 0, 0, time.UTC)
	url := upcomingURL(espnBase, "usa.1", from)
	if !strings.Contains(url, "20260427") {
		t.Errorf("upcomingURL %q missing end date 8 days ahead (20260427)", url)
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

func TestLeagueSlugs_ContainsMLS(t *testing.T) {
	found := false
	for _, s := range leagueSlugs {
		if s == "usa.1" {
			found = true
		}
	}
	if !found {
		t.Error("leagueSlugs must include usa.1 (MLS)")
	}
}

func TestLeagueSlugs_ContainsOpenCup(t *testing.T) {
	found := false
	for _, s := range leagueSlugs {
		if s == "usa.open" {
			found = true
		}
	}
	if !found {
		t.Error("leagueSlugs must include usa.open (US Open Cup)")
	}
}

func TestLeagueSlugs_DoesNotContainFriendly(t *testing.T) {
	for _, s := range leagueSlugs {
		if strings.Contains(s, "friendly") {
			t.Errorf("leagueSlugs must not include friendlies, found %q", s)
		}
	}
}

func TestScoreField_ParsesObjectForm(t *testing.T) {
	var s scoreField
	if err := json.Unmarshal([]byte(`{"displayValue":"2","value":2.0}`), &s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Display != "2" {
		t.Errorf("expected display '2', got %q", s.Display)
	}
}

func TestScoreField_ParsesIntegerForm(t *testing.T) {
	var s scoreField
	if err := json.Unmarshal([]byte(`0`), &s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Display != "" {
		t.Errorf("expected empty display for zero integer, got %q", s.Display)
	}
}

func TestParseEvents_SkipsEventWithNoCompetitions(t *testing.T) {
	data := espnResponse{}
	data.Events = append(data.Events, struct {
		ID           string `json:"id"`
		Date         string `json:"date"`
		Competitions []struct {
			Competitors []struct {
				HomeAway string     `json:"homeAway"`
				Score    scoreField `json:"score"`
				Team     struct {
					DisplayName string `json:"displayName"`
				} `json:"team"`
			} `json:"competitors"`
			Status struct {
				State string `json:"state"`
				Type  struct {
					Name string `json:"name"`
				} `json:"type"`
			} `json:"status"`
		} `json:"competitions"`
	}{ID: "e1", Date: "2026-05-01T20:00Z"})

	records := parseEvents(data)
	if len(records) != 0 {
		t.Errorf("expected 0 records for event with no competitions, got %d", len(records))
	}
}
