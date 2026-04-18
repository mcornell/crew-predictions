package templates_test

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/mcornell/crew-predictions/internal/models"
	"github.com/mcornell/crew-predictions/internal/repository"
	"github.com/mcornell/crew-predictions/templates"
)

var noPredictions = map[string]*repository.Prediction{}

func TestMatchList_RendersTitle(t *testing.T) {
	var buf bytes.Buffer
	templates.MatchList([]models.Match{}, "", noPredictions).Render(context.Background(), &buf)
	if !strings.Contains(buf.String(), "<title>Crew Predictions</title>") {
		t.Errorf("expected title, got: %s", buf.String())
	}
}

func TestMatchList_ShowsSignInWhenLoggedOut(t *testing.T) {
	var buf bytes.Buffer
	templates.MatchList([]models.Match{}, "", noPredictions).Render(context.Background(), &buf)
	body := buf.String()
	if !strings.Contains(body, "Sign in with Google") {
		t.Errorf("expected sign-in link for unauthenticated user")
	}
}

func TestMatchList_ShowsHandleWhenLoggedIn(t *testing.T) {
	var buf bytes.Buffer
	templates.MatchList([]models.Match{}, "BlackAndGold@bsky.mock", noPredictions).Render(context.Background(), &buf)
	body := buf.String()
	if !strings.Contains(body, "BlackAndGold@bsky.mock") {
		t.Errorf("expected handle in output")
	}
	if strings.Contains(body, "Sign in with Google") {
		t.Errorf("sign-in link should be hidden for authenticated user")
	}
}

func TestMatchList_ShowsSigninNudgeWhenLoggedOutWithMatches(t *testing.T) {
	matches := []models.Match{
		{ID: "1", HomeTeam: "Columbus Crew", AwayTeam: "Atlanta United", Kickoff: time.Now(), Status: "scheduled"},
	}
	var buf bytes.Buffer
	templates.MatchList(matches, "", noPredictions).Render(context.Background(), &buf)
	body := buf.String()
	if !strings.Contains(body, "Sign in") {
		t.Errorf("expected sign-in nudge for unauthenticated user with matches")
	}
	if strings.Contains(body, "Lock In") {
		t.Errorf("prediction form should not appear for unauthenticated user")
	}
}

func TestMatchList_ShowsPredictionFormWhenLoggedIn(t *testing.T) {
	matches := []models.Match{
		{ID: "1", HomeTeam: "Columbus Crew", AwayTeam: "Atlanta United", Kickoff: time.Now(), Status: "scheduled"},
	}
	var buf bytes.Buffer
	templates.MatchList(matches, "BlackAndGold@bsky.mock", noPredictions).Render(context.Background(), &buf)
	body := buf.String()
	if !strings.Contains(body, "Lock In") {
		t.Errorf("expected prediction form for authenticated user")
	}
	if strings.Contains(body, "Sign in") {
		t.Errorf("sign-in nudge should not appear for authenticated user")
	}
}

func TestMatchList_PredictionFormHasNameAttributes(t *testing.T) {
	matches := []models.Match{
		{ID: "match1", HomeTeam: "Columbus Crew", AwayTeam: "Atlanta United", Kickoff: time.Now(), Status: "scheduled"},
	}
	var buf bytes.Buffer
	templates.MatchList(matches, "BlackAndGold@bsky.mock", noPredictions).Render(context.Background(), &buf)
	body := buf.String()
	if !strings.Contains(body, `name="home_goals"`) {
		t.Errorf("expected home_goals input name")
	}
	if !strings.Contains(body, `name="away_goals"`) {
		t.Errorf("expected away_goals input name")
	}
	if !strings.Contains(body, `name="match_id"`) {
		t.Errorf("expected match_id hidden input")
	}
}

func TestMatchList_RendersMatchCards(t *testing.T) {
	matches := []models.Match{
		{ID: "1", HomeTeam: "Columbus Crew", AwayTeam: "Atlanta United", Kickoff: time.Now(), Status: "scheduled"},
	}
	var buf bytes.Buffer
	templates.MatchList(matches, "", noPredictions).Render(context.Background(), &buf)
	body := buf.String()
	if !strings.Contains(body, "Columbus Crew") {
		t.Errorf("expected Columbus Crew in output")
	}
	if !strings.Contains(body, `data-testid="match-card"`) {
		t.Errorf("expected match-card testid")
	}
}

func TestMatchList_ShowsSavedScoreWhenPredictionExists(t *testing.T) {
	matches := []models.Match{
		{ID: "match-99", HomeTeam: "Columbus Crew", AwayTeam: "FC Cincinnati", Kickoff: time.Now(), Status: "scheduled"},
	}
	predictions := map[string]*repository.Prediction{
		"match-99": {MatchID: "match-99", Handle: "BlackAndGold@bsky.mock", HomeGoals: 3, AwayGoals: 1},
	}
	var buf bytes.Buffer
	templates.MatchList(matches, "BlackAndGold@bsky.mock", predictions).Render(context.Background(), &buf)
	body := buf.String()
	if !strings.Contains(body, "3") || !strings.Contains(body, "Your Pick") {
		t.Errorf("expected saved score and label, got: %s", body)
	}
	if strings.Contains(body, "Lock In") {
		t.Errorf("prediction form should not appear when prediction already saved")
	}
}
