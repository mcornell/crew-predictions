package templates_test

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/mcornell/crew-predictions/internal/models"
	"github.com/mcornell/crew-predictions/templates"
)

func TestMatchList_RendersTitle(t *testing.T) {
	var buf bytes.Buffer
	templates.MatchList([]models.Match{}, "").Render(context.Background(), &buf)
	if !strings.Contains(buf.String(), "<title>Crew Predictions</title>") {
		t.Errorf("expected title, got: %s", buf.String())
	}
}

func TestMatchList_ShowsSignInWhenLoggedOut(t *testing.T) {
	var buf bytes.Buffer
	templates.MatchList([]models.Match{}, "").Render(context.Background(), &buf)
	body := buf.String()
	if !strings.Contains(body, "Sign in with Google") {
		t.Errorf("expected sign-in link for unauthenticated user")
	}
}

func TestMatchList_ShowsHandleWhenLoggedIn(t *testing.T) {
	var buf bytes.Buffer
	templates.MatchList([]models.Match{}, "BlackYellow@bsky.social").Render(context.Background(), &buf)
	body := buf.String()
	if !strings.Contains(body, "BlackYellow@bsky.social") {
		t.Errorf("expected handle in output")
	}
	if strings.Contains(body, "Sign in with Google") {
		t.Errorf("sign-in link should be hidden for authenticated user")
	}
}

func TestMatchList_RendersMatchCards(t *testing.T) {
	matches := []models.Match{
		{ID: "1", HomeTeam: "Columbus Crew", AwayTeam: "Atlanta United", Kickoff: time.Now(), Status: "scheduled"},
	}
	var buf bytes.Buffer
	templates.MatchList(matches, "").Render(context.Background(), &buf)
	body := buf.String()
	if !strings.Contains(body, "Columbus Crew") {
		t.Errorf("expected Columbus Crew in output")
	}
	if !strings.Contains(body, `data-testid="match-card"`) {
		t.Errorf("expected match-card testid")
	}
}
