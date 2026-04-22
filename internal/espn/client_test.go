package espn

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func serveESPN(t *testing.T, payload espnResponse) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(payload)
	}))
}

func crewEvent(id, date, status string) espnResponse {
	return espnResponse{Events: []struct {
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
				Type struct {
					Name string `json:"name"`
				} `json:"type"`
			} `json:"status"`
		} `json:"competitions"`
	}{
		{
			ID:   id,
			Date: date,
			Competitions: []struct {
				Competitors []struct {
					HomeAway string     `json:"homeAway"`
					Score    scoreField `json:"score"`
					Team     struct {
						DisplayName string `json:"displayName"`
					} `json:"team"`
				} `json:"competitors"`
				Status struct {
					Type struct {
						Name string `json:"name"`
					} `json:"type"`
				} `json:"status"`
			}{
				{
					Competitors: []struct {
						HomeAway string     `json:"homeAway"`
						Score    scoreField `json:"score"`
						Team     struct {
							DisplayName string `json:"displayName"`
						} `json:"team"`
					}{
						{HomeAway: "home", Team: struct {
							DisplayName string `json:"displayName"`
						}{DisplayName: "Columbus Crew"}},
						{HomeAway: "away", Team: struct {
							DisplayName string `json:"displayName"`
						}{DisplayName: "Portland Timbers"}},
					},
					Status: struct {
						Type struct {
							Name string `json:"name"`
						} `json:"type"`
					}{Type: struct {
						Name string `json:"name"`
					}{Name: status}},
				},
			},
		},
	}}
}

func TestFetchAndParse_ReturnsMatchFromESPN(t *testing.T) {
	srv := serveESPN(t, crewEvent("evt-1", "2026-05-01T23:00Z", "STATUS_SCHEDULED"))
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
	srv := serveESPN(t, espnResponse{})
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
			json.NewEncoder(w).Encode(espnResponse{})
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
			json.NewEncoder(w).Encode(espnResponse{})
		}
	}))
	defer srv.Close()

	_, err := fetchCrewMatchesFrom(srv.URL)
	if err == nil {
		t.Error("expected error when scoreboard fetch returns invalid JSON")
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
			json.NewEncoder(w).Encode(espnResponse{})
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
