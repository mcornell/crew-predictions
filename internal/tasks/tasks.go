// Package tasks abstracts the Cloud Tasks queue used to drive the
// match-polling chain externally to the Cloud Run container. The interface
// has a real (Cloud Tasks SDK) implementation and a fake (in-memory)
// implementation for tests.
//
// Design rationale lives in docs/match-polling-architecture.md.
package tasks

import (
	"context"
	"sync"
	"time"
)

// Enqueuer schedules a future POST to /admin/poll-scores?matchID=... at the
// given run-at time. The implementation is responsible for naming the task
// deterministically (matchID + run-at unix seconds) so duplicate enqueues
// within the Cloud Tasks 1h dedup window are no-ops.
type Enqueuer interface {
	EnqueuePoll(ctx context.Context, matchID string, runAt time.Time) error
}

// Call is a record of one EnqueuePoll invocation, used by FakeEnqueuer for
// test assertions.
type Call struct {
	MatchID string
	RunAt   time.Time
}

// FakeEnqueuer records EnqueuePoll calls in memory. Tests inspect the
// recorded calls via Calls() to assert chain-seeding behavior without
// hitting real Cloud Tasks.
type FakeEnqueuer struct {
	mu    sync.Mutex
	calls []Call
}

func NewFakeEnqueuer() *FakeEnqueuer {
	return &FakeEnqueuer{}
}

func (f *FakeEnqueuer) EnqueuePoll(_ context.Context, matchID string, runAt time.Time) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls = append(f.calls, Call{MatchID: matchID, RunAt: runAt})
	return nil
}

// Calls returns a copy of the recorded enqueue calls, in order.
func (f *FakeEnqueuer) Calls() []Call {
	f.mu.Lock()
	defer f.mu.Unlock()
	out := make([]Call, len(f.calls))
	copy(out, f.calls)
	return out
}

// Reset clears the recorded calls. Useful between scenarios.
func (f *FakeEnqueuer) Reset() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls = nil
}
