package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mcornell/crew-predictions/internal/seasons"
)

type SeasonsHandler struct{}

func NewSeasonsHandler() *SeasonsHandler {
	return &SeasonsHandler{}
}

type seasonResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	IsCurrent bool   `json:"isCurrent"`
}

func (h *SeasonsHandler) APIList(w http.ResponseWriter, r *http.Request) {
	all := seasons.AllSeasons()
	out := make([]seasonResponse, len(all))
	for i, s := range all {
		out[i] = seasonResponse{ID: s.ID, Name: s.Name}
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]any{"seasons": out}); err != nil {
		log.Printf("seasons: encode response: %v", err)
	}
}
