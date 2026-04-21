package main

import (
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

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
