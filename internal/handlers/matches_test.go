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
