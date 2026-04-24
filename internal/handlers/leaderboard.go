package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"

	"github.com/mcornell/crew-predictions/internal/repository"
)

type LeaderboardHandler struct {
	users repository.UserStore
}

func NewLeaderboardHandler(_ repository.PredictionStore, _ repository.ResultStore, users repository.UserStore, _ string) *LeaderboardHandler {
	return &LeaderboardHandler{users: users}
}

type leaderboardEntry struct {
	UserID          string `json:"userID"`
	Handle          string `json:"handle"`
	AcesRadioPoints int    `json:"acesRadioPoints"`
	Upper90Points   int    `json:"upper90ClubPoints"`
	GrouchyPoints   int    `json:"grouchyPoints"`
	HasProfile      bool   `json:"hasProfile"`
}

func (h *LeaderboardHandler) APIList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	allUsers, err := h.users.GetAll(ctx)
	if err != nil {
		http.Error(w, "couldn't load leaderboard", http.StatusInternalServerError)
		return
	}

	entries := make([]leaderboardEntry, 0, len(allUsers))
	for _, u := range allUsers {
		if u.PredictionCount == 0 {
			continue
		}
		entries = append(entries, leaderboardEntry{
			UserID:          u.UserID,
			Handle:          u.Handle,
			AcesRadioPoints: u.AcesRadioPoints,
			Upper90Points:   u.Upper90Points,
			GrouchyPoints:   u.GrouchyPoints,
			HasProfile:      true,
		})
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].AcesRadioPoints > entries[j].AcesRadioPoints })

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]any{"entries": entries}); err != nil {
		log.Printf("leaderboard: encode response: %v", err)
	}
}
