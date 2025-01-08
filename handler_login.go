package main

import (
	"net/http"
	"time"
	"encoding/json"
	"chirpy/internal/auth"
	"log"
)

// check if a user is already logged in
func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}
	type response struct {
		ID         string    `json:"id"`
		Created_at time.Time `json:"created_at"`
		Updated_at time.Time `json:"updated_at"`
		Email      string    `json:"email"`
		Token      string    `json:"token"`
	}

	var params parameters
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "could not decode json body")
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, 400, "a user with that email doesnt exist")
		return
	}

	if err = auth.CheckPasswordHash(params.Password, user.HashedPassword); err != nil {
		log.Println(err, params.Password, user.HashedPassword)
		respondWithError(w, 401, "Incorrect email or password")
		return
	}
	expiresIn := time.Hour
	if params.ExpiresInSeconds != 0 && params.ExpiresInSeconds < 86400{
		expiresIn = time.Duration(params.ExpiresInSeconds) * time.Second
	}
	token, err := auth.MakeJWT(user.ID, cfg.secretKey, expiresIn)
	if err != nil {
		respondWithError(w, 500, "Could not generate jwt")
		return
	}
	cfg.jwt = token

	userAndToken := response{
		ID:         user.ID.String(),
		Created_at: user.CreatedAt,
		Updated_at: user.UpdatedAt,
		Email:      user.Email,
		Token:      token,
	}

	respondWithJSON(w, 200, userAndToken)
}