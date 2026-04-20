package espn

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
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
