package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"

	"github.com/mcornell/crew-predictions/internal/repository"
	"github.com/mcornell/crew-predictions/internal/scoring"
)

type ProfileHandler struct {
	predictions repository.PredictionStore
	results     repository.ResultStore
	users       repository.UserStore
	targetTeam  string
}

func NewProfileHandler(predictions repository.PredictionStore, results repository.ResultStore, users repository.UserStore, targetTeam string) *ProfileHandler {
	return &ProfileHandler{predictions: predictions, results: results, users: users, targetTeam: targetTeam}
}

type standingEntry struct {
	Points int `json:"points"`
	Rank   int `json:"rank"`
}

func (h *ProfileHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := r.PathValue("userID")

	user, err := h.users.GetByID(ctx, userID)
	if err != nil {
		http.Error(w, "could not load user", http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	allPredictions, err := h.predictions.GetAll(ctx)
	if err != nil {
		http.Error(w, "could not load predictions", http.StatusInternalServerError)
		return
	}

	// Count this user's predictions and compute full leaderboard for rank.
	predictionCount := 0
	acesTotals := map[string]int{}
	u90Totals := map[string]int{}
	grouchyTotals := map[string]int{}

	for _, p := range allPredictions {
		if p.UserID == userID {
			predictionCount++
		}
		result, err := h.results.GetResult(ctx, p.MatchID)
		if err != nil || result == nil {
			continue
		}
		key := p.UserID
		if key == "" {
			key = p.Handle
		}
		pred := scoring.Prediction{Home: p.HomeGoals, Away: p.AwayGoals}
		res := scoring.Result{Home: result.HomeGoals, Away: result.AwayGoals}
		targetIsHome := result.HomeTeam == h.targetTeam
		acesTotals[key] += scoring.AcesRadio(res, pred)
		u90Totals[key] += scoring.Upper90Club(res, pred, targetIsHome)
		grouchyTotals[key] += scoring.Grouchy(res, pred, targetIsHome)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]any{
		"userID":          user.UserID,
		"handle":          user.Handle,
		"location":        user.Location,
		"predictionCount": predictionCount,
		"acesRadio":       rankFor(userID, acesTotals),
		"upper90Club":     rankFor(userID, u90Totals),
		"grouchy":         rankFor(userID, grouchyTotals),
	}); err != nil {
		log.Printf("profile: encode response: %v", err)
	}
}

func rankFor(userID string, totals map[string]int) standingEntry {
	points, ok := totals[userID]
	if !ok {
		return standingEntry{Points: 0, Rank: 0}
	}
	type kv struct {
		key    string
		points int
	}
	sorted := make([]kv, 0, len(totals))
	for k, v := range totals {
		sorted = append(sorted, kv{k, v})
	}
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].points > sorted[j].points })
	rank := 1
	for i, entry := range sorted {
		if i > 0 && entry.points < sorted[i-1].points {
			rank = i + 1
		}
		if entry.key == userID {
			return standingEntry{Points: points, Rank: rank}
		}
	}
	return standingEntry{Points: points, Rank: 0}
}
