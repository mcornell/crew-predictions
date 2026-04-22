package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestServeFirebaseConfig_ReturnsJavaScriptWithEnvVars(t *testing.T) {
	t.Setenv("FIREBASE_API_KEY", "test-api-key")
	t.Setenv("FIREBASE_PROJECT_ID", "test-project")
	t.Setenv("FIREBASE_AUTH_DOMAIN", "test-project.firebaseapp.com")

	req := httptest.NewRequest(http.MethodGet, "/auth/config.js", nil)
	rr := httptest.NewRecorder()
	serveFirebaseConfig(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "application/javascript" {
		t.Errorf("expected Content-Type application/javascript, got %q", ct)
	}
	body := rr.Body.String()
	if !strings.Contains(body, "window.__firebaseConfig") {
		t.Errorf("expected window.__firebaseConfig assignment, got %q", body)
	}
	if !strings.Contains(body, "test-api-key") {
		t.Errorf("expected API key in response, got %q", body)
	}
	if !strings.Contains(body, "test-project") {
		t.Errorf("expected project ID in response, got %q", body)
	}
}

func TestSPAHandlerNoCacheHeader(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "index.html"), []byte("<html></html>"), 0644); err != nil {
		t.Fatal(err)
	}

	h := spaHandler(dir)
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	cc := rr.Header().Get("Cache-Control")
	if cc != "no-cache" {
		t.Errorf("expected Cache-Control: no-cache, got %q", cc)
	}
}

func TestAssetsImmutableCacheHeader(t *testing.T) {
	dir := t.TempDir()
	assetsDir := filepath.Join(dir, "assets")
	if err := os.MkdirAll(assetsDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(assetsDir, "index-abc123.css"), []byte("body{}"), 0644); err != nil {
		t.Fatal(err)
	}

	h := assetsHandler(dir)
	req := httptest.NewRequest("GET", "/assets/index-abc123.css", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	cc := rr.Header().Get("Cache-Control")
	if cc != "public, max-age=31536000, immutable" {
		t.Errorf("expected immutable cache header, got %q", cc)
	}
}
