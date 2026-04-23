package handlers

import (
	"crypto/subtle"
	"net/http"
	"os"
)

// AdminAuth guards admin endpoints with an X-Admin-Key header checked against
// the ADMIN_KEY env var. TEST_MODE bypasses the check so e2e tests work without
// threading the key through Playwright config.
func AdminAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if os.Getenv("TEST_MODE") == "1" {
			next(w, r)
			return
		}
		expected := os.Getenv("ADMIN_KEY")
		provided := r.Header.Get("X-Admin-Key")
		if expected == "" || subtle.ConstantTimeCompare([]byte(provided), []byte(expected)) != 1 {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		next(w, r)
	}
}
