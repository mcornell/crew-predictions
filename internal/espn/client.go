package espn

import (
	"encoding/json"
	"fmt"
	"net/http"
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

func FetchCrewMatches() ([]models.Match, error) {
	resp, err := http.Get(scheduleURL)
	if err != nil {
		return nil, fmt.Errorf("espn fetch: %w", err)
	}
	defer resp.Body.Close()

	var data espnResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("espn decode: %w", err)
	}

	var matches []models.Match
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

		kickoff, _ := time.Parse(time.RFC3339, event.Date)
		matches = append(matches, models.Match{
			ID:       event.ID,
			HomeTeam: home,
			AwayTeam: away,
			Kickoff:  kickoff,
			Status:   comp.Status.Type.Name,
		})
	}

	return matches, nil
}
