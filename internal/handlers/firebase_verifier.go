package handlers

import (
	"context"

	firebaseauth "firebase.google.com/go/v4/auth"
)

type FirebaseTokenVerifier struct {
	client *firebaseauth.Client
}

func NewFirebaseTokenVerifier(client *firebaseauth.Client) *FirebaseTokenVerifier {
	return &FirebaseTokenVerifier{client: client}
}

func (v *FirebaseTokenVerifier) VerifyIDToken(ctx context.Context, idToken string) (*FirebaseToken, error) {
	token, err := v.client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, err
	}
	email, _ := token.Claims["email"].(string)
	displayName, _ := token.Claims["name"].(string)
	emailVerified, _ := token.Claims["email_verified"].(bool)
	return &FirebaseToken{
		UID:           token.UID,
		Email:         email,
		DisplayName:   displayName,
		EmailVerified: emailVerified,
		Provider:      token.Firebase.SignInProvider,
	}, nil
}
