package handlers_test

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mcornell/crew-predictions/internal/handlers"
	"github.com/mcornell/crew-predictions/internal/repository"
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

func TestSubmitPrediction_RejectsNonIntegerAwayGoals(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/predictions",
		strings.NewReader("match_id=abc&home_goals=2&away_goals=one"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(sessionCookie("BlackAndGold@bsky.mock"))
	w := httptest.NewRecorder()

	handlers.SubmitPrediction(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestSubmitPrediction_SavesPredictionToStore(t *testing.T) {
	store := repository.NewMemoryPredictionStore()
	req := httptest.NewRequest(http.MethodPost, "/predictions",
		strings.NewReader("match_id=match1&home_goals=3&away_goals=1"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("HX-Request", "true")
	req.AddCookie(sessionCookie("BlackAndGold@bsky.mock"))
	w := httptest.NewRecorder()

	handlers.NewPredictionsHandler(store).Submit(w, req)

	got, err := store.GetByMatchAndHandle(req.Context(), "match1", "BlackAndGold@bsky.mock")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil {
		t.Fatal("expected prediction to be saved, got nil")
	}
	if got.HomeGoals != 3 || got.AwayGoals != 1 {
		t.Errorf("expected 3-1, got %d-%d", got.HomeGoals, got.AwayGoals)
	}
}

func TestSubmitPrediction_ReturnsMatchCardWithSavedScore(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/predictions",
		strings.NewReader("match_id=match1&home_goals=3&away_goals=1"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("HX-Request", "true")
	req.AddCookie(sessionCookie("BlackAndGold@bsky.mock"))
	w := httptest.NewRecorder()

	handlers.SubmitPrediction(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "3") || !strings.Contains(body, "1") {
		t.Errorf("expected saved score in response, got: %s", body)
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
