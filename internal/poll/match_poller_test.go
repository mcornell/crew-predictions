package poll_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/mcornell/crew-predictions/internal/models"
	"github.com/mcornell/crew-predictions/internal/poll"
	"github.com/mcornell/crew-predictions/internal/repository"
)

// immediateTimer fires the callback synchronously, ignoring the delay.
func immediateTimer(_ time.Duration, f func()) { f() }

// capturingTimer records delays without firing callbacks.
func capturingTimer(delays *[]time.Duration) func(time.Duration, func()) {
	return func(d time.Duration, _ func()) {
		*delays = append(*delays, d)
	}
}

func newPoller(matchStore repository.MatchStore, resultStore repository.ResultStore, matches []models.Match, timer func(time.Duration, func())) *poll.MatchPoller {
	fetcher := func() ([]models.Match, error) { return matches, nil }
	return poll.NewMatchPoller(matchStore, resultStore, fetcher, timer)
}

func TestMatchPoller_Tick_FetchesSummaryAndPersistsEventsForLiveMatch(t *testing.T) {
	matchStore := repository.NewMemoryMatchStore()
	matchStore.Seed([]models.Match{{
		ID: "m-live", HomeTeam: "Crew", AwayTeam: "FC Dallas",
		Kickoff: time.Now().Add(-30 * time.Minute), State: "in", Status: "STATUS_IN_PROGRESS",
	}})

	scoreboardMatch := models.Match{
		ID: "m-live", HomeTeam: "Crew", AwayTeam: "FC Dallas",
		Kickoff: time.Now().Add(-30 * time.Minute), State: "in", Status: "STATUS_IN_PROGRESS",
		HomeScore: "1", AwayScore: "0",
	}
	scoreboardFetcher := func() ([]models.Match, error) { return []models.Match{scoreboardMatch}, nil }
	summaryCalls := 0
	summaryFetcher := func(id string) (models.MatchSummary, error) {
		summaryCalls++
		return models.MatchSummary{
			Events: []models.MatchEvent{
				{Clock: "23'", TypeID: "goal", Team: "Crew", Players: []string{"Hugo Picard"}},
			},
		}, nil
	}

	p := poll.NewMatchPoller(matchStore, repository.NewMemoryResultStore(), scoreboardFetcher, immediateTimer)
	p.SetSummaryFetcher(summaryFetcher)
	p.Schedule([]models.Match{scoreboardMatch}) // activates the match
	p.Tick(context.Background())

	if summaryCalls != 1 {
		t.Errorf("expected summaryFetcher called exactly once, got %d", summaryCalls)
	}
	stored, _ := matchStore.GetAll()
	if len(stored) == 0 || len(stored[0].Events) != 1 || stored[0].Events[0].Players[0] != "Hugo Picard" {
		t.Errorf("expected events persisted to store; got %+v", stored)
	}
}

func TestMatchPoller_Tick_DoesNotCallSummaryWhenNoneConfigured(t *testing.T) {
	matchStore := repository.NewMemoryMatchStore()
	scoreboardMatch := models.Match{
		ID: "m-no-sum", HomeTeam: "Crew", AwayTeam: "FC Dallas",
		Kickoff: time.Now().Add(-30 * time.Minute), State: "in", Status: "STATUS_IN_PROGRESS",
	}
	matchStore.Seed([]models.Match{scoreboardMatch})
	fetcher := func() ([]models.Match, error) { return []models.Match{scoreboardMatch}, nil }
	p := poll.NewMatchPoller(matchStore, repository.NewMemoryResultStore(), fetcher, immediateTimer)
	p.Schedule([]models.Match{scoreboardMatch})
	// Should not panic with nil summaryFetcher.
	p.Tick(context.Background())
}

