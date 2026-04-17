package handlers

import "net/http"

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
}
