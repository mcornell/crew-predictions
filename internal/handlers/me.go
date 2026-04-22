package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mcornell/crew-predictions/internal/repository"
)

type MeHandler struct {
	users repository.UserStore
}

func NewMeHandler(users repository.UserStore) *MeHandler {
	return &MeHandler{users: users}
}

func (h *MeHandler) Get(w http.ResponseWriter, r *http.Request) {
	user := UserFromSession(r)
	if user == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if user.UserID != "" {
		if err := h.users.Upsert(r.Context(), repository.User{UserID: user.UserID, Handle: user.Handle}); err != nil {
			log.Printf("me: lazy upsert failed for %s: %v", user.UserID, err)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"handle":        user.Handle,
		"emailVerified": user.EmailVerified,
	})
}
