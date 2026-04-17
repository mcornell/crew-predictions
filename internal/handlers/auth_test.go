package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mcornell/crew-predictions/internal/handlers"
)

func TestLoginHandler_RedirectsToGoogle(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/auth/login", nil)
	w := httptest.NewRecorder()

	handlers.Login(w, req)

	if w.Code != http.StatusFound {
		t.Errorf("expected 302 redirect, got %d", w.Code)
	}
	loc := w.Header().Get("Location")
	if loc == "" {
		t.Error("expected Location header, got none")
	}
	// Must redirect to Google's OAuth endpoint
	if len(loc) < 4 || loc[:4] != "http" {
		t.Errorf("expected redirect URL, got: %s", loc)
	}
}
