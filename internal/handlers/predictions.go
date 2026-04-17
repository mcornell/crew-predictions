package handlers

import (
	"net/http"
	"strconv"
)

func SubmitPrediction(w http.ResponseWriter, r *http.Request) {
	if userFromSession(r) == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	r.ParseForm()
	if r.FormValue("match_id") == "" || r.FormValue("home_goals") == "" || r.FormValue("away_goals") == "" {
		http.Error(w, "missing fields", http.StatusBadRequest)
		return
	}
	if _, err := strconv.Atoi(r.FormValue("home_goals")); err != nil {
		http.Error(w, "goals must be integers", http.StatusBadRequest)
		return
	}
	if _, err := strconv.Atoi(r.FormValue("away_goals")); err != nil {
		http.Error(w, "goals must be integers", http.StatusBadRequest)
		return
	}
}
