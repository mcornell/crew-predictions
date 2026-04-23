package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mcornell/crew-predictions/internal/handlers"
	"github.com/mcornell/crew-predictions/internal/repository"
)

func TestSeedMatchHandler_SavesMatch(t *testing.T) {
	store := repository.NewMemoryMatchStore()
	h := handlers.NewSeedMatchHandler(store)

	req := httptest.NewRequest(http.MethodPost, "/admin/seed-match",
		strings.NewReader("id=m1&home_team=Columbus+Crew&away_team=LA+Galaxy&kickoff=2026-05-01T19:30:00Z&status=STATUS_SCHEDULED"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	h.Submit(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", w.Code)
	}
	matches, err := store.GetAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(matches) != 1 || matches[0].ID != "m1" {
		t.Errorf("expected match m1, got %+v", matches)
	}
}

func TestSeedMatchHandler_SavesMatchWithState(t *testing.T) {
	store := repository.NewMemoryMatchStore()
	h := handlers.NewSeedMatchHandler(store)

	req := httptest.NewRequest(http.MethodPost, "/admin/seed-match",
		strings.NewReader("id=m2&home_team=Columbus+Crew&away_team=LA+Galaxy&kickoff=2026-05-01T19:30:00Z&status=STATUS_SCHEDULED&state=in"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	h.Submit(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", w.Code)
	}
	matches, _ := store.GetAll()
	if len(matches) != 1 || matches[0].State != "in" {
		t.Errorf("expected state=in, got %+v", matches)
	}
}

func TestSeedMatchHandler_SavesMatchWithScores(t *testing.T) {
	store := repository.NewMemoryMatchStore()
	h := handlers.NewSeedMatchHandler(store)

	req := httptest.NewRequest(http.MethodPost, "/admin/seed-match",
		strings.NewReader("id=m3&home_team=Columbus+Crew&away_team=FC+Dallas&kickoff=2026-05-01T19:30:00Z&status=STATUS_FULL_TIME&state=post&home_score=3&away_score=1"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	h.Submit(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", w.Code)
	}
	matches, _ := store.GetAll()
	if len(matches) != 1 || matches[0].HomeScore != "3" || matches[0].AwayScore != "1" {
		t.Errorf("expected HomeScore=3 AwayScore=1, got %+v", matches)
	}
}

func TestSeedMatchHandler_RejectsEmptyID(t *testing.T) {
	h := handlers.NewSeedMatchHandler(repository.NewMemoryMatchStore())
	req := httptest.NewRequest(http.MethodPost, "/admin/seed-match",
		strings.NewReader("id=&home_team=Columbus+Crew&away_team=LA+Galaxy&kickoff=2026-05-01T19:30:00Z&status=STATUS_SCHEDULED"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	h.Submit(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for empty id, got %d", w.Code)
	}
}

func TestSeedMatchHandler_RejectsBadKickoff(t *testing.T) {
	h := handlers.NewSeedMatchHandler(repository.NewMemoryMatchStore())
	req := httptest.NewRequest(http.MethodPost, "/admin/seed-match",
		strings.NewReader("id=m1&home_team=Columbus+Crew&away_team=LA+Galaxy&kickoff=not-a-date&status=STATUS_SCHEDULED"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	h.Submit(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}
