package main

import (
	"chirpy/internal/auth"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func (cfg *apiConfig) deleteChirpHandler(w http.ResponseWriter, r *http.Request) {
	chirpIDStr := strings.TrimPrefix(r.URL.Path, "/api/chirps/")
	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		respondWithError(w, 400, "uuid format incorrect")
		return
	}
	token, err := auth.GetAuthorizationHeader("Bearer", r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No access token in header")
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.secretKey)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not verify token")
		return
	}
	chirpAuthor, err := cfg.db.GetChirpAuthor(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not get the author of chirp")
		return
	}
	if chirpAuthor.ID != userID {
		respondWithError(w, 403, "Access token doesnt belong to this chrips author")
		return
	}
	err = cfg.db.RemoveChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, 500, "Failed to delete chirp")
		return
	}
	respondWithJSON(w, http.StatusNoContent, nil)
}