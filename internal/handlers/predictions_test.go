package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mcornell/crew-predictions/internal/handlers"
)

func TestSubmitPrediction_RejectsUnauthenticated(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/predictions", nil)
	w := httptest.NewRecorder()

	handlers.SubmitPrediction(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}
