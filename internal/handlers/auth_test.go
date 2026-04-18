package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mcornell/crew-predictions/internal/handlers"
)

func TestLogout_ClearsSessionCookieAndRedirects(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/auth/logout", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: "somevalue"})
	w := httptest.NewRecorder()

	handlers.Logout(w, req)

	if w.Code != http.StatusFound {
		t.Errorf("expected 302 redirect, got %d", w.Code)
	}
	if w.Header().Get("Location") != "/matches" {
		t.Errorf("expected redirect to /matches, got %s", w.Header().Get("Location"))
	}
	var found bool
	for _, c := range w.Result().Cookies() {
		if c.Name == "session" && c.MaxAge == -1 {
			found = true
		}
	}
	if !found {
		t.Error("expected session cookie to be cleared (MaxAge=-1)")
	}
}
