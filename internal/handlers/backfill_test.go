package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mcornell/crew-predictions/internal/repository"
)

func TestBackfillUsersHandler_CreatesUserRecordsFromPredictions(t *testing.T) {
	predictions := repository.NewMemoryPredictionStore()
	users := repository.NewMemoryUserStore()
	ctx := context.Background()

	predictions.Save(ctx, repository.Prediction{MatchID: "m1", UserID: "firebase:u1", })
	predictions.Save(ctx, repository.Prediction{MatchID: "m2", UserID: "firebase:u1", })
	predictions.Save(ctx, repository.Prediction{MatchID: "m1", UserID: "firebase:u2", })

	h := NewBackfillUsersHandler(predictions, users)
	req := httptest.NewRequest(http.MethodPost, "/admin/backfill-users", nil)
	w := httptest.NewRecorder()
	h.Backfill(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	u1, _ := users.GetByID(ctx, "firebase:u1")
	if u1 == nil {
		t.Error("expected UserStore entry for firebase:u1")
	}
	u2, _ := users.GetByID(ctx, "firebase:u2")
	if u2 == nil {
		t.Error("expected UserStore entry for firebase:u2")
	}
}

func TestBackfillUsersHandler_SkipsPredictionsWithNoUserID(t *testing.T) {
	predictions := repository.NewMemoryPredictionStore()
	users := repository.NewMemoryUserStore()
	ctx := context.Background()

	predictions.Save(ctx, repository.Prediction{MatchID: "m1", UserID: "", })

	h := NewBackfillUsersHandler(predictions, users)
	req := httptest.NewRequest(http.MethodPost, "/admin/backfill-users", nil)
	httptest.NewRecorder()
	h.Backfill(httptest.NewRecorder(), req)

	all, _ := users.GetAll(ctx)
	if len(all) != 0 {
		t.Errorf("expected no users for legacy predictions, got %d", len(all))
	}
}
