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

func UserFromSession(r *http.Request) *repository.User {
	cookie, err := r.Cookie("session")
	if err != nil {
		return nil
	}
	data, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		return nil
	}
	var session map[string]string
	if err := json.Unmarshal(data, &session); err != nil {
		return nil
	}
	userID := session["userID"]
	handle := session["handle"]
	if userID == "" || handle == "" {
		return nil
	}
	return &repository.User{
		UserID:   userID,
		Handle:   handle,
		Provider: session["provider"],
	}
}
