package poll

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/mcornell/crew-predictions/internal/models"
	"github.com/mcornell/crew-predictions/internal/repository"
)

type MatchPoller struct {
	matchStore  repository.MatchStore
	resultStore repository.ResultStore
	fetcher     func() ([]models.Match, error)

	mu        sync.Mutex
	scheduled map[string]bool
	active    map[string]bool
	timerFunc func(time.Duration, func())
}

func NewMatchPoller(
	matchStore repository.MatchStore,
	resultStore repository.ResultStore,
	fetcher func() ([]models.Match, error),
	timerFunc func(time.Duration, func()),
) *MatchPoller {
	return &MatchPoller{
		matchStore:  matchStore,
		resultStore: resultStore,
		fetcher:     fetcher,
		timerFunc:   timerFunc,
		scheduled:   make(map[string]bool),
		active:      make(map[string]bool),
	}
}

// SetTimerFunc replaces the timer function; used in tests to switch from
// immediate firing to delay-capturing between Schedule calls.
func (p *MatchPoller) SetTimerFunc(f func(time.Duration, func())) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.timerFunc = f
}

type scheduledMatch struct {
	matchID string
	delay   time.Duration
}

func (p *MatchPoller) Schedule(matches []models.Match) {
	p.mu.Lock()
	var toSchedule []scheduledMatch
	for _, m := range matches {
		if p.scheduled[m.ID] || terminalStatuses[m.Status] {
			continue
		}
		p.scheduled[m.ID] = true
		delay := time.Until(m.Kickoff)
		if delay < 0 {
			delay = 0
		}
		toSchedule = append(toSchedule, scheduledMatch{m.ID, delay})
	}
	p.mu.Unlock()

	for _, s := range toSchedule {
		matchID := s.matchID
		delay := s.delay
		slog.Info("poller: match scheduled", "matchID", matchID, "delaySeconds", int(delay.Seconds()))
		p.timerFunc(delay, func() {
			p.mu.Lock()
			p.active[matchID] = true
			p.mu.Unlock()
			slog.Info("poller: kickoff reached, entering active polling", "matchID", matchID)
		})
	}
}

// Reset cancels all active pollers, clears state, and reschedules from matches.
func (p *MatchPoller) Reset(matches []models.Match) {
	p.mu.Lock()
	p.active = make(map[string]bool)
	p.scheduled = make(map[string]bool)
	p.mu.Unlock()
	p.Schedule(matches)
}

// Tick fetches current match data, updates the store, and records results for
// any active matches that have reached a terminal status. It is a no-op when
// no matches are active.
func (p *MatchPoller) Tick(ctx context.Context) {
	p.mu.Lock()
	if len(p.active) == 0 {
		p.mu.Unlock()
		return
	}
	p.mu.Unlock()

	p.mu.Lock()
	activeCount := len(p.active)
	p.mu.Unlock()
	slog.Info("poller: tick", "activeMatches", activeCount)

	matches, err := p.fetcher()
	if err != nil {
		slog.Error("poller: ESPN fetch failed", "error", err)
		return
	}
	if err := p.matchStore.SaveAll(matches); err != nil {
		slog.Error("poller: store update failed", "error", err)
		return
	}

	byID := make(map[string]models.Match, len(matches))
	for _, m := range matches {
		byID[m.ID] = m
	}

	p.mu.Lock()
	var toSave []models.Match
	for matchID := range p.active {
		m, ok := byID[matchID]
		if !ok {
			continue
		}
		if terminalStatuses[m.Status] {
			toSave = append(toSave, m)
			delete(p.active, matchID)
			delete(p.scheduled, matchID)
		} else {
			slog.Info("poller: match still in progress", "matchID", matchID, "status", m.Status, "homeScore", m.HomeScore, "awayScore", m.AwayScore)
		}
	}
	p.mu.Unlock()

	for _, m := range toSave {
		slog.Info("poller: match finished, saving result", "matchID", m.ID, "status", m.Status, "homeScore", m.HomeScore, "awayScore", m.AwayScore)
		if err := saveResult(ctx, p.resultStore, m); err != nil {
			slog.Error("poller: save result failed", "matchID", m.ID, "error", err)
		}
	}
}

// Run starts the polling loop. Call in a goroutine; cancel ctx to stop.
func (p *MatchPoller) Run(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			p.Tick(ctx)
		case <-ctx.Done():
			return
		}
	}
}

// Backfill saves results for any terminal-status matches in the provided list
// without re-fetching from ESPN. Called after a daily refresh to catch matches
// that finished while no active poller was running for them.
func (p *MatchPoller) Backfill(ctx context.Context, matches []models.Match) {
	for _, m := range matches {
		if !terminalStatuses[m.Status] {
			continue
		}
		if err := saveResult(ctx, p.resultStore, m); err != nil {
			slog.Error("poller: backfill result failed", "matchID", m.ID, "error", err)
		}
	}
}

// ActiveCount returns the number of matches currently being polled (for testing).
func (p *MatchPoller) ActiveCount() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.active)
}
