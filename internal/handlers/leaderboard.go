package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"

	"github.com/mcornell/crew-predictions/internal/repository"
	"github.com/mcornell/crew-predictions/internal/scoring"
)

type LeaderboardHandler struct {
	predictions repository.PredictionStore
	results     repository.ResultStore
	users       repository.UserStore
	targetTeam  string
}

func NewLeaderboardHandler(predictions repository.PredictionStore, results repository.ResultStore, users repository.UserStore, targetTeam string) *LeaderboardHandler {
	return &LeaderboardHandler{predictions: predictions, results: results, users: users, targetTeam: targetTeam}
}

type leaderboardEntry struct {
	UserID     string `json:"userID"`
	Handle     string `json:"handle"`
	Points     int    `json:"points"`
	HasProfile bool   `json:"hasProfile"`
}

func (h *LeaderboardHandler) APIList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	allPredictions, err := h.predictions.GetAll(ctx)
	if err != nil {
		http.Error(w, "couldn't load predictions", http.StatusInternalServerError)
		return
	}

	// Build userID → current handle map from UserStore; fall back to prediction handle when absent.
	allUsers, _ := h.users.GetAll(ctx)
	handleByUserID := make(map[string]string, len(allUsers))
	knownUsers := make(map[string]bool, len(allUsers))
	for _, u := range allUsers {
		handleByUserID[u.UserID] = u.Handle
		knownUsers[u.UserID] = true
	}

	// Seed totals from all predictions so users with ≥1 prediction appear at 0 pts before results land.
	acesTotals := map[string]int{}
	u90Totals := map[string]int{}
	keyHandle := map[string]string{} // key → display handle

	for _, p := range allPredictions {
		key := p.UserID
		if key == "" {
			key = p.Handle
		}
		if _, seen := keyHandle[key]; !seen {
			if h, ok := handleByUserID[key]; ok {
				keyHandle[key] = h
			} else {
				keyHandle[key] = p.Handle
			}
		}
		if _, seen := acesTotals[key]; !seen {
			acesTotals[key] = 0
			u90Totals[key] = 0
		}
		result, err := h.results.GetResult(ctx, p.MatchID)
		if err != nil || result == nil {
			continue
		}
		pred := scoring.Prediction{Home: p.HomeGoals, Away: p.AwayGoals}
		res := scoring.Result{Home: result.HomeGoals, Away: result.AwayGoals}
		targetIsHome := result.HomeTeam == h.targetTeam
		acesTotals[key] += scoring.AcesRadio(res, pred)
		u90Totals[key] += scoring.Upper90Club(res, pred, targetIsHome)
	}

	keys := allKeys(acesTotals, u90Totals)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]any{
		"acesRadio":   toEntries(keys, keyHandle, acesTotals, knownUsers),
		"upper90Club": toEntries(keys, keyHandle, u90Totals, knownUsers),
	}); err != nil {
		log.Printf("leaderboard: encode response: %v", err)
	}
}

func toEntries(keys []string, keyHandle map[string]string, totals map[string]int, knownUsers map[string]bool) []leaderboardEntry {
	entries := make([]leaderboardEntry, 0, len(keys))
	for _, k := range keys {
		entries = append(entries, leaderboardEntry{UserID: k, Handle: keyHandle[k], Points: totals[k], HasProfile: knownUsers[k]})
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Points > entries[j].Points })
	return entries
}

func allKeys(maps ...map[string]int) []string {
	seen := map[string]struct{}{}
	for _, m := range maps {
		for k := range m {
			seen[k] = struct{}{}
		}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	return keys
}
