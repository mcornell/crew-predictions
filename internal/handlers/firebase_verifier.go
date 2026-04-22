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

func tokenToFirebaseToken(uid string, claims map[string]interface{}, provider string) *FirebaseToken {
	email, _ := claims["email"].(string)
	displayName, _ := claims["name"].(string)
	emailVerified, _ := claims["email_verified"].(bool)
	return &FirebaseToken{
		UID:           uid,
		Email:         email,
		DisplayName:   displayName,
		EmailVerified: emailVerified,
		Provider:      provider,
	}
}

func (v *FirebaseTokenVerifier) VerifyIDToken(ctx context.Context, idToken string) (*FirebaseToken, error) {
	token, err := v.client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, err
	}
	return tokenToFirebaseToken(token.UID, token.Claims, token.Firebase.SignInProvider), nil
}