func TestMatchPoller_ScheduleSetsPositiveDelayForFutureKickoff(t *testing.T) {
	var delays []time.Duration
	p := newPoller(repository.NewMemoryMatchStore(), repository.NewMemoryResultStore(), nil, capturingTimer(&delays))

	future := time.Now().Add(2 * time.Hour)
	p.Schedule([]models.Match{{ID: "m1", HomeTeam: "Columbus Crew", AwayTeam: "FC Dallas", Kickoff: future, Status: "STATUS_SCHEDULED"}})

	if len(delays) != 1 {
		t.Fatalf("expected 1 timer, got %d", len(delays))
	}
	if delays[0] <= 0 {
		t.Errorf("expected positive delay for future kickoff, got %v", delays[0])
	}
}

func TestMatchPoller_ScheduleUsesZeroDelayForPastKickoff(t *testing.T) {
	var delays []time.Duration
	p := newPoller(repository.NewMemoryMatchStore(), repository.NewMemoryResultStore(), nil, capturingTimer(&delays))

	past := time.Now().Add(-1 * time.Hour)
	p.Schedule([]models.Match{{ID: "m1", HomeTeam: "Columbus Crew", AwayTeam: "FC Dallas", Kickoff: past, Status: "STATUS_SCHEDULED", State: "in"}})

	if len(delays) != 1 {
		t.Fatalf("expected 1 timer, got %d", len(delays))
	}
	if delays[0] != 0 {
		t.Errorf("expected zero delay for past kickoff, got %v", delays[0])
	}
}

func TestMatchPoller_ScheduleSkipsTerminalMatch(t *testing.T) {
	var delays []time.Duration
	p := newPoller(repository.NewMemoryMatchStore(), repository.NewMemoryResultStore(), nil, capturingTimer(&delays))

	p.Schedule([]models.Match{{ID: "m1", HomeTeam: "Columbus Crew", AwayTeam: "FC Dallas", Kickoff: time.Now().Add(-1 * time.Hour), Status: "STATUS_FULL_TIME", State: "post"}})

	if len(delays) != 0 {
		t.Errorf("expected no timer for terminal match, got %d", len(delays))
	}
}

func TestMatchPoller_ScheduleSkipsAlreadyScheduled(t *testing.T) {
	var delays []time.Duration
	p := newPoller(repository.NewMemoryMatchStore(), repository.NewMemoryResultStore(), nil, capturingTimer(&delays))

	m := models.Match{ID: "m1", HomeTeam: "Columbus Crew", AwayTeam: "FC Dallas", Kickoff: time.Now().Add(1 * time.Hour), Status: "STATUS_SCHEDULED"}
	p.Schedule([]models.Match{m})
	p.Schedule([]models.Match{m})

	if len(delays) != 1 {
		t.Errorf("expected 1 timer after two Schedule calls for same match, got %d", len(delays))
	}
}

func TestMatchPoller_Tick_NoopWhenNoActiveMatches(t *testing.T) {
	calls := 0
	matchStore := repository.NewMemoryMatchStore()
	resultStore := repository.NewMemoryResultStore()
	fetcher := func() ([]models.Match, error) { calls++; return nil, nil }
	p := poll.NewMatchPoller(matchStore, resultStore, fetcher, capturingTimer(new([]time.Duration)))

	p.Tick(context.Background())

	if calls != 0 {
		t.Errorf("expected no fetcher calls with no active matches, got %d", calls)
	}
}

