package espn

import (
	"testing"
	"time"
)

func TestParseKickoff_RFC3339WithSeconds(t *testing.T) {
	k, err := parseKickoff("2026-05-01T20:00:00Z")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if k.Year() != 2026 || k.Month() != time.May || k.Day() != 1 {
		t.Errorf("unexpected date: %v", k)
	}
}

func TestParseKickoff_ESPNFormatNoSeconds(t *testing.T) {
	k, err := parseKickoff("2026-04-19T23:00Z")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if k.Year() != 2026 || k.Month() != time.April || k.Day() != 19 {
		t.Errorf("unexpected date: %v", k)
	}
}

func TestParseKickoff_InvalidReturnsError(t *testing.T) {
	_, err := parseKickoff("not-a-date")
	if err == nil {
		t.Error("expected error for invalid date")
	}
}
