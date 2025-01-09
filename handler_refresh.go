package main

import (
	"chirpy/internal/auth"
	"net/http"
)

func (cfg *apiConfig) refreshHandler(w http.ResponseWriter, r *http.Request) {
	type resp struct {
		Token string `json:"token"`
	}
	// requires refresh token to be present in the header
	token, err := auth.GetAuthorizationHeader("Bearer", r.Header)
	if err != nil {
		respondWithError(w, 500, "No refresh token present in headers")
	}
	// token was present
	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Refresh token doesnt exist")
		return
	}
	_, err = cfg.db.CheckAndFetchRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token is expired")
		return	
	}
	access_token, err := auth.MakeJWT(user.ID, cfg.secretKey)
	if err != nil {
		respondWithError(w, 500, "Unable to make access token")
		return
	}
	respondWithJSON(w, http.StatusOK, resp{
		Token: access_token,
	})
}