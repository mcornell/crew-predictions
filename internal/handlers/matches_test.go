package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mcornell/crew-predictions/internal/handlers"
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

func TestMatchesHandler_ContainsMatchCards(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/matches", nil)
	w := httptest.NewRecorder()

	handlers.Matches(w, req)

	body := w.Body.String()
	if !strings.Contains(body, "Columbus Crew") {
		t.Errorf("expected body to contain 'Columbus Crew', got: %s", body)
	}
	if !strings.Contains(body, `data-testid="match-card"`) {
		t.Errorf("expected body to contain match-card testid, got: %s", body)
	}
}
