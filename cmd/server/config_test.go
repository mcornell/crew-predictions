package main

import (
	"testing"
)

func TestLoadConfig_DefaultPort(t *testing.T) {
	t.Setenv("PORT", "")
	cfg := loadConfig()
	if cfg.Port != "8080" {
		t.Errorf("expected default port 8080, got %q", cfg.Port)
	}
}

func TestLoadConfig_ExplicitPort(t *testing.T) {
	t.Setenv("PORT", "9090")
	cfg := loadConfig()
	if cfg.Port != "9090" {
		t.Errorf("expected port 9090, got %q", cfg.Port)
	}
}

func TestLoadConfig_TestModeDetected(t *testing.T) {
	t.Setenv("TEST_MODE", "1")
	cfg := loadConfig()
	if !cfg.TestMode {
		t.Error("expected TestMode true when TEST_MODE=1")
	}
}

func TestLoadConfig_TestModeOff(t *testing.T) {
	t.Setenv("TEST_MODE", "")
	cfg := loadConfig()
	if cfg.TestMode {
		t.Error("expected TestMode false when TEST_MODE empty")
	}
}

func TestLoadConfig_FirebaseProjectFallsBackToFirestore(t *testing.T) {
	t.Setenv("FIREBASE_PROJECT_ID", "")
	t.Setenv("GOOGLE_CLOUD_PROJECT", "my-project")
	cfg := loadConfig()
	if cfg.FirebaseProject != "my-project" {
		t.Errorf("expected FirebaseProject to fall back to GOOGLE_CLOUD_PROJECT, got %q", cfg.FirebaseProject)
	}
	if cfg.FirestoreProject != "my-project" {
		t.Errorf("expected FirestoreProject 'my-project', got %q", cfg.FirestoreProject)
	}
}

func TestLoadConfig_FirebaseProjectExplicitWins(t *testing.T) {
	t.Setenv("FIREBASE_PROJECT_ID", "fb-project")
	t.Setenv("GOOGLE_CLOUD_PROJECT", "gc-project")
	cfg := loadConfig()
	if cfg.FirebaseProject != "fb-project" {
		t.Errorf("expected FirebaseProject 'fb-project', got %q", cfg.FirebaseProject)
	}
	if cfg.FirestoreProject != "gc-project" {
		t.Errorf("expected FirestoreProject 'gc-project', got %q", cfg.FirestoreProject)
	}
}

func TestLoadConfig_TargetTeamHardcoded(t *testing.T) {
	cfg := loadConfig()
	if cfg.TargetTeam != "Columbus Crew" {
		t.Errorf("expected TargetTeam 'Columbus Crew' (hardcoded), got %q", cfg.TargetTeam)
	}
}

func TestLoadConfig_PassesThroughSecrets(t *testing.T) {
	t.Setenv("ADMIN_KEY", "admin-secret")
	t.Setenv("SESSION_SECRET", "session-secret")
	t.Setenv("FIREBASE_AUTH_EMULATOR_HOST", "localhost:9099")
	cfg := loadConfig()
	if cfg.AdminKey != "admin-secret" {
		t.Errorf("AdminKey: got %q", cfg.AdminKey)
	}
	if cfg.SessionSecret != "session-secret" {
		t.Errorf("SessionSecret: got %q", cfg.SessionSecret)
	}
	if cfg.AuthEmulatorHost != "localhost:9099" {
		t.Errorf("AuthEmulatorHost: got %q", cfg.AuthEmulatorHost)
	}
}
