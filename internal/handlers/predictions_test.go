package handlers_test

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mcornell/crew-predictions/internal/handlers"
)

func sessionCookie(handle string) *http.Cookie {
	data, _ := json.Marshal(map[string]string{"handle": handle})
	return &http.Cookie{Name: "session", Value: base64.StdEncoding.EncodeToString(data)}
}

func TestSubmitPrediction_RejectsMissingFields(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/predictions", strings.NewReader(""))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(sessionCookie("BlackAndGold@bsky.mock"))
	w := httptest.NewRecorder()

	handlers.SubmitPrediction(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestSubmitPrediction_RejectsNonIntegerGoals(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/predictions",
		strings.NewReader("match_id=abc&home_goals=two&away_goals=one"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(sessionCookie("BlackAndGold@bsky.mock"))
	w := httptest.NewRecorder()

	handlers.SubmitPrediction(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestSubmitPrediction_AcceptsValidRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/predictions",
		strings.NewReader("match_id=match1&home_goals=3&away_goals=1"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(sessionCookie("BlackAndGold@bsky.mock"))
	w := httptest.NewRecorder()

	handlers.SubmitPrediction(w, req)

	if w.Code != http.StatusFound {
		t.Errorf("expected 302 redirect, got %d", w.Code)
	}
	if w.Header().Get("Location") != "/matches" {
		t.Errorf("expected redirect to /matches, got %s", w.Header().Get("Location"))
	}
}

func TestSubmitPrediction_RejectsUnauthenticated(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/predictions", nil)
	w := httptest.NewRecorder()

	handlers.SubmitPrediction(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}
