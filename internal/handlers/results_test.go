package handlers_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mcornell/crew-predictions/internal/handlers"
	"github.com/mcornell/crew-predictions/internal/repository"
)

func TestResultsHandler_SavesResult(t *testing.T) {
	store := repository.NewMemoryResultStore()
	rh := handlers.NewResultsHandler(store)

	req := httptest.NewRequest(http.MethodPost, "/admin/results",
		strings.NewReader("match_id=m1&home_goals=2&away_goals=0"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	rh.Submit(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", w.Code)
	}
	got, err := store.GetResult(context.Background(), "m1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil || got.HomeGoals != 2 || got.AwayGoals != 0 {
		t.Errorf("expected result 2-0, got %+v", got)
	}
}

func TestResultsHandler_RejectsBadHomeGoals(t *testing.T) {
	rh := handlers.NewResultsHandler(repository.NewMemoryResultStore())
	req := httptest.NewRequest(http.MethodPost, "/admin/results",
		strings.NewReader("match_id=m1&home_goals=bad&away_goals=0"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	rh.Submit(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestResultsHandler_RejectsNegativeGoals(t *testing.T) {
	rh := handlers.NewResultsHandler(repository.NewMemoryResultStore())
	req := httptest.NewRequest(http.MethodPost, "/admin/results",
		strings.NewReader("match_id=m1&home_goals=-1&away_goals=0"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	rh.Submit(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for negative goals, got %d", w.Code)
	}
}

func TestResultsHandler_RejectsTooLargeGoals(t *testing.T) {
	rh := handlers.NewResultsHandler(repository.NewMemoryResultStore())
	req := httptest.NewRequest(http.MethodPost, "/admin/results",
		strings.NewReader("match_id=m1&home_goals=0&away_goals=100"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	rh.Submit(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for goals > 99, got %d", w.Code)
	}
}

func TestResultsHandler_RejectsEmptyMatchID(t *testing.T) {
	rh := handlers.NewResultsHandler(repository.NewMemoryResultStore())
	req := httptest.NewRequest(http.MethodPost, "/admin/results",
		strings.NewReader("match_id=&home_goals=2&away_goals=0"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	rh.Submit(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for empty match_id, got %d", w.Code)
	}
}

func TestResultsHandler_RejectsBadAwayGoals(t *testing.T) {
	rh := handlers.NewResultsHandler(repository.NewMemoryResultStore())
	req := httptest.NewRequest(http.MethodPost, "/admin/results",
		strings.NewReader("match_id=m1&home_goals=2&away_goals=bad"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	rh.Submit(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}
