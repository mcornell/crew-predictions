package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

type FirebaseToken struct {
	UID           string
	Email         string
	DisplayName   string
	EmailVerified bool
	Provider      string
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
		log.Printf("VerifyIDToken failed: %v", err)
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	handle := tok.DisplayName
	if handle == "" {
		handle = tok.Email
	}
	writeSessionCookie(w, sessionPayload{
		UserID:        "firebase:" + tok.UID,
		Handle:        handle,
		Provider:      tok.Provider,
		EmailVerified: tok.EmailVerified,
	})
	w.WriteHeader(http.StatusOK)
}
