package espn

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/mcornell/crew-predictions/internal/models"
)

const scheduleURL = "https://site.api.espn.com/apis/site/v2/sports/soccer/usa.1/teams/183/schedule"

type espnResponse struct {
	Events []struct {
		ID           string `json:"id"`
		Date         string `json:"date"`
		Competitions []struct {
			Competitors []struct {
				HomeAway string `json:"homeAway"`
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
	} `json:"events"`
}

type matchRecord struct {
	id      string
	kickoff time.Time
	home    string
	away    string
	status  string
}

func parseKickoff(s string) (time.Time, error) {
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}
	return time.Parse("2006-01-02T15:04Z07:00", s)
}

func upcomingURL(from time.Time) string {
	end := from.AddDate(1, 0, 0)
	return fmt.Sprintf(
		"https://site.api.espn.com/apis/site/v2/sports/soccer/usa.1/scoreboard?dates=%s-%s&limit=500",
		from.Format("20060102"),
		end.Format("20060102"),
	)
}

func dedupeByID(records []matchRecord) []matchRecord {
	seen := map[string]bool{}
	out := make([]matchRecord, 0, len(records))
	for _, r := range records {
		if !seen[r.id] {
			seen[r.id] = true
			out = append(out, r)
		}
	}
	return out
}

func parseEvents(data espnResponse) []matchRecord {
	var records []matchRecord
	for _, event := range data.Events {
		if len(event.Competitions) == 0 {
			continue
		}
		comp := event.Competitions[0]
		var home, away string
		for _, c := range comp.Competitors {
			if c.HomeAway == "home" {
				home = c.Team.DisplayName
			} else {
				away = c.Team.DisplayName
			}
		}
		kickoff, _ := parseKickoff(event.Date)
		records = append(records, matchRecord{
			id:      event.ID,
			kickoff: kickoff,
			home:    home,
			away:    away,
			status:  comp.Status.Type.Name,
		})
	}
	return records
}

func fetchAndParse(url string) ([]matchRecord, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var data espnResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return parseEvents(data), nil
}

func FetchCrewMatches() ([]models.Match, error) {
	past, err := fetchAndParse(scheduleURL)
	if err != nil {
		return nil, fmt.Errorf("espn schedule fetch: %w", err)
	}

	upcoming, err := fetchAndParse(upcomingURL(time.Now().UTC()))
	if err != nil {
		return nil, fmt.Errorf("espn scoreboard fetch: %w", err)
	}

	// Filter scoreboard to Columbus Crew matches only
	crewUpcoming := make([]matchRecord, 0, len(upcoming))
	for _, r := range upcoming {
		if r.home == "Columbus Crew" || r.away == "Columbus Crew" {
			crewUpcoming = append(crewUpcoming, r)
		}
	}

	all := dedupeByID(append(past, crewUpcoming...))
	sort.Slice(all, func(i, j int) bool {
		return all[i].kickoff.Before(all[j].kickoff)
	})

	matches := make([]models.Match, len(all))
	for i, r := range all {
		matches[i] = models.Match{
			ID:       r.id,
			HomeTeam: r.home,
			AwayTeam: r.away,
			Kickoff:  r.kickoff,
			Status:   r.status,
		}
	}
	return matches, nil
}
