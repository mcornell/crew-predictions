package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/mcornell/crew-predictions/internal/repository"
)

type errUserStore struct{ repository.UserStore }

func (e *errUserStore) Upsert(_ context.Context, _ repository.User) error {
	return fmt.Errorf("store failed")
}
func (e *errUserStore) GetAll(_ context.Context) ([]repository.User, error) {
	return nil, fmt.Errorf("store failed")
}

func TestHandleHandler_Returns500WhenUpsertFails(t *testing.T) {
	h := NewHandleHandler(&errUserStore{UserStore: repository.NewMemoryUserStore()})
	form := url.Values{"handle": {"CrewForever"}}
	req := httptest.NewRequest(http.MethodPost, "/auth/handle", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(sessionCookie("firebase:abc", "oldfan"))
	w := httptest.NewRecorder()
	h.Update(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestBackfillUsersHandler_Returns500WhenGetAllFails(t *testing.T) {
	errPreds := &errGetAllStore{}
	h := NewBackfillUsersHandler(errPreds, repository.NewMemoryUserStore())
	req := httptest.NewRequest(http.MethodPost, "/admin/backfill-users", nil)
	w := httptest.NewRecorder()
	h.Backfill(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestBackfillUsersHandler_Returns500WhenUpsertFails(t *testing.T) {
	predictions := repository.NewMemoryPredictionStore()
	predictions.Save(context.Background(), repository.Prediction{MatchID: "m1", UserID: "firebase:u1", Handle: "fan"})
	h := NewBackfillUsersHandler(predictions, &errUserStore{UserStore: repository.NewMemoryUserStore()})
	req := httptest.NewRequest(http.MethodPost, "/admin/backfill-users", nil)
	w := httptest.NewRecorder()
	h.Backfill(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

type errGetAllStore struct{ repository.PredictionStore }

func (e *errGetAllStore) GetAll(_ context.Context) ([]repository.Prediction, error) {
	return nil, fmt.Errorf("db down")
}
func (e *errGetAllStore) Save(_ context.Context, _ repository.Prediction) error { return nil }
func (e *errGetAllStore) GetByMatchAndUser(_ context.Context, _, _ string) (*repository.Prediction, error) {
	return nil, nil
}
