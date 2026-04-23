package handlers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestAdminAuth_NoHeader_Returns403(t *testing.T) {
	t.Setenv("ADMIN_KEY", "secret")
	t.Setenv("TEST_MODE", "")

	handler := AdminAuth(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	req := httptest.NewRequest(http.MethodPost, "/admin/results", nil)
	rr := httptest.NewRecorder()
	handler(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", rr.Code)
	}
}

func TestAdminAuth_WrongKey_Returns403(t *testing.T) {
	t.Setenv("ADMIN_KEY", "secret")
	t.Setenv("TEST_MODE", "")

	handler := AdminAuth(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	req := httptest.NewRequest(http.MethodPost, "/admin/results", nil)
	req.Header.Set("X-Admin-Key", "wrong")
	rr := httptest.NewRecorder()
	handler(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", rr.Code)
	}
}

func TestAdminAuth_CorrectKey_Passes(t *testing.T) {
	t.Setenv("ADMIN_KEY", "secret")
	t.Setenv("TEST_MODE", "")

	handler := AdminAuth(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	req := httptest.NewRequest(http.MethodPost, "/admin/results", nil)
	req.Header.Set("X-Admin-Key", "secret")
	rr := httptest.NewRecorder()
	handler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestAdminAuth_TestMode_Bypasses(t *testing.T) {
	t.Setenv("ADMIN_KEY", "secret")
	t.Setenv("TEST_MODE", "1")

	handler := AdminAuth(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	req := httptest.NewRequest(http.MethodPost, "/admin/results", nil)
	// No header — should still pass in TEST_MODE
	rr := httptest.NewRecorder()
	handler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200 in TEST_MODE, got %d", rr.Code)
	}
}

func TestAdminAuth_EmptyEnvKey_AlwaysDenies(t *testing.T) {
	orig := os.Getenv("ADMIN_KEY")
	os.Unsetenv("ADMIN_KEY")
	defer os.Setenv("ADMIN_KEY", orig)
	t.Setenv("TEST_MODE", "")

	handler := AdminAuth(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	req := httptest.NewRequest(http.MethodPost, "/admin/results", nil)
	req.Header.Set("X-Admin-Key", "")
	rr := httptest.NewRecorder()
	handler(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Errorf("expected 403 when ADMIN_KEY unset, got %d", rr.Code)
	}
}
