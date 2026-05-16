package main

import "os"

// Config groups all environment-derived settings the server needs at startup.
// Keep it minimal — fields here should map 1-to-1 with env vars or be derived
// trivially from them. Anything more complex (constructed stores, handlers,
// etc.) belongs in Stores or Deps.
type Config struct {
	Port             string // PORT, defaults to 8080
	TestMode         bool   // TEST_MODE=1 → true; switches stores to in-memory + enables seed/reset endpoints
	FirestoreProject string // GOOGLE_CLOUD_PROJECT — empty means use in-memory stores
	FirebaseProject  string // FIREBASE_PROJECT_ID, falls back to FirestoreProject
	AuthEmulatorHost string // FIREBASE_AUTH_EMULATOR_HOST — when set, skip Firebase Auth credentials
	AdminKey         string // ADMIN_KEY — required in production, used for AdminAuth middleware
	SessionSecret    string // SESSION_SECRET — required in production, used to HMAC the session cookie
	TargetTeam       string // hardcoded "Columbus Crew" for now; not env-driven

	// Cloud Tasks routing — when all four are set, the server enqueues
	// chain tasks for live-match polling (see docs/match-polling-architecture.md).
	// Unset (any blank) → no chain enqueueing, falls back to the in-process
	// MatchPoller goroutine only.
	CloudTasksProject  string // CLOUD_TASKS_PROJECT, usually == FirestoreProject
	CloudTasksLocation string // CLOUD_TASKS_LOCATION, e.g. "us-east5"
	CloudTasksQueue    string // CLOUD_TASKS_QUEUE, e.g. "match-polling"
	CloudTasksTarget   string // CLOUD_TASKS_TARGET_URL, public URL of /admin/poll-scores on this Cloud Run service
}

func loadConfig() Config {
	cfg := Config{
		Port:             os.Getenv("PORT"),
		TestMode:         os.Getenv("TEST_MODE") == "1",
		FirestoreProject: os.Getenv("GOOGLE_CLOUD_PROJECT"),
		FirebaseProject:  os.Getenv("FIREBASE_PROJECT_ID"),
		AuthEmulatorHost: os.Getenv("FIREBASE_AUTH_EMULATOR_HOST"),
		AdminKey:           os.Getenv("ADMIN_KEY"),
		SessionSecret:      os.Getenv("SESSION_SECRET"),
		TargetTeam:         "Columbus Crew",
		CloudTasksProject:  os.Getenv("CLOUD_TASKS_PROJECT"),
		CloudTasksLocation: os.Getenv("CLOUD_TASKS_LOCATION"),
		CloudTasksQueue:    os.Getenv("CLOUD_TASKS_QUEUE"),
		CloudTasksTarget:   os.Getenv("CLOUD_TASKS_TARGET_URL"),
	}
	if cfg.Port == "" {
		cfg.Port = "8080"
	}
	if cfg.FirebaseProject == "" {
		cfg.FirebaseProject = cfg.FirestoreProject
	}
	return cfg
}
