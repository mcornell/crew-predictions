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
	matches := []models.Match{}
	var buf bytes.Buffer
	err := templates.MatchList(matches).Render(context.Background(), &buf)
	if err != nil {
		t.Fatalf("render error: %v", err)
	}
	if !strings.Contains(buf.String(), "<title>Crew Predictions</title>") {
		t.Errorf("expected title, got: %s", buf.String())
	}
}

func TestMatchList_RendersHeader(t *testing.T) {
	matches := []models.Match{}
	var buf bytes.Buffer
	templates.MatchList(matches).Render(context.Background(), &buf)
	if !strings.Contains(buf.String(), "Crew Predictions") {
		t.Errorf("expected header text, got: %s", buf.String())
	}
}

func TestMatchList_RendersMatchCards(t *testing.T) {
	matches := []models.Match{
		{ID: "1", HomeTeam: "Columbus Crew", AwayTeam: "Atlanta United", Kickoff: time.Now(), Status: "scheduled"},
	}
	var buf bytes.Buffer
	templates.MatchList(matches).Render(context.Background(), &buf)
	body := buf.String()
	if !strings.Contains(body, "Columbus Crew") {
		t.Errorf("expected Columbus Crew in output, got: %s", body)
	}
	if !strings.Contains(body, `data-testid="match-card"`) {
		t.Errorf("expected match-card testid, got: %s", body)
	}
}
