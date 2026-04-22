package handlers_test

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mcornell/crew-predictions/internal/handlers"
)

func cookieWithPayload(payload map[string]string) *http.Cookie {
	data, _ := json.Marshal(payload)
	return &http.Cookie{Name: "__session", Value: base64.StdEncoding.EncodeToString(data)}
}

func TestUserFromSession_ReturnsNilForEmptyUserID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(cookieWithPayload(map[string]string{"userID": "", "handle": "fan@bsky.mock"}))
	if got := handlers.UserFromSession(req); got != nil {
		t.Errorf("expected nil for empty userID, got %+v", got)
	}
}

func TestUserFromSession_ReturnsNilForEmptyHandle(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(cookieWithPayload(map[string]string{"userID": "google:abc", "handle": ""}))
	if got := handlers.UserFromSession(req); got != nil {
		t.Errorf("expected nil for empty handle, got %+v", got)
	}
}

func TestLogout_ClearsSessionCookieAndRedirects(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/auth/logout", nil)
	req.AddCookie(&http.Cookie{Name: "__session", Value: "somevalue"})
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
		if c.Name == "__session" && c.MaxAge == -1 {
			found = true
		}
	}
	if !found {
		t.Error("expected session cookie to be cleared (MaxAge=-1)")
	}
}
