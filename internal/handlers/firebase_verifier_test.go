package handlers

import "testing"

func TestTokenToFirebaseToken_MapsAllClaims(t *testing.T) {
	claims := map[string]interface{}{
		"email":          "fan@nordecke.com",
		"name":           "Nordecke Fan",
		"email_verified": true,
	}
	got := tokenToFirebaseToken("uid-1", claims, "google.com")

	if got.UID != "uid-1" {
		t.Errorf("expected UID uid-1, got %q", got.UID)
	}
	if got.Email != "fan@nordecke.com" {
		t.Errorf("expected email fan@nordecke.com, got %q", got.Email)
	}
	if got.DisplayName != "Nordecke Fan" {
		t.Errorf("expected display name Nordecke Fan, got %q", got.DisplayName)
	}
	if !got.EmailVerified {
		t.Error("expected EmailVerified true")
	}
	if got.Provider != "google.com" {
		t.Errorf("expected provider google.com, got %q", got.Provider)
	}
}

func TestTokenToFirebaseToken_MissingClaimsDefaultToZeroValues(t *testing.T) {
	got := tokenToFirebaseToken("uid-2", map[string]interface{}{}, "password")

	if got.Email != "" {
		t.Errorf("expected empty email, got %q", got.Email)
	}
	if got.DisplayName != "" {
		t.Errorf("expected empty display name, got %q", got.DisplayName)
	}
	if got.EmailVerified {
		t.Error("expected EmailVerified false when claim missing")
	}
}

func TestTokenToFirebaseToken_WrongClaimTypesDefaultToZeroValues(t *testing.T) {
	claims := map[string]interface{}{
		"email":          42,
		"name":           false,
		"email_verified": "yes",
	}
	got := tokenToFirebaseToken("uid-3", claims, "password")

	if got.Email != "" {
		t.Errorf("expected empty email on wrong type, got %q", got.Email)
	}
	if got.DisplayName != "" {
		t.Errorf("expected empty display name on wrong type, got %q", got.DisplayName)
	}
	if got.EmailVerified {
		t.Error("expected EmailVerified false on wrong type")
	}
}
