package espn

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/mcornell/crew-predictions/internal/models"
)

// scoreField handles ESPN's polymorphic score value: either an integer (upcoming)
// or an object with displayValue (completed matches).
type scoreField struct {
	Display string
}

func (s *scoreField) UnmarshalJSON(data []byte) error {
	// Completed matches send {"displayValue":"2",...}
	var obj struct {
		DisplayValue string `json:"displayValue"`
	}
	if err := json.Unmarshal(data, &obj); err == nil && obj.DisplayValue != "" {
		s.Display = obj.DisplayValue
		return nil
	}
	// Live matches send a plain string "2"
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		s.Display = str
	}
	return nil
}

// leagueSlugs are the ESPN league identifiers to check for Columbus Crew matches.
// Friendlies are excluded by omission.
var leagueSlugs = []string{
	"usa.1",              // MLS
	"usa.open",           // US Open Cup
	"concacaf.leagues.cup", // Leagues Cup
	"concacaf.champions", // CONCACAF Champions Cup
}

const teamID = "183"
const espnBase = "https://site.api.espn.com/apis/site/v2/sports/soccer"

type espnResponse struct {
	Events []struct {
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
			Venue struct {
				FullName string `json:"fullName"`
			} `json:"venue"`
			Status struct {
				DisplayClock string `json:"displayClock"`
				State        string `json:"state"`
				Type         struct {
					Name string `json:"name"`
				} `json:"type"`
			} `json:"status"`
		} `json:"competitions"`
	} `json:"events"`
}

type matchRecord struct {
	id           string
	kickoff      time.Time
	home         string
	away         string
	homeScore    string
	awayScore    string
	status       string
	state        string
	displayClock string
	venue        string
}

func deriveState(espnState, statusName string) string {
	if espnState != "" {
		return espnState
	}
	switch statusName {
	case "STATUS_FIRST_HALF", "STATUS_SECOND_HALF", "STATUS_HALFTIME",
		"STATUS_IN_PROGRESS", "STATUS_END_PERIOD", "STATUS_OVERTIME",
		"STATUS_EXTRA_TIME", "STATUS_SHOOTOUT":
		return "in"
	case "STATUS_FULL_TIME", "STATUS_FINAL", "STATUS_FT",
		"STATUS_FULL_PEN", "STATUS_ABANDONED":
		return "post"
	default:
		return "pre"
	}
}

func parseKickoff(s string) (time.Time, error) {
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}
	return time.Parse("2006-01-02T15:04Z07:00", s)
}

func scheduleURL(base, league string) string {
	return fmt.Sprintf("%s/%s/teams/%s/schedule", base, league, teamID)
}

func upcomingURL(base, league string, from time.Time) string {
	// Start 2 days back so delayed/rescheduled matches aren't dropped when
	// ESPN still indexes them under their original kickoff date.
	start := from.AddDate(0, 0, -2)
	end := from.AddDate(0, 0, 8)
	return fmt.Sprintf(
		"%s/%s/scoreboard?dates=%s-%s&limit=500",
		base, league,
		start.Format("20060102"),
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
		var home, away, homeScore, awayScore string
		for _, c := range comp.Competitors {
			if c.HomeAway == "home" {
				home = c.Team.DisplayName
				homeScore = c.Score.Display
			} else {
				away = c.Team.DisplayName
				awayScore = c.Score.Display
			}
		}
		kickoff, _ := parseKickoff(event.Date)
		records = append(records, matchRecord{
			id:           event.ID,
			kickoff:      kickoff,
			home:         home,
			away:         away,
			homeScore:    homeScore,
			awayScore:    awayScore,
			status:       comp.Status.Type.Name,
			state:        deriveState(comp.Status.State, comp.Status.Type.Name),
			displayClock: comp.Status.DisplayClock,
			venue:        comp.Venue.FullName,
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
	if resp.StatusCode != http.StatusOK {
		return nil, nil // league may not have data; skip silently
	}
	var data espnResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return parseEvents(data), nil
}

func isCrewMatch(r matchRecord) bool {
	return r.home == "Columbus Crew" || r.away == "Columbus Crew"
}

func fetchCrewMatchesFrom(base string) ([]models.Match, error) {
	now := time.Now().UTC()
	var all []matchRecord

	for _, league := range leagueSlugs {
		past, err := fetchAndParse(scheduleURL(base, league))
		if err != nil {
			return nil, fmt.Errorf("espn schedule fetch (%s): %w", league, err)
		}
		for _, r := range past {
			if isCrewMatch(r) {
				all = append(all, r)
			}
		}

		upcoming, err := fetchAndParse(upcomingURL(base, league, now))
		if err != nil {
			return nil, fmt.Errorf("espn scoreboard fetch (%s): %w", league, err)
		}
		for _, r := range upcoming {
			if isCrewMatch(r) {
				all = append(all, r)
			}
		}
	}

	all = dedupeByID(all)
	sort.Slice(all, func(i, j int) bool {
		return all[i].kickoff.Before(all[j].kickoff)
	})

	matches := make([]models.Match, len(all))
	for i, r := range all {
		matches[i] = models.Match{
			ID:           r.id,
			HomeTeam:     r.home,
			AwayTeam:     r.away,
			Kickoff:      r.kickoff,
			Status:       r.status,
			HomeScore:    r.homeScore,
			AwayScore:    r.awayScore,
			State:        r.state,
			DisplayClock: r.displayClock,
		}
	}
	return matches, nil
}

func FetchCrewMatches() ([]models.Match, error) {
	return fetchCrewMatchesFrom(espnBase)
}
