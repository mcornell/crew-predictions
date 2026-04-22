package handlers

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/mcornell/crew-predictions/internal/repository"
)

type HandleHandler struct {
	users repository.UserStore
}

func NewHandleHandler(users repository.UserStore) *HandleHandler {
	return &HandleHandler{users: users}
}

func (h *HandleHandler) Update(w http.ResponseWriter, r *http.Request) {
	user := UserFromSession(r)
	if user == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	r.ParseForm()
	handle := r.FormValue("handle")
	if handle == "" {
		http.Error(w, "handle is required", http.StatusBadRequest)
		return
	}

	if err := h.users.Upsert(r.Context(), repository.User{
		UserID:   user.UserID,
		Handle:   handle,
		Provider: user.Provider,
	}); err != nil {
		http.Error(w, "could not save handle", http.StatusInternalServerError)
		return
	}

	writeSessionCookie(w, sessionPayload{
		UserID:        user.UserID,
		Handle:        handle,
		Provider:      user.Provider,
		EmailVerified: user.EmailVerified,
	})
	w.WriteHeader(http.StatusOK)
}

func writeSessionCookie(w http.ResponseWriter, payload sessionPayload) {
	data, _ := json.Marshal(payload)
	http.SetCookie(w, &http.Cookie{
		Name:     "__session",
		Value:    base64.StdEncoding.EncodeToString(data),
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}