func TestMatchPoller_Tick_SavesResultAndDeactivatesTerminalMatch(t *testing.T) {
	matchStore := repository.NewMemoryMatchStore()
	resultStore := repository.NewMemoryResultStore()

	// Fetcher simulates ESPN returning a now-terminal match
	terminalMatch := models.Match{
		ID: "m-done", HomeTeam: "Columbus Crew", AwayTeam: "FC Dallas",
		Status: "STATUS_FULL_TIME", State: "post", HomeScore: "2", AwayScore: "1",
		Kickoff: time.Now().Add(-2 * time.Hour),
	}
	p := newPoller(matchStore, resultStore, []models.Match{terminalMatch}, immediateTimer)

	// Schedule with pre-match state (as it would have been at kickoff)
	preMatch := models.Match{
		ID: "m-done", HomeTeam: "Columbus Crew", AwayTeam: "FC Dallas",
		Status: "STATUS_SCHEDULED", State: "pre",
		Kickoff: time.Now().Add(-2 * time.Hour),
	}
	p.Schedule([]models.Match{preMatch})

	p.Tick(context.Background())

	result, err := resultStore.GetResult(context.Background(), "m-done")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || result.HomeGoals != 2 || result.AwayGoals != 1 {
		t.Errorf("expected result 2-1, got %+v", result)
	}
	if p.ActiveCount() != 0 {
		t.Errorf("expected match removed from active after terminal, got %d active", p.ActiveCount())
	}
}

func TestMatchPoller_Tick_KeepsActiveForLiveMatch(t *testing.T) {
	matchStore := repository.NewMemoryMatchStore()
	resultStore := repository.NewMemoryResultStore()

	liveMatch := models.Match{
		ID: "m-live", HomeTeam: "Columbus Crew", AwayTeam: "FC Dallas",
		Status: "STATUS_IN_PROGRESS", State: "in", HomeScore: "1", AwayScore: "0",
		Kickoff: time.Now().Add(-1 * time.Hour),
	}
	p := newPoller(matchStore, resultStore, []models.Match{liveMatch}, immediateTimer)
	p.Schedule([]models.Match{liveMatch})

	p.Tick(context.Background())

	if p.ActiveCount() != 1 {
		t.Errorf("expected match to remain active for live match, got %d active", p.ActiveCount())
	}
}

func TestMatchPoller_Tick_KeepsActiveForUnknownStatus(t *testing.T) {
	matchStore := repository.NewMemoryMatchStore()
	resultStore := repository.NewMemoryResultStore()

	unknownMatch := models.Match{
		ID: "m-weird", HomeTeam: "Columbus Crew", AwayTeam: "FC Dallas",
		Status: "STATUS_POSTPONED", State: "pre",
		Kickoff: time.Now().Add(-1 * time.Hour),
	}
	p := newPoller(matchStore, resultStore, []models.Match{unknownMatch}, immediateTimer)
	p.Schedule([]models.Match{unknownMatch})

	p.Tick(context.Background())

	if p.ActiveCount() != 1 {
		t.Errorf("expected match to remain active for unknown status (4am reset will clean up), got %d active", p.ActiveCount())
	}
}

func TestMatchPoller_Tick_NoopWhenFetcherFails(t *testing.T) {
	matchStore := repository.NewMemoryMatchStore()
	resultStore := repository.NewMemoryResultStore()

	callCount := 0
	fetcher := func() ([]models.Match, error) {
		callCount++
		return nil, fmt.Errorf("ESPN down")
	}
	p := poll.NewMatchPoller(matchStore, resultStore, fetcher, immediateTimer)

	pastMatch := models.Match{
		ID: "m1", HomeTeam: "Columbus Crew", AwayTeam: "FC Dallas",
		Status: "STATUS_SCHEDULED", State: "in",
		Kickoff: time.Now().Add(-1 * time.Hour),
	}
	p.Schedule([]models.Match{pastMatch})
	p.Tick(context.Background())

	result, _ := resultStore.GetResult(context.Background(), "m1")
	if result != nil {
		t.Errorf("expected no result saved when fetcher fails, got %+v", result)
	}
}

func TestMatchPoller_Run_StopsOnContextCancel(t *testing.T) {
	matchStore := repository.NewMemoryMatchStore()
	resultStore := repository.NewMemoryResultStore()
	fetcher := func() ([]models.Match, error) { return nil, nil }
	p := poll.NewMatchPoller(matchStore, resultStore, fetcher, func(time.Duration, func()) {})

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		p.Run(ctx, 10*time.Millisecond)
		close(done)
	}()

	cancel()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Error("Run did not stop after context cancel within 1 second")
	}
}

