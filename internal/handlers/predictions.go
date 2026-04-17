package handlers

import "net/http"

func SubmitPrediction(w http.ResponseWriter, r *http.Request) {
	if userFromSession(r) == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
}
