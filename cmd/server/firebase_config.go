package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// serveFirebaseConfig writes a tiny JS snippet that publishes the Firebase
// web SDK config on window.__firebaseConfig. Values are pulled from env at
// request time (rather than baked in at build time) so the same image can
// run against staging, prod, and the auth emulator.
func serveFirebaseConfig(w http.ResponseWriter, _ *http.Request) {
	cfg, _ := json.Marshal(map[string]string{
		"apiKey":           os.Getenv("FIREBASE_API_KEY"),
		"authDomain":       os.Getenv("FIREBASE_AUTH_DOMAIN"),
		"projectId":        os.Getenv("FIREBASE_PROJECT_ID"),
		"appId":            os.Getenv("FIREBASE_APP_ID"),
		"measurementId":    os.Getenv("FIREBASE_MEASUREMENT_ID"),
		"authEmulatorHost": os.Getenv("FIREBASE_AUTH_EMULATOR_HOST"),
	})
	w.Header().Set("Content-Type", "application/javascript")
	fmt.Fprintf(w, "window.__firebaseConfig = %s;", cfg)
}
