package main

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// check if a user is already logged in
func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
	}
	type response struct {
		ID         string    `json:"id"`
		Created_at time.Time `json:"created_at"`
		Updated_at time.Time `json:"updated_at"`
		Email      string    `json:"email"`
		Token      string    `json:"token"`
		Is_Chirpy_Red bool `json:"is_chirpy_red"`
		Refresh_token string `json:"refresh_token"`
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
	token, err := auth.MakeJWT(user.ID, cfg.secretKey)
	if err != nil {
		respondWithError(w, 500, "Could not generate jwt")
		return
	}

	refresh_token, err := auth.MakeRefreshToken()
	cfg.db.RevokeRefreshTokenFromUser(r.Context(), user.ID)
	cfg.db.AssignRefreshTokenToUser(r.Context(), database.AssignRefreshTokenToUserParams{
		Token: refresh_token,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, 500, "Could not generate refresh token")
		return
	}

	userAndToken := response{
		ID:         user.ID.String(),
		Created_at: user.CreatedAt,
		Updated_at: user.UpdatedAt,
		Email:      user.Email,
		Is_Chirpy_Red: user.IsChirpyRed.Bool,
		Token:      token,
		Refresh_token: refresh_token,
	}

	respondWithJSON(w, 200, userAndToken)
}