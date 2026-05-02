package main

import (
	"reflect"
	"testing"
	"time"

	"github.com/mcornell/crew-predictions/internal/espn"
	"github.com/mcornell/crew-predictions/internal/handlers"
	"github.com/mcornell/crew-predictions/internal/models"
	"github.com/mcornell/crew-predictions/internal/repository"
)

// fnPtr returns the underlying code pointer of a function value so we can
// assert "this is the same function" without invoking it. Useful for
// distinguishing the live ESPN fetcher from the fixture-backed one.
func fnPtr(f any) uintptr {
	return reflect.ValueOf(f).Pointer()
}

func TestChooseSummaryFetcher_ProductionReturnsLiveFetcher(t *testing.T) {
	got := chooseSummaryFetcher(Config{TestMode: false})
	if fnPtr(got) != fnPtr(espn.FetchSummary) {
		t.Errorf("expected espn.FetchSummary in production mode, got a different function")
	}
}

func TestChooseSummaryFetcher_TestModeReturnsFixtureBacked(t *testing.T) {
	got := chooseSummaryFetcher(Config{TestMode: true})
	if got == nil {
		t.Fatal("expected non-nil fetcher")
	}
	// FixtureFetcher returns a closure, not espn.FetchSummary directly.
	if fnPtr(got) == fnPtr(espn.FetchSummary) {
		t.Error("expected fixture-backed fetcher in test mode, got live FetchSummary")
	}
}

func TestChooseRefreshFetcher_ProductionReturnsLiveFetcher(t *testing.T) {
	got := chooseRefreshFetcher(Config{TestMode: false}, repository.NewMemoryMatchStore())
	if fnPtr(got) != fnPtr(espn.FetchCrewMatches) {
		t.Errorf("expected espn.FetchCrewMatches in production mode, got a different function")
	}
}

func TestChooseRefreshFetcher_TestModeReadsFromStore(t *testing.T) {
	store := repository.NewMemoryMatchStore()
	seeded := []models.Match{{ID: "seeded-1", Kickoff: time.Now()}}
	if err := store.SaveAll(seeded); err != nil {
		t.Fatalf("seed store: %v", err)
	}
	fetcher := chooseRefreshFetcher(Config{TestMode: true}, store)
	got, err := fetcher()
	if err != nil {
		t.Fatalf("fetcher() error: %v", err)
	}
	if len(got) != 1 || got[0].ID != "seeded-1" {
		t.Errorf("expected refresh fetcher to read from store, got %+v", got)
	}
}

func TestBuildDeps_TestModeOmitsMatchPoller(t *testing.T) {
	cfg := Config{TestMode: true, TargetTeam: "Columbus Crew"}
	stores, err := buildStores(t.Context(), cfg)
	if err != nil {
		t.Fatalf("buildStores: %v", err)
	}
	deps := buildDeps(cfg, stores, handlers.NoopTokenVerifier{})

	if deps.MatchPoller != nil {
		t.Errorf("expected MatchPoller nil in test mode, got %T", deps.MatchPoller)
	}
	if deps.SummaryFetcher == nil {
		t.Error("expected SummaryFetcher populated")
	}
	if deps.RefreshFetcher == nil {
		t.Error("expected RefreshFetcher populated")
	}
	if deps.RecalcFn == nil {
		t.Error("expected RecalcFn populated")
	}
	if deps.TwoOneBot == nil {
		t.Error("expected TwoOneBot populated")
	}
}

func TestBuildDeps_ProductionPopulatesMatchPoller(t *testing.T) {
	cfg := Config{TestMode: false, FirestoreProject: "", TargetTeam: "Columbus Crew"}
	stores, err := buildStores(t.Context(), cfg)
	if err != nil {
		t.Fatalf("buildStores: %v", err)
	}
	deps := buildDeps(cfg, stores, handlers.NoopTokenVerifier{})

	if deps.MatchPoller == nil {
		t.Error("expected MatchPoller populated outside test mode")
	}
}

func TestBuildDeps_PassesThroughCfgStoresVerifier(t *testing.T) {
	cfg := Config{TestMode: true, TargetTeam: "Columbus Crew", AdminKey: "marker"}
	stores, err := buildStores(t.Context(), cfg)
	if err != nil {
		t.Fatalf("buildStores: %v", err)
	}
	verifier := handlers.NoopTokenVerifier{}
	deps := buildDeps(cfg, stores, verifier)

	if deps.Cfg.AdminKey != "marker" {
		t.Errorf("Cfg not passed through: AdminKey=%q", deps.Cfg.AdminKey)
	}
	if deps.Stores.Match == nil {
		t.Error("Stores not passed through: Match nil")
	}
	if _, ok := deps.Verifier.(handlers.NoopTokenVerifier); !ok {
		t.Errorf("Verifier not passed through: got %T", deps.Verifier)
	}
}

func TestBuildDeps_RecalcFnDoesNotPanicOnEmptyStores(t *testing.T) {
	cfg := Config{TestMode: true, TargetTeam: "Columbus Crew"}
	stores, err := buildStores(t.Context(), cfg)
	if err != nil {
		t.Fatalf("buildStores: %v", err)
	}
	deps := buildDeps(cfg, stores, handlers.NoopTokenVerifier{})
	deps.RecalcFn(t.Context()) // empty stores: should be a no-op, not a panic
}
