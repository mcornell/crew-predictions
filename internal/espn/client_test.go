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

func TestFetchSummaryFrom_ReturnsAttendance(t *testing.T) {
	payload := `{"gameInfo":{"attendance":19903},"keyEvents":[]}`
	srv := serveJSON(t, payload)
	defer srv.Close()

	summary, err := fetchSummaryFrom(srv.URL, "761573")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if summary.Attendance != 19903 {
		t.Errorf("expected Attendance 19903, got %d", summary.Attendance)
	}
}

func TestFetchSummaryFrom_ParsesKeyEvent(t *testing.T) {
	payload := `{"gameInfo":{"attendance":0},"keyEvents":[{"clock":{"displayValue":"15'"},"type":{"type":"goal"},"team":{"displayName":"Columbus Crew"},"participants":[{"athlete":{"displayName":"Cucho Hernandez"}}]}]}`
	srv := serveJSON(t, payload)
	defer srv.Close()

	summary, err := fetchSummaryFrom(srv.URL, "x")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(summary.Events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(summary.Events))
	}
	e := summary.Events[0]
	if e.Clock != "15'" {
		t.Errorf("Clock: got %q, want %q", e.Clock, "15'")
	}
	if e.TypeID != "goal" {
		t.Errorf("TypeID: got %q, want %q", e.TypeID, "goal")
	}
	if e.Team != "Columbus Crew" {
		t.Errorf("Team: got %q, want %q", e.Team, "Columbus Crew")
	}
	if len(e.Players) != 1 || e.Players[0] != "Cucho Hernandez" {
		t.Errorf("Players: got %v, want [Cucho Hernandez]", e.Players)
	}
}

// TestFetchSummaryFrom_RealFixtures parses captured ESPN /summary responses
// for six completed Crew matches that span the variations we observed:
// goal subtypes (header, volley), penalty---scored/saved, own-goal + red-card,
// and a USOC match where ESPN reports no attendance. Re-fetch fixtures with
//
//	curl -s 'https://site.api.espn.com/apis/site/v2/sports/soccer/usa.1/summary?event=<id>' \
//	  | jq '{gameInfo,keyEvents,boxscore,header,rosters,standings,leaders,headToHeadGames,meta,format}' \
//	  > internal/espn/testdata/summary_<id>.json
func TestFetchSummaryFrom_RealFixtures(t *testing.T) {
	cases := []struct {
		matchID            string
		label              string
		expectedAttendance int
		expectedReferee    string // empty when ESPN returned no officials
		minEvents          int
		requiredTypes      []string
	}{
		{"761573", "Philadelphia (own goal + red card)", 19903, "", 19, []string{"goal", "own-goal", "red-card", "yellow-card"}},
		{"761499", "Toronto (header goal)", 15384, "Pierre-Luc Lauziere", 20, []string{"goal---header"}},
		{"761451", "Portland (volley goal)", 22210, "", 21, []string{"goal---volley"}},
		{"761552", "Revolution (penalty scored)", 16257, "Timothy Ford", 18, []string{"penalty---scored"}},
		{"761461", "Sporting KC (penalty saved)", 18522, "Sergii Demianchuk", 19, []string{"penalty---saved"}},
		{"401869714", "USOC vs Knoxville (no attendance)", 0, "Nabil Bensalah", 20, []string{"goal---header", "substitution"}},
	}

	for _, tc := range cases {
		t.Run(tc.label, func(t *testing.T) {
			payload, err := os.ReadFile("testdata/summary_" + tc.matchID + ".json")
			if err != nil {
				t.Fatalf("read fixture: %v", err)
			}
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write(payload)
			}))
			defer srv.Close()

			summary, err := fetchSummaryFrom(srv.URL, tc.matchID)
			if err != nil {
				t.Fatalf("fetchSummaryFrom: %v", err)
			}
			if summary.Attendance != tc.expectedAttendance {
				t.Errorf("Attendance: got %d, want %d", summary.Attendance, tc.expectedAttendance)
			}
			if summary.Referee != tc.expectedReferee {
				t.Errorf("Referee: got %q, want %q", summary.Referee, tc.expectedReferee)
			}
			if !strings.Contains(summary.HomeLogo, "espncdn.com/i/teamlogos/soccer") {
				t.Errorf("HomeLogo: expected ESPN team logo URL, got %q", summary.HomeLogo)
			}
			if !strings.Contains(summary.AwayLogo, "espncdn.com/i/teamlogos/soccer") {
				t.Errorf("AwayLogo: expected ESPN team logo URL, got %q", summary.AwayLogo)
			}
			if summary.HomeLogo == summary.AwayLogo {
				t.Errorf("HomeLogo and AwayLogo are identical: %q", summary.HomeLogo)
			}
			if len(summary.Events) < tc.minEvents {
				t.Errorf("Events: got %d, want at least %d", len(summary.Events), tc.minEvents)
			}
			seen := map[string]bool{}
			for _, e := range summary.Events {
				seen[e.TypeID] = true
			}
			for _, want := range tc.requiredTypes {
				if !seen[want] {
					t.Errorf("expected event type %q not found; got types %v", want, mapKeys(seen))
				}
			}
		})
	}
}

func mapKeys(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

func TestFetchSummaryFrom_ReturnsErrorOnNetworkFailure(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	srv.Close()

	_, err := fetchSummaryFrom(srv.URL, "x")
	if err == nil {
		t.Error("expected error on network failure, got nil")
	}
}

func TestFetchSummaryFrom_ReturnsEmptySummaryOnNonOK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	summary, err := fetchSummaryFrom(srv.URL, "x")
	if err != nil {
		t.Fatalf("expected nil error on non-OK, got %v", err)
	}
	if summary.Attendance != 0 {
		t.Errorf("expected zero attendance on non-OK, got %d", summary.Attendance)
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