func TestMatchPoller_Backfill_SavesResultForTerminalMatch(t *testing.T) {
	resultStore := repository.NewMemoryResultStore()
	p := newPoller(repository.NewMemoryMatchStore(), resultStore, nil, capturingTimer(new([]time.Duration)))

	p.Backfill(context.Background(), []models.Match{
		{ID: "m-done", HomeTeam: "Columbus Crew", AwayTeam: "LA Galaxy",
			Status: "STATUS_FULL_TIME", HomeScore: "2", AwayScore: "1"},
	})

	result, err := resultStore.GetResult(context.Background(), "m-done")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || result.HomeGoals != 2 || result.AwayGoals != 1 {
		t.Errorf("expected result 2-1, got %+v", result)
	}
}

func TestMatchPoller_Backfill_SkipsNonTerminalMatch(t *testing.T) {
	resultStore := repository.NewMemoryResultStore()
	p := newPoller(repository.NewMemoryMatchStore(), resultStore, nil, capturingTimer(new([]time.Duration)))

	p.Backfill(context.Background(), []models.Match{
		{ID: "m-live", HomeTeam: "Columbus Crew", AwayTeam: "LA Galaxy",
			Status: "STATUS_IN_PROGRESS", State: "in", HomeScore: "1", AwayScore: "0"},
	})

	result, _ := resultStore.GetResult(context.Background(), "m-live")
	if result != nil {
		t.Errorf("expected no result for non-terminal match, got %+v", result)
	}
}

func TestMatchPoller_Backfill_LogsErrorWhenSaveResultFails(t *testing.T) {
	p := newPoller(repository.NewMemoryMatchStore(), repository.NewErrorResultStore(), nil, capturingTimer(new([]time.Duration)))

	// Should not panic — error is logged
	p.Backfill(context.Background(), []models.Match{
		{ID: "m1", HomeTeam: "Columbus Crew", AwayTeam: "FC Dallas", Status: "STATUS_FULL_TIME", HomeScore: "2", AwayScore: "1"},
	})
}

func TestMatchPoller_Tick_LogsErrorWhenSaveAllFails(t *testing.T) {
	resultStore := repository.NewMemoryResultStore()
	fetcher := func() ([]models.Match, error) {
		return []models.Match{{ID: "m1", HomeTeam: "Columbus Crew", AwayTeam: "FC Dallas", Status: "STATUS_IN_PROGRESS", HomeScore: "1", AwayScore: "0"}}, nil
	}
	p := poll.NewMatchPoller(&errMatchStore{}, resultStore, fetcher, immediateTimer)

	pastMatch := models.Match{ID: "m1", HomeTeam: "Columbus Crew", AwayTeam: "FC Dallas", Status: "STATUS_SCHEDULED", Kickoff: time.Now().Add(-1 * time.Hour)}
	p.Schedule([]models.Match{pastMatch})

	// Should not panic — error is logged
	p.Tick(context.Background())
}

func TestMatchPoller_Tick_LogsErrorWhenSaveResultFails(t *testing.T) {
	fetcher := func() ([]models.Match, error) {
		return []models.Match{{ID: "m1", HomeTeam: "Columbus Crew", AwayTeam: "FC Dallas", Status: "STATUS_FULL_TIME", HomeScore: "2", AwayScore: "1", Kickoff: time.Now().Add(-2 * time.Hour)}}, nil
	}
	p := poll.NewMatchPoller(repository.NewMemoryMatchStore(), repository.NewErrorResultStore(), fetcher, immediateTimer)

	preMatch := models.Match{ID: "m1", HomeTeam: "Columbus Crew", AwayTeam: "FC Dallas", Status: "STATUS_SCHEDULED", Kickoff: time.Now().Add(-2 * time.Hour)}
	p.Schedule([]models.Match{preMatch})

	// Should not panic — error is logged
	p.Tick(context.Background())
}

