package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/mcornell/crew-predictions/internal/models"
	"github.com/mcornell/crew-predictions/internal/repository"
)

type MatchesHandler struct {
	store   repository.PredictionStore
	fetcher func() ([]models.Match, error)
}

func NewMatchesHandler(store repository.PredictionStore, fetcher func() ([]models.Match, error)) *MatchesHandler {
	return &MatchesHandler{store: store, fetcher: fetcher}
}

type apiMatch struct {
	ID        string `json:"id"`
	HomeTeam  string `json:"homeTeam"`
	AwayTeam  string `json:"awayTeam"`
	Kickoff   string `json:"kickoff"`
	Status    string `json:"status"`
	HomeScore string `json:"homeScore"`
	AwayScore string `json:"awayScore"`
}

type apiPrediction struct {
	HomeGoals int `json:"homeGoals"`
	AwayGoals int `json:"awayGoals"`
}

func (h *MatchesHandler) APIList(w http.ResponseWriter, r *http.Request) {
	matches, err := h.fetcher()
	if err != nil {
		http.Error(w, "couldn't fetch matches", http.StatusInternalServerError)
		return
	}

	out := make([]apiMatch, len(matches))
	for i, m := range matches {
		out[i] = apiMatch{
			ID:        m.ID,
			HomeTeam:  m.HomeTeam,
			AwayTeam:  m.AwayTeam,
			Kickoff:   m.Kickoff.Format("2006-01-02T15:04:05Z"),
			Status:    m.Status,
			HomeScore: m.HomeScore,
			AwayScore: m.AwayScore,
		}
	}

	preds := map[string]apiPrediction{}
	if user := UserFromSession(r); user != nil {
		for _, m := range matches {
			p, _ := h.store.GetByMatchAndUser(r.Context(), m.ID, user.UserID)
			if p != nil {
				preds[m.ID] = apiPrediction{HomeGoals: p.HomeGoals, AwayGoals: p.AwayGoals}
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"matches": out, "predictions": preds})
}
