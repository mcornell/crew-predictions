package handlers

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserFromSession_NoCookie(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	if got := userFromSession(req); got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestUserFromSession_InvalidBase64(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: "not-valid-base64!!!"})
	if got := userFromSession(req); got != "" {
		t.Errorf("expected empty string for bad base64, got %q", got)
	}
}

func TestUserFromSession_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: base64.StdEncoding.EncodeToString([]byte("not json"))})
	if got := userFromSession(req); got != "" {
		t.Errorf("expected empty string for bad JSON, got %q", got)
	}
}

func TestUserFromSession_ReturnsHandle(t *testing.T) {
	data, _ := json.Marshal(map[string]string{"handle": "BlackYellow@bsky.social"})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: base64.StdEncoding.EncodeToString(data)})
	if got := userFromSession(req); got != "BlackYellow@bsky.social" {
		t.Errorf("expected handle, got %q", got)
	}
}
