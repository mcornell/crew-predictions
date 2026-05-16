package tasks

import (
	"context"
	"testing"
)

func TestNewCloudTasksEnqueuer_RejectsIncompleteConfig(t *testing.T) {
	cases := []struct {
		name string
		cfg  CloudTasksConfig
	}{
		{"missing ProjectID", CloudTasksConfig{Location: "us-east5", QueueID: "match-polling", TargetURL: "https://x/p"}},
		{"missing Location", CloudTasksConfig{ProjectID: "p", QueueID: "match-polling", TargetURL: "https://x/p"}},
		{"missing QueueID", CloudTasksConfig{ProjectID: "p", Location: "us-east5", TargetURL: "https://x/p"}},
		{"missing TargetURL", CloudTasksConfig{ProjectID: "p", Location: "us-east5", QueueID: "match-polling"}},
		{"all empty", CloudTasksConfig{}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewCloudTasksEnqueuer(context.Background(), tc.cfg)
			if err == nil {
				t.Errorf("expected error for incomplete config, got nil")
			}
		})
	}
}

func TestSanitizeMatchID(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"761621", "761621"},
		{"m-live-123", "m-live-123"},
		{"with/slash", "with_slash"},
		{"weird:char", "weird_char"},
		{"unicode-café", "unicode-caf_"}, // é is one rune → one underscore
	}
	for _, tc := range cases {
		got := sanitizeMatchID(tc.in)
		if got != tc.want {
			t.Errorf("sanitizeMatchID(%q): got %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestIsAlreadyExists(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want bool
	}{
		{"nil", nil, false},
		{"unrelated error", &stringErr{"something went wrong"}, false},
		{"AlreadyExists in message", &stringErr{"rpc error: code = AlreadyExists desc = task already exists"}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := isAlreadyExists(tc.err); got != tc.want {
				t.Errorf("isAlreadyExists(%v): got %v, want %v", tc.err, got, tc.want)
			}
		})
	}
}

type stringErr struct{ s string }

func (e *stringErr) Error() string { return e.s }
