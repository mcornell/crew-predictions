package espn

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func serveJSON(t *testing.T, payload string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(payload))
	}))
}

func crewEventJSON(id, date, status string) string {
	return `{"events":[{"id":"` + id + `","date":"` + date + `","competitions":[{"competitors":[{"homeAway":"home","score":{},"team":{"displayName":"Columbus Crew"}},{"homeAway":"away","score":{},"team":{"displayName":"Portland Timbers"}}],"status":{"displayClock":"","state":"","type":{"name":"` + status + `"}}}]}]}`
}

func TestFetchAndParse_ReturnsMatchFromESPN(t *testing.T) {
	srv := serveJSON(t, crewEventJSON("evt-1", "2026-05-01T23:00Z", "STATUS_SCHEDULED"))
	defer srv.Close()

	records, err := fetchAndParse(srv.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}
	if records[0].id != "evt-1" {
		t.Errorf("expected id evt-1, got %q", records[0].id)
	}
	if records[0].home != "Columbus Crew" {
		t.Errorf("expected Columbus Crew at home, got %q", records[0].home)
	}
}

func TestFetchAndParse_ReturnsEmptyOnEmptyEvents(t *testing.T) {
	srv := serveJSON(t, `{"events":[]}`)
	defer srv.Close()

	records, err := fetchAndParse(srv.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(records) != 0 {
		t.Errorf("expected 0 records, got %d", len(records))
	}
}

func TestFetchAndParse_Returns404Silently(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	records, err := fetchAndParse(srv.URL)
	if err != nil {
		t.Fatalf("expected nil error on 404, got %v", err)
	}
	if records != nil {
		t.Errorf("expected nil records on 404, got %v", records)
	}
}

func TestIsCrewMatch_TrueWhenHome(t *testing.T) {
	if !isCrewMatch(matchRecord{home: "Columbus Crew", away: "Portland Timbers"}) {
		t.Error("expected true when Columbus Crew is home")
	}
}

func TestIsCrewMatch_TrueWhenAway(t *testing.T) {
	if !isCrewMatch(matchRecord{home: "Portland Timbers", away: "Columbus Crew"}) {
		t.Error("expected true when Columbus Crew is away")
	}
}

func TestIsCrewMatch_FalseWhenNeither(t *testing.T) {
	if isCrewMatch(matchRecord{home: "Portland Timbers", away: "Atlanta United"}) {
		t.Error("expected false when Columbus Crew is not in the match")
	}
}

func TestFetchAndParse_ReturnsErrorOnNetworkFailure(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	srv.Close()

	_, err := fetchAndParse(srv.URL)
	if err == nil {
		t.Error("expected error on network failure, got nil")
	}
}

func TestFetchAndParse_ReturnsErrorOnInvalidJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("not json"))
	}))
	defer srv.Close()

	_, err := fetchAndParse(srv.URL)
	if err == nil {
		t.Error("expected error on invalid JSON, got nil")
	}
}

func TestFetchCrewMatchesFrom_ReturnsErrorWhenScheduleFetchFails(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "schedule") {
			w.Write([]byte("not json"))
		} else {
			w.Write([]byte(`{"events":[]}`))

		}
	}))
	defer srv.Close()

	_, err := fetchCrewMatchesFrom(srv.URL)
	if err == nil {
		t.Error("expected error when schedule fetch returns invalid JSON")
	}
}

func TestFetchCrewMatchesFrom_ReturnsErrorWhenScoreboardFetchFails(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "scoreboard") {
			w.Write([]byte("not json"))
		} else {
			w.Write([]byte(`{"events":[]}`))

		}
	}))
	defer srv.Close()

	_, err := fetchCrewMatchesFrom(srv.URL)
	if err == nil {
		t.Error("expected error when scoreboard fetch returns invalid JSON")
	}
}