func TestMatchPoller_Tick_CallsOnResultSavedWhenMatchFinishes(t *testing.T) {
	matchStore := repository.NewMemoryMatchStore()
	resultStore := repository.NewMemoryResultStore()

	terminalMatch := models.Match{
		ID: "m-done", HomeTeam: "Columbus Crew", AwayTeam: "FC Dallas",
		Status: "STATUS_FULL_TIME", State: "post", HomeScore: "2", AwayScore: "1",
		Kickoff: time.Now().Add(-2 * time.Hour),
	}
	p := newPoller(matchStore, resultStore, []models.Match{terminalMatch}, immediateTimer)
	p.Schedule([]models.Match{{ID: "m-done", HomeTeam: "Columbus Crew", AwayTeam: "FC Dallas", Status: "STATUS_SCHEDULED", Kickoff: time.Now().Add(-2 * time.Hour)}})

	called := 0
	p.SetOnResultSaved(func(_ context.Context) { called++ })

	p.Tick(context.Background())

	if called != 1 {
		t.Errorf("expected onResultSaved called once, got %d", called)
	}
}

func TestMatchPoller_Tick_DoesNotCallOnResultSavedWhenMatchStillLive(t *testing.T) {
	matchStore := repository.NewMemoryMatchStore()
	resultStore := repository.NewMemoryResultStore()

	liveMatch := models.Match{
		ID: "m-live", HomeTeam: "Columbus Crew", AwayTeam: "FC Dallas",
		Status: "STATUS_IN_PROGRESS", State: "in", HomeScore: "1", AwayScore: "0",
		Kickoff: time.Now().Add(-1 * time.Hour),
	}
	p := newPoller(matchStore, resultStore, []models.Match{liveMatch}, immediateTimer)
	p.Schedule([]models.Match{liveMatch})

	called := 0
	p.SetOnResultSaved(func(_ context.Context) { called++ })

	p.Tick(context.Background())

	if called != 0 {
		t.Errorf("expected onResultSaved not called for live match, got %d", called)
	}
}

func TestMatchPoller_Reset_ClearsActiveAndReschedules(t *testing.T) {
	matchStore := repository.NewMemoryMatchStore()
	resultStore := repository.NewMemoryResultStore()

	liveMatch := models.Match{
		ID: "m-live", HomeTeam: "Columbus Crew", AwayTeam: "FC Dallas",
		Status: "STATUS_IN_PROGRESS", State: "in",
		Kickoff: time.Now().Add(-1 * time.Hour),
	}

	var delays []time.Duration
	p := newPoller(matchStore, resultStore, []models.Match{liveMatch}, immediateTimer)
	p.Schedule([]models.Match{liveMatch})

	if p.ActiveCount() != 1 {
		t.Fatalf("expected 1 active match before reset, got %d", p.ActiveCount())
	}

	// After reset, old active cleared; new schedule applied with capturing timer
	newMatch := models.Match{
		ID: "m-new", HomeTeam: "Columbus Crew", AwayTeam: "Portland Timbers",
		Status: "STATUS_SCHEDULED", State: "pre",
		Kickoff: time.Now().Add(3 * time.Hour),
	}
	p.SetTimerFunc(capturingTimer(&delays))
	p.Reset([]models.Match{newMatch})

	if p.ActiveCount() != 0 {
		t.Errorf("expected active cleared after reset, got %d", p.ActiveCount())
	}
	if len(delays) != 1 {
		t.Errorf("expected 1 new timer after reset, got %d", len(delays))
	}
}

func TestMatchPoller_Run_ExitsOnContextCancel(t *testing.T) {
	p := newPoller(repository.NewMemoryMatchStore(), repository.NewMemoryResultStore(), nil, capturingTimer(new([]time.Duration)))
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		p.Run(ctx, 10*time.Millisecond)
		close(done)
	}()
	cancel()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Error("Run did not exit after context cancellation")
	}
}
