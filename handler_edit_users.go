package main

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) editUsersHandler(w http.ResponseWriter, r *http.Request) {
	type body struct {
		NewPassword string `json:"password"`
		NewEmail string `json:"email"`
	}
	type resp struct {
		ID string `json:"id"`
		Created_at string `json:"created_at"`
		Updated_at string `json:"updated_at"`
		Email string `json:"email"`
	}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Could not get access token")
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secretKey)
	if err != nil {
		respondWithError(w, 401, "Could not validate jwt")
		return
	}

	var b body
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&b)
	if err != nil {
		respondWithError(w, 401, "Could not get body")
		return
	}

	hashed, err := auth.HashPassword(b.NewPassword)
	if err != nil {
		respondWithError(w, 401, "Could not hash password")
		return
	}

	err = cfg.db.ChangeUserPassword(r.Context(), database.ChangeUserPasswordParams{
		ID: userID,
		HashedPassword: hashed,
	})
	if err != nil {
		respondWithError(w, 401, "Could not change password")
		return
	}

	err = cfg.db.ChangeUserEmail(r.Context(), database.ChangeUserEmailParams{
		ID: userID,
		Email: b.NewEmail,
	})
	if err != nil {
		respondWithError(w, 401, "Could not change email")
		return
	}

	user, err := cfg.db.GetUser(r.Context(), userID)
	if err != nil {
		respondWithError(w, 401, "Credentials were change but unable to return user")
		return
	}

	respondWithJSON(w, 200, resp{
		ID: userID.String(),
		Created_at: user.CreatedAt.String(),
		Updated_at: user.UpdatedAt.String(),
		Email: user.Email,
	})
}