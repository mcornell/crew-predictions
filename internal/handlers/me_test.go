package handlers_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mcornell/crew-predictions/internal/handlers"
	"github.com/mcornell/crew-predictions/internal/repository"
)

func TestMeHandler_UpsertsProviderFromSession(t *testing.T) {
	users := repository.NewMemoryUserStore()
	h := handlers.NewMeHandler(users)
	data, _ := json.Marshal(map[string]interface{}{"userID": "google:abc", "handle": "Fan", "provider": "google.com"})
	req := httptest.NewRequest(http.MethodGet, "/api/me", nil)
	req.AddCookie(&http.Cookie{Name: "__session", Value: base64.StdEncoding.EncodeToString(data)})
	h.Get(httptest.NewRecorder(), req)

	u, _ := users.GetByID(context.Background(), "google:abc")
	if u == nil || u.Provider != "google.com" {
		t.Errorf("expected provider google.com stored, got %+v", u)
	}
}

func TestMeHandler_UpsertsUserToStoreOnValidSession(t *testing.T) {
	users := repository.NewMemoryUserStore()
	h := handlers.NewMeHandler(users)
	req := httptest.NewRequest(http.MethodGet, "/api/me", nil)
	req.AddCookie(sessionCookie("google:abc123", "BlackAndGold@bsky.mock"))
	w := httptest.NewRecorder()
	h.Get(w, req)

	u, _ := users.GetByID(context.Background(), "google:abc123")
	if u == nil {
		t.Fatal("expected user upserted into store, got nil")
	}
	if u.Handle != "BlackAndGold@bsky.mock" {
		t.Errorf("expected handle BlackAndGold@bsky.mock, got %s", u.Handle)
	}
}

func TestMeHandler_ReturnsUserWhenLoggedIn(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/me", nil)
	req.AddCookie(sessionCookie("google:abc123", "BlackAndGold@bsky.mock"))
	w := httptest.NewRecorder()

	handlers.NewMeHandler(repository.NewMemoryUserStore()).Get(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var body struct {
		Handle string `json:"handle"`
	}
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if body.Handle != "BlackAndGold@bsky.mock" {
		t.Errorf("expected handle BlackAndGold@bsky.mock, got %s", body.Handle)
	}
}

func TestMeHandler_Returns401WhenNotLoggedIn(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/me", nil)
	w := httptest.NewRecorder()

	handlers.NewMeHandler(repository.NewMemoryUserStore()).Get(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestMeHandler_ReturnsEmailVerifiedInResponse(t *testing.T) {
	data, _ := json.Marshal(map[string]interface{}{"userID": "google:abc", "handle": "Fan", "emailVerified": true})
	cookie := &http.Cookie{Name: "__session", Value: base64.StdEncoding.EncodeToString(data)}
	req := httptest.NewRequest(http.MethodGet, "/api/me", nil)
	req.AddCookie(cookie)
	w := httptest.NewRecorder()

	handlers.NewMeHandler(repository.NewMemoryUserStore()).Get(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var body struct {
		Handle        string `json:"handle"`
		EmailVerified bool   `json:"emailVerified"`
	}
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if !body.EmailVerified {
		t.Error("expected emailVerified to be true in response")
	}
}
