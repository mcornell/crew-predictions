package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/mcornell/crew-predictions/internal/handlers"
	"github.com/mcornell/crew-predictions/internal/models"
)

func TestMatchesHandler_ReturnsOK(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/matches", nil)
	w := httptest.NewRecorder()

	handlers.Matches(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestMatchesHandler_ContainsHeading(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/matches", nil)
	w := httptest.NewRecorder()

	handlers.Matches(w, req)

	if !strings.Contains(w.Body.String(), "Upcoming Matches") {
		t.Errorf("expected body to contain 'Upcoming Matches', got: %s", w.Body.String())
	}
}

func TestMatchesHandler_RendersMatchCards(t *testing.T) {
	matches := []models.Match{
		{
			ID:       "1",
			HomeTeam: "Columbus Crew",
			AwayTeam: "Atlanta United",
			Kickoff:  time.Now().Add(48 * time.Hour),
			Status:   "scheduled",
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/matches", nil)
	w := httptest.NewRecorder()

	handlers.MatchesWithData(w, req, matches)

	body := w.Body.String()
	if !strings.Contains(body, "Columbus Crew") {
		t.Errorf("expected body to contain 'Columbus Crew', got: %s", body)
	}
	if !strings.Contains(body, `data-testid="match-card"`) {
		t.Errorf("expected body to contain match-card testid, got: %s", body)
	}
}
