package handlers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
)

type FirebaseToken struct {
	UID      string
	Email    string
	Provider string
}

type TokenVerifier interface {
	VerifyIDToken(ctx context.Context, idToken string) (*FirebaseToken, error)
}

// NoopTokenVerifier is used in local dev when Firebase credentials are absent.
type NoopTokenVerifier struct{}

func (NoopTokenVerifier) VerifyIDToken(ctx context.Context, idToken string) (*FirebaseToken, error) {
	return nil, fmt.Errorf("Firebase Auth not configured")
}

type SessionHandler struct {
	verifier TokenVerifier
}

func NewSessionHandler(v TokenVerifier) *SessionHandler {
	return &SessionHandler{verifier: v}
}

func (h *SessionHandler) Create(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	idToken := r.FormValue("idToken")
	if idToken == "" {
		http.Error(w, "missing idToken", http.StatusBadRequest)
		return
	}

	tok, err := h.verifier.VerifyIDToken(r.Context(), idToken)
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	sessionData, _ := json.Marshal(map[string]string{
		"userID":   "firebase:" + tok.UID,
		"handle":   tok.Email,
		"provider": tok.Provider,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    base64.StdEncoding.EncodeToString(sessionData),
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, "/matches", http.StatusFound)
}
