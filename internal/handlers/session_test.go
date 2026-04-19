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
	if got := UserFromSession(req); got != nil {
		t.Errorf("expected nil, got %+v", got)
	}
}

func TestUserFromSession_InvalidBase64(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: "not-valid-base64!!!"})
	if got := UserFromSession(req); got != nil {
		t.Errorf("expected nil for bad base64, got %+v", got)
	}
}

func TestUserFromSession_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: base64.StdEncoding.EncodeToString([]byte("not json"))})
	if got := UserFromSession(req); got != nil {
		t.Errorf("expected nil for bad JSON, got %+v", got)
	}
}

func TestUserFromSession_ReturnsUser(t *testing.T) {
	data, _ := json.Marshal(map[string]string{
		"userID":   "google:110048215615",
		"handle":   "BlackAndGold@bsky.mock",
		"provider": "google",
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: base64.StdEncoding.EncodeToString(data)})
	got := UserFromSession(req)
	if got == nil {
		t.Fatal("expected user, got nil")
	}
	if got.UserID != "google:110048215615" {
		t.Errorf("expected userID google:110048215615, got %q", got.UserID)
	}
	if got.Handle != "BlackAndGold@bsky.mock" {
		t.Errorf("expected handle BlackAndGold@bsky.mock, got %q", got.Handle)
	}
}
