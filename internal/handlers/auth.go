package handlers

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/mcornell/crew-predictions/internal/repository"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, "/matches", http.StatusFound)
}

type sessionPayload struct {
	UserID        string `json:"userID"`
	Handle        string `json:"handle"`
	Provider      string `json:"provider"`
	EmailVerified bool   `json:"emailVerified"`
}

func UserFromSession(r *http.Request) *repository.User {
	cookie, err := r.Cookie("session")
	if err != nil {
		return nil
	}
	data, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		return nil
	}
	var session sessionPayload
	if err := json.Unmarshal(data, &session); err != nil {
		return nil
	}
	if session.UserID == "" || session.Handle == "" {
		return nil
	}
	return &repository.User{
		UserID:        session.UserID,
		Handle:        session.Handle,
		Provider:      session.Provider,
		EmailVerified: session.EmailVerified,
	}
}
