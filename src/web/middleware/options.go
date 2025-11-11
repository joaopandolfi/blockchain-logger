package middleware

import (
	"net/http"

	"github.com/joaopandolfi/blackwhale/handlers"
)

// Options allow to
func Options(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	if r.Method == http.MethodOptions {
		return
	}

	handlers.Response(w, "", 200)
}
