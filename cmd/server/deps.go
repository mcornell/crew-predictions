package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/mcornell/crew-predictions/internal/bot"
	"github.com/mcornell/crew-predictions/internal/espn"
	"github.com/mcornell/crew-predictions/internal/handlers"
	"github.com/mcornell/crew-predictions/internal/models"
	"github.com/mcornell/crew-predictions/internal/poll"
	"github.com/mcornell/crew-predictions/internal/recalculator"
)

// Deps is the bag of dependencies passed into route registration and
// background-job startup. Everything in here is constructed at server
// startup and lives for the process lifetime.
type Deps struct {
	Cfg            Config
	Stores         Stores
	Verifier       handlers.TokenVerifier
	SummaryFetcher func(matchID string) (models.MatchSummary, error)
	RefreshFetcher func() ([]models.Match, error)
	RecalcFn       func(ctx context.Context)
	TwoOneBot      *bot.TwoOneBot
	MatchPoller    *poll.MatchPoller // nil in test mode (no background polling)
}

// buildDeps wires the runtime dependencies that bridge stores and handlers:
// the summary fetcher (live ESPN vs. fixture-backed in test mode), the
// match-refresh fetcher (live ESPN vs. read-from-store in test mode), the
// recalculation closure that handlers invoke after persisting state, the
// 2-1 bot, and the background match poller (production only).
func buildDeps(cfg Config, stores Stores, verifier handlers.TokenVerifier) Deps {
	summaryFetcher := chooseSummaryFetcher(cfg)
	refreshFetcher := chooseRefreshFetcher(cfg, stores.Match)

	recalcFn := func(ctx context.Context) {
		if err := recalculator.Recalculate(ctx, stores.Prediction, stores.Result, stores.User, nil, recalculator.SeasonWindow{}, cfg.TargetTeam); err != nil {
			slog.Error("recalculate failed", "error", err)
		}
	}

	twoOneBot := bot.New(stores.Prediction, stores.User, cfg.TargetTeam)

	deps := Deps{
		Cfg:            cfg,
		Stores:         stores,
		Verifier:       verifier,
		SummaryFetcher: summaryFetcher,
		RefreshFetcher: refreshFetcher,
		RecalcFn:       recalcFn,
		TwoOneBot:      twoOneBot,
	}

	if !cfg.TestMode {
		deps.MatchPoller = poll.NewMatchPoller(
			stores.Match, stores.Result, refreshFetcher,
			func(d time.Duration, f func()) { time.AfterFunc(d, f) },
		)
		deps.MatchPoller.SetSummaryFetcher(summaryFetcher)
		deps.MatchPoller.SetOnResultSaved(recalcFn)
	}

	return deps
}

// chooseSummaryFetcher returns the live ESPN summary fetcher in production
// or a fixture-backed one in TEST_MODE so e2e tests don't depend on ESPN.
func chooseSummaryFetcher(cfg Config) func(matchID string) (models.MatchSummary, error) {
	if cfg.TestMode {
		return espn.FixtureFetcher("internal/espn/testdata")
	}
	return espn.FetchSummary
}

// chooseRefreshFetcher returns the live ESPN scoreboard fetcher in production
// or a no-op (read-from-store) in TEST_MODE so the seed endpoints fully
// control match data.
func chooseRefreshFetcher(cfg Config, matchStore interface {
	GetAll() ([]models.Match, error)
}) func() ([]models.Match, error) {
	if cfg.TestMode {
		return matchStore.GetAll
	}
	return espn.FetchCrewMatches
}
