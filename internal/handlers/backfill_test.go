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

	predictions.Save(ctx, repository.Prediction{MatchID: "m1", UserID: "firebase:u1", Handle: "crewfan"})
	predictions.Save(ctx, repository.Prediction{MatchID: "m2", UserID: "firebase:u1", Handle: "crewfan"})
	predictions.Save(ctx, repository.Prediction{MatchID: "m1", UserID: "firebase:u2", Handle: "northend96"})

	h := NewBackfillUsersHandler(predictions, users)
	req := httptest.NewRequest(http.MethodPost, "/admin/backfill-users", nil)
	w := httptest.NewRecorder()
	h.Backfill(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	u1, _ := users.GetByID(ctx, "firebase:u1")
	if u1 == nil || u1.Handle != "crewfan" {
		t.Errorf("expected u1 crewfan, got %+v", u1)
	}
	u2, _ := users.GetByID(ctx, "firebase:u2")
	if u2 == nil || u2.Handle != "northend96" {
		t.Errorf("expected u2 northend96, got %+v", u2)
	}
}

func TestBackfillUsersHandler_SkipsPredictionsWithNoUserID(t *testing.T) {
	predictions := repository.NewMemoryPredictionStore()
	users := repository.NewMemoryUserStore()
	ctx := context.Background()

	predictions.Save(ctx, repository.Prediction{MatchID: "m1", UserID: "", Handle: "legacyfan"})

	h := NewBackfillUsersHandler(predictions, users)
	req := httptest.NewRequest(http.MethodPost, "/admin/backfill-users", nil)
	httptest.NewRecorder()
	h.Backfill(httptest.NewRecorder(), req)

	all, _ := users.GetAll(ctx)
	if len(all) != 0 {
		t.Errorf("expected no users for legacy predictions, got %d", len(all))
	}
}
