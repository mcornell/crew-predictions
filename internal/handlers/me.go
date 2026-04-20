package handlers

import (
	"encoding/json"
	"net/http"
)

func Me(w http.ResponseWriter, r *http.Request) {
	user := UserFromSession(r)
	if user == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"handle":        user.Handle,
		"emailVerified": user.EmailVerified,
	})
}
