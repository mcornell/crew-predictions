package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mcornell/crew-predictions/internal/bot"
	"github.com/mcornell/crew-predictions/internal/models"
	"github.com/mcornell/crew-predictions/internal/poll"
	"github.com/mcornell/crew-predictions/internal/repository"
)

func TestNext4amET_BeforeFourAM_ReturnsTodayAt4AM(t *testing.T) {
	etLoc, _ := time.LoadLocation("America/New_York")
	now := time.Date(2026, 4, 22, 2, 0, 0, 0, etLoc) // 2am ET
	next := next4amET(now, etLoc)
	want := time.Date(2026, 4, 22, 4, 0, 0, 0, etLoc)
	if !next.Equal(want) {
		t.Errorf("expected %v, got %v", want, next)
	}
}

func TestNext4amET_AfterFourAM_ReturnsTomorrowAt4AM(t *testing.T) {
	etLoc, _ := time.LoadLocation("America/New_York")
	now := time.Date(2026, 4, 22, 5, 0, 0, 0, etLoc) // 5am ET
	next := next4amET(now, etLoc)
	want := time.Date(2026, 4, 23, 4, 0, 0, 0, etLoc)
	if !next.Equal(want) {
		t.Errorf("expected %v, got %v", want, next)
	}
}

func TestNext4amET_ExactlyFourAM_ReturnsTomorrowAt4AM(t *testing.T) {
	etLoc, _ := time.LoadLocation("America/New_York")
	now := time.Date(2026, 4, 22, 4, 0, 0, 0, etLoc)
	next := next4amET(now, etLoc)
	want := time.Date(2026, 4, 23, 4, 0, 0, 0, etLoc)
	if !next.Equal(want) {
		t.Errorf("expected %v, got %v", want, next)
	}
}

func TestStartDailyRefresh_PopulatesStoreImmediately(t *testing.T) {
	etLoc, _ := time.LoadLocation("America/New_York")
	store := repository.NewMemoryMatchStore()
	called := make(chan struct{}, 1)
	fetcher := func() ([]models.Match, error) {
		called <- struct{}{}
		return []models.Match{{ID: "bg-match", Kickoff: time.Now().Add(24 * time.Hour)}}, nil
	}
	poller := poll.NewMatchPoller(store, repository.NewMemoryResultStore(), fetcher,
		func(time.Duration, func()) {})
	noopBot := bot.New(repository.NewMemoryPredictionStore(), repository.NewMemoryUserStore(), "Columbus Crew")

	stop := startDailyRefresh(store, fetcher, poller, noopBot, etLoc)
	defer close(stop)

	select {
	case <-called:
	case <-time.After(time.Second):
		t.Fatal("fetcher was not called within 1 second of startup")
	}

	matches, _ := store.GetAll()
	if len(matches) != 1 || matches[0].ID != "bg-match" {
		t.Errorf("expected bg-match in store, got %+v", matches)
	}
}

func TestStartDailyRefresh_BackfillsCompletedMatchResults(t *testing.T) {
	etLoc, _ := time.LoadLocation("America/New_York")
	matchStore := repository.NewMemoryMatchStore()
	resultStore := repository.NewMemoryResultStore()
	completedMatch := models.Match{
		ID: "m-finished", HomeTeam: "Columbus Crew", AwayTeam: "LA Galaxy",
		Status: "STATUS_FULL_TIME", HomeScore: "2", AwayScore: "1",
		Kickoff: time.Now().Add(-3 * time.Hour),
	}
	fetcher := func() ([]models.Match, error) { return []models.Match{completedMatch}, nil }
	poller := poll.NewMatchPoller(matchStore, resultStore, fetcher, func(time.Duration, func()) {})
	noopBot := bot.New(repository.NewMemoryPredictionStore(), repository.NewMemoryUserStore(), "Columbus Crew")

	stop := startDailyRefresh(matchStore, fetcher, poller, noopBot, etLoc)
	defer close(stop)

	// startDailyRefresh calls refresh() immediately on startup; give it a moment
	time.Sleep(50 * time.Millisecond)

	result, err := resultStore.GetResult(context.Background(), "m-finished")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected backfill to save result for completed match, got nil")
	}
	if result.HomeGoals != 2 || result.AwayGoals != 1 {
		t.Errorf("expected 2-1, got %d-%d", result.HomeGoals, result.AwayGoals)
	}
}

func TestStartDailyRefresh_BotPredictsUpcomingMatches(t *testing.T) {
	etLoc, _ := time.LoadLocation("America/New_York")
	matchStore := repository.NewMemoryMatchStore()
	resultStore := repository.NewMemoryResultStore()
	predStore := repository.NewMemoryPredictionStore()
	userStore := repository.NewMemoryUserStore()
	upcomingMatch := models.Match{
		ID: "m-upcoming", HomeTeam: "Columbus Crew", AwayTeam: "FC Dallas",
		Status: "STATUS_SCHEDULED", Kickoff: time.Now().Add(24 * time.Hour),
	}
	fetcher := func() ([]models.Match, error) { return []models.Match{upcomingMatch}, nil }
	poller := poll.NewMatchPoller(matchStore, resultStore, fetcher, func(time.Duration, func()) {})
	twoOneBot := bot.New(predStore, userStore, "Columbus Crew")

	stop := startDailyRefresh(matchStore, fetcher, poller, twoOneBot, etLoc)
	defer close(stop)

	time.Sleep(50 * time.Millisecond)

	p, err := predStore.GetByMatchAndUser(context.Background(), "m-upcoming", bot.UserID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected bot prediction for upcoming match, got nil")
	}
	if p.HomeGoals != 2 || p.AwayGoals != 1 {
		t.Errorf("expected 2-1, got %d-%d", p.HomeGoals, p.AwayGoals)
	}
}

func TestServeFirebaseConfig_ReturnsJavaScriptWithEnvVars(t *testing.T) {
	t.Setenv("FIREBASE_API_KEY", "test-api-key")
	t.Setenv("FIREBASE_PROJECT_ID", "test-project")
	t.Setenv("FIREBASE_AUTH_DOMAIN", "test-project.firebaseapp.com")

	req := httptest.NewRequest(http.MethodGet, "/auth/config.js", nil)
	rr := httptest.NewRecorder()
	serveFirebaseConfig(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "application/javascript" {
		t.Errorf("expected Content-Type application/javascript, got %q", ct)
	}
	body := rr.Body.String()
	if !strings.Contains(body, "window.__firebaseConfig") {
		t.Errorf("expected window.__firebaseConfig assignment, got %q", body)
	}
	if !strings.Contains(body, "test-api-key") {
		t.Errorf("expected API key in response, got %q", body)
	}
	if !strings.Contains(body, "test-project") {
		t.Errorf("expected project ID in response, got %q", body)
	}
}

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
