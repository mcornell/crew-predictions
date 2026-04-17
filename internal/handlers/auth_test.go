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
	if len(loc) < 4 || loc[:4] != "http" {
		t.Errorf("expected redirect URL, got: %s", loc)
	}
}

func TestCallbackHandler_RejectsMissingState(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/auth/callback?code=abc", nil)
	w := httptest.NewRecorder()

	handlers.Callback(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCallbackHandler_RejectsMismatchedState(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/auth/callback?code=abc&state=wrong", nil)
	req.AddCookie(&http.Cookie{Name: "oauth_state", Value: "correct"})
	w := httptest.NewRecorder()

	handlers.Callback(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}
