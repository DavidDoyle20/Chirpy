package main

import (
	"net/http"
)

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Swap(0)
	err := cfg.db.ResetUsers(r.Context())
	if err != nil {
		respondWithError(w, 500, "could not reset the users")
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Reset Metrics"))
}