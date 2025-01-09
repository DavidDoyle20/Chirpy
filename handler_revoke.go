package main

import (
	"chirpy/internal/auth"
	"log"
	"net/http"
)

func (cfg *apiConfig) revokeHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 500, "No refresh token present in headers")
	}

	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), token)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusUnauthorized, "Refresh token doesnt exist")
		return
	}
	err = cfg.db.RevokeRefreshTokenFromUser(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, 500, "Could not revoke token")
		return
	}
	respondWithJSON(w, http.StatusNoContent, nil)
}