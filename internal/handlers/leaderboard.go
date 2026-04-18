package handlers

import (
	"net/http"
	"sort"

	"github.com/mcornell/crew-predictions/internal/repository"
	"github.com/mcornell/crew-predictions/internal/scoring"
	"github.com/mcornell/crew-predictions/templates"
)

type LeaderboardHandler struct {
	predictions repository.PredictionStore
	results     repository.ResultStore
}

func NewLeaderboardHandler(predictions repository.PredictionStore, results repository.ResultStore) *LeaderboardHandler {
	return &LeaderboardHandler{predictions: predictions, results: results}
}

func (h *LeaderboardHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	allPredictions, err := h.predictions.GetAll(ctx)
	if err != nil {
		http.Error(w, "couldn't load predictions", http.StatusInternalServerError)
		return
	}

	totals := map[string]int{}
	for _, p := range allPredictions {
		result, err := h.results.GetResult(ctx, p.MatchID)
		if err != nil || result == nil {
			continue
		}
		pts := scoring.AcesRadio(
			scoring.Result{Home: result.HomeGoals, Away: result.AwayGoals},
			scoring.Prediction{Home: p.HomeGoals, Away: p.AwayGoals},
		)
		totals[p.Handle] += pts
	}

	entries := make([]templates.LeaderboardEntry, 0, len(totals))
	for handle, pts := range totals {
		entries = append(entries, templates.LeaderboardEntry{Handle: handle, Points: pts})
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Points > entries[j].Points })

	handle := userFromSession(r)
	templates.Leaderboard(entries, handle).Render(ctx, w)
}