func TestFetchCrewMatchesFrom_DerivesStateFromStatusWhenMissing(t *testing.T) {
	cases := []struct {
		statusName    string
		expectedState string
	}{
		{"STATUS_FIRST_HALF", "in"},
		{"STATUS_SECOND_HALF", "in"},
		{"STATUS_HALFTIME", "in"},
		{"STATUS_IN_PROGRESS", "in"},
		{"STATUS_FULL_TIME", "post"},
		{"STATUS_FINAL", "post"},
		{"STATUS_SCHEDULED", "pre"},
	}
	for _, tc := range cases {
		json := `{"events":[{"id":"x","date":"2026-05-01T23:00Z","competitions":[{"competitors":[{"homeAway":"home","score":{},"team":{"displayName":"Columbus Crew"}},{"homeAway":"away","score":{},"team":{"displayName":"Portland Timbers"}}],"status":{"state":"","type":{"name":"` + tc.statusName + `"}}}]}]}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte(json)) }))
		matches, err := fetchCrewMatchesFrom(srv.URL)
		srv.Close()
		if err != nil {
			t.Fatalf("%s: unexpected error: %v", tc.statusName, err)
		}
		if len(matches) == 0 {
			t.Fatalf("%s: expected match", tc.statusName)
		}
		if matches[0].State != tc.expectedState {
			t.Errorf("%s: expected state=%q, got %q", tc.statusName, tc.expectedState, matches[0].State)
		}
	}
}

func TestFetchCrewMatchesFrom_PopulatesMatchState(t *testing.T) {
	liveJSON := `{"events":[{"id":"live-1","date":"2026-05-01T23:00Z","competitions":[{"competitors":[{"homeAway":"home","score":{},"team":{"displayName":"Columbus Crew"}},{"homeAway":"away","score":{},"team":{"displayName":"Portland Timbers"}}],"status":{"state":"in","type":{"name":"STATUS_IN_PROGRESS"}}}]}]}`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(liveJSON))
	}))
	defer srv.Close()

	matches, err := fetchCrewMatchesFrom(srv.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(matches) == 0 {
		t.Fatal("expected at least one match")
	}
	if matches[0].State != "in" {
		t.Errorf("expected state=in, got %q", matches[0].State)
	}
}

func TestFetchCrewMatchesFrom_CapturesDisplayClock(t *testing.T) {
	liveJSON := `{"events":[{"id":"c1","date":"2026-05-01T23:00Z","competitions":[{"competitors":[{"homeAway":"home","score":"2","team":{"displayName":"Columbus Crew"}},{"homeAway":"away","score":"0","team":{"displayName":"Portland Timbers"}}],"status":{"displayClock":"48'","state":"in","type":{"name":"STATUS_SECOND_HALF"}}}]}]}`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(liveJSON))
	}))
	defer srv.Close()

	matches, err := fetchCrewMatchesFrom(srv.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(matches) == 0 {
		t.Fatal("expected at least one match")
	}
	if matches[0].DisplayClock != "48'" {
		t.Errorf("expected DisplayClock '48\\'', got %q", matches[0].DisplayClock)
	}
}

func TestFetchCrewMatchesFrom_HalftimeDisplayClock(t *testing.T) {
	htJSON := `{"events":[{"id":"ht1","date":"2026-05-01T23:00Z","competitions":[{"competitors":[{"homeAway":"home","score":"1","team":{"displayName":"Columbus Crew"}},{"homeAway":"away","score":"0","team":{"displayName":"Portland Timbers"}}],"status":{"displayClock":"HT","state":"in","type":{"name":"STATUS_HALFTIME"}}}]}]}`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(htJSON))
	}))
	defer srv.Close()

	matches, err := fetchCrewMatchesFrom(srv.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(matches) == 0 {
		t.Fatal("expected at least one match")
	}
	if matches[0].DisplayClock != "HT" {
		t.Errorf("expected DisplayClock 'HT', got %q", matches[0].DisplayClock)
	}
}

func TestFetchCrewMatchesFrom_PopulatesVenue(t *testing.T) {
	payload := `{"events":[{"id":"v1","date":"2026-05-01T23:00Z","competitions":[{"competitors":[{"homeAway":"home","score":{},"team":{"displayName":"Columbus Crew"}},{"homeAway":"away","score":{},"team":{"displayName":"FC Dallas"}}],"venue":{"fullName":"ScottsMiracle-Gro Field"},"status":{"state":"pre","type":{"name":"STATUS_SCHEDULED"}}}]}]}`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(payload))
	}))
	defer srv.Close()

	matches, err := fetchCrewMatchesFrom(srv.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(matches) == 0 {
		t.Fatal("expected at least one match")
	}
	if matches[0].Venue != "ScottsMiracle-Gro Field" {
		t.Errorf("expected Venue 'ScottsMiracle-Gro Field', got %q", matches[0].Venue)
	}
}

func TestFetchCrewMatchesFrom_PopulatesRecordsAndForm(t *testing.T) {
	payload := `{"events":[{"id":"r1","date":"2026-05-01T23:00Z","competitions":[{"competitors":[{"homeAway":"home","score":{},"team":{"displayName":"Columbus Crew"},"records":[{"type":"total","summary":"5-3-2"}],"form":"WWWLL"},{"homeAway":"away","score":{},"team":{"displayName":"FC Dallas"},"records":[{"type":"total","summary":"4-4-2"}],"form":"LWDWL"}],"venue":{},"status":{"state":"pre","type":{"name":"STATUS_SCHEDULED"}}}]}]}`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(payload))
	}))
	defer srv.Close()

	matches, err := fetchCrewMatchesFrom(srv.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(matches) == 0 {
		t.Fatal("expected at least one match")
	}
	m := matches[0]
	if m.HomeRecord != "5-3-2" {
		t.Errorf("HomeRecord: got %q, want %q", m.HomeRecord, "5-3-2")
	}
	if m.AwayRecord != "4-4-2" {
		t.Errorf("AwayRecord: got %q, want %q", m.AwayRecord, "4-4-2")
	}
	if m.HomeForm != "WWWLL" {
		t.Errorf("HomeForm: got %q, want %q", m.HomeForm, "WWWLL")
	}
	if m.AwayForm != "LWDWL" {
		t.Errorf("AwayForm: got %q, want %q", m.AwayForm, "LWDWL")
	}
}

func TestFetchCrewMatchesFrom_ReturnsCrewMatchesFromFixtures(t *testing.T) {
	schedule, err := os.ReadFile("testdata/mls_schedule.json")
	if err != nil {
		t.Fatalf("reading schedule fixture: %v", err)
	}
	scoreboard, err := os.ReadFile("testdata/mls_scoreboard.json")
	if err != nil {
		t.Fatalf("reading scoreboard fixture: %v", err)
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "schedule"):
			w.Write(schedule)
		case strings.Contains(r.URL.Path, "scoreboard"):
			w.Write(scoreboard)
		default:
			w.Write([]byte(`{"events":[]}`))

		}
	}))
	defer srv.Close()

	matches, err := fetchCrewMatchesFrom(srv.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(matches) == 0 {
		t.Fatal("expected at least one match from fixtures")
	}
	for _, m := range matches {
		if m.HomeTeam != "Columbus Crew" && m.AwayTeam != "Columbus Crew" {
			t.Errorf("non-Crew match returned: %q vs %q", m.HomeTeam, m.AwayTeam)
		}
	}
}
