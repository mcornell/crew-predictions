package handlers

import (
	"fmt"
	"net/http"
)

func Matches(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "<h1>Upcoming Matches</h1>")
}
