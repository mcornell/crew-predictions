package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/mcornell/crew-predictions/internal/handlers"
)

func main() {
	godotenv.Load()
	setupLogging()

	cfg := loadConfig()
	ctx := context.Background()

	stores, err := buildStores(ctx, cfg)
	if err != nil {
		log.Fatalf("buildStores: %v", err)
	}

	verifier := buildVerifier(ctx, cfg)
	deps := buildDeps(cfg, stores, verifier)

	if cfg.SessionSecret != "" {
		handlers.SetSessionSecret([]byte(cfg.SessionSecret))
	}

	mux := http.NewServeMux()
	registerRoutes(mux, deps)
	if cfg.TestMode {
		registerTestRoutes(mux, deps)
	}

	if !cfg.TestMode {
		etLoc, err := time.LoadLocation("America/New_York")
		if err != nil {
			log.Fatalf("failed to load ET timezone: %v", err)
		}
		stop := startDailyRefresh(stores.Match, deps.RefreshFetcher, deps.MatchPoller, deps.TwoOneBot, deps.RecalcFn, etLoc)
		defer close(stop)

		pollerCtx, cancelPoller := context.WithCancel(ctx)
		defer cancelPoller()
		go deps.MatchPoller.Run(pollerCtx, 2*time.Minute)
	}

	slog.Info("server listening", "port", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		log.Fatal(err)
	}
}

// setupLogging configures slog to emit JSON that Cloud Logging understands:
// the "level" field becomes "severity" (with WARN remapped to Google's
// "WARNING" spelling), and "msg" becomes "message".
func setupLogging() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == slog.LevelKey {
				a.Key = "severity"
				if a.Value.String() == "WARN" {
					a.Value = slog.StringValue("WARNING")
				}
			}
			if a.Key == slog.MessageKey {
				a.Key = "message"
			}
			return a
		},
	})))
}
