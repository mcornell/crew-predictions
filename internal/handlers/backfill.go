package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mcornell/crew-predictions/internal/repository"
)

type BackfillUsersHandler struct {
	predictions repository.PredictionStore
	users       repository.UserStore
}

func NewBackfillUsersHandler(predictions repository.PredictionStore, users repository.UserStore) *BackfillUsersHandler {
	return &BackfillUsersHandler{predictions: predictions, users: users}
}

func (h *BackfillUsersHandler) Backfill(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	preds, err := h.predictions.GetAll(ctx)
	if err != nil {
		http.Error(w, "could not load predictions", http.StatusInternalServerError)
		return
	}

	seen := map[string]struct{}{}
	for _, p := range preds {
		if p.UserID == "" {
			continue
		}
		seen[p.UserID] = struct{}{}
	}

	for userID := range seen {
		if err := h.users.Upsert(ctx, repository.User{UserID: userID}); err != nil {
			http.Error(w, "upsert failed", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]int{"backfilled": len(seen)}); err != nil {
		log.Printf("backfill: encode response: %v", err)
	}
}
