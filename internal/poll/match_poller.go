package poll

import (
	"context"
	"log"
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
		p.timerFunc(s.delay, func() {
			p.mu.Lock()
			p.active[matchID] = true
			p.mu.Unlock()
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

	matches, err := p.fetcher()
	if err != nil {
		log.Printf("score poll fetch failed: %v", err)
		return
	}
	if err := p.matchStore.SaveAll(matches); err != nil {
		log.Printf("score poll store update failed: %v", err)
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
		if ok && terminalStatuses[m.Status] {
			toSave = append(toSave, m)
			delete(p.active, matchID)
			delete(p.scheduled, matchID)
		}
	}
	p.mu.Unlock()

	for _, m := range toSave {
		if err := saveResult(ctx, p.resultStore, m); err != nil {
			log.Printf("score poll save result failed for %s: %v", m.ID, err)
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

// ActiveCount returns the number of matches currently being polled (for testing).
func (p *MatchPoller) ActiveCount() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.active)
}
