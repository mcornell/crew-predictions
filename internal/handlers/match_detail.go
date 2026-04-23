package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mcornell/crew-predictions/internal/repository"
	"github.com/mcornell/crew-predictions/internal/scoring"
)

type MatchDetailHandler struct {
	predictions repository.PredictionStore
	results     repository.ResultStore
	matches     repository.MatchStore
	targetTeam  string
}

func NewMatchDetailHandler(predictions repository.PredictionStore, results repository.ResultStore, matches repository.MatchStore, targetTeam string) *MatchDetailHandler {
	return &MatchDetailHandler{predictions: predictions, results: results, matches: matches, targetTeam: targetTeam}
}

type matchDetailPrediction struct {
	UserID            string `json:"userID"`
	Handle            string `json:"handle"`
	HomeGoals         int    `json:"homeGoals"`
	AwayGoals         int    `json:"awayGoals"`
	AcesRadioPoints   int    `json:"acesRadioPoints"`
	Upper90ClubPoints int    `json:"upper90ClubPoints"`
}

type matchDetailMatch struct {
	ID        string `json:"id"`
	HomeTeam  string `json:"homeTeam"`
	AwayTeam  string `json:"awayTeam"`
	Kickoff   string `json:"kickoff"`
	HomeScore string `json:"homeScore"`
	AwayScore string `json:"awayScore"`
}

func (h *MatchDetailHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	matchID := r.PathValue("matchId")

	allMatches, err := h.matches.GetAll()
	if err != nil {
		http.Error(w, "could not load matches", http.StatusInternalServerError)
		return
	}
	var found *matchDetailMatch
	for _, m := range allMatches {
		if m.ID == matchID {
			found = &matchDetailMatch{
				ID:        m.ID,
				HomeTeam:  m.HomeTeam,
				AwayTeam:  m.AwayTeam,
				Kickoff:   m.Kickoff.Format("2006-01-02T15:04:05Z07:00"),
				HomeScore: m.HomeScore,
				AwayScore: m.AwayScore,
			}
			break
		}
	}
	if found == nil {
		http.Error(w, "match not found", http.StatusNotFound)
		return
	}

	preds, err := h.predictions.GetByMatch(ctx, matchID)
	if err != nil {
		http.Error(w, "could not load predictions", http.StatusInternalServerError)
		return
	}

	result, _ := h.results.GetResult(ctx, matchID)

	entries := make([]matchDetailPrediction, 0, len(preds))
	for _, p := range preds {
		entry := matchDetailPrediction{
			UserID:    p.UserID,
			Handle:    p.Handle,
			HomeGoals: p.HomeGoals,
			AwayGoals: p.AwayGoals,
		}
		if result != nil {
			pred := scoring.Prediction{Home: p.HomeGoals, Away: p.AwayGoals}
			res := scoring.Result{Home: result.HomeGoals, Away: result.AwayGoals}
			targetIsHome := result.HomeTeam == h.targetTeam
			entry.AcesRadioPoints = scoring.AcesRadio(res, pred)
			entry.Upper90ClubPoints = scoring.Upper90Club(res, pred, targetIsHome)
		}
		entries = append(entries, entry)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]any{
		"match": found,
		"scoringFormats": []map[string]string{
			{"key": "acesRadio", "label": "Aces Radio"},
			{"key": "upper90Club", "label": "Upper 90 Club"},
		},
		"predictions": entries,
	}); err != nil {
		log.Printf("match_detail: encode response: %v", err)
	}
}
