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
	targetTeam  string
}

func NewLeaderboardHandler(predictions repository.PredictionStore, results repository.ResultStore, targetTeam string) *LeaderboardHandler {
	return &LeaderboardHandler{predictions: predictions, results: results, targetTeam: targetTeam}
}

func (h *LeaderboardHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	allPredictions, err := h.predictions.GetAll(ctx)
	if err != nil {
		http.Error(w, "couldn't load predictions", http.StatusInternalServerError)
		return
	}

	acesTotals := map[string]int{}
	u90Totals := map[string]int{}

	for _, p := range allPredictions {
		result, err := h.results.GetResult(ctx, p.MatchID)
		if err != nil || result == nil {
			continue
		}
		pred := scoring.Prediction{Home: p.HomeGoals, Away: p.AwayGoals}
		res := scoring.Result{Home: result.HomeGoals, Away: result.AwayGoals}
		targetIsHome := result.HomeTeam == h.targetTeam

		acesTotals[p.Handle] += scoring.AcesRadio(res, pred)
		u90Totals[p.Handle] += scoring.Upper90Club(res, pred, targetIsHome)
	}

	handles := allHandles(acesTotals, u90Totals)

	acesEntries := rankEntries(handles, acesTotals)
	u90Entries := rankEntries(handles, u90Totals)

	handle := userFromSession(r)
	templates.Leaderboard(acesEntries, u90Entries, handle).Render(ctx, w)
}

func allHandles(maps ...map[string]int) []string {
	seen := map[string]struct{}{}
	for _, m := range maps {
		for k := range m {
			seen[k] = struct{}{}
		}
	}
	handles := make([]string, 0, len(seen))
	for k := range seen {
		handles = append(handles, k)
	}
	return handles
}

func rankEntries(handles []string, totals map[string]int) []templates.LeaderboardEntry {
	entries := make([]templates.LeaderboardEntry, 0, len(handles))
	for _, h := range handles {
		entries = append(entries, templates.LeaderboardEntry{Handle: h, Points: totals[h]})
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Points > entries[j].Points })
	return entries
}
