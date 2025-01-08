package main

import (
	"net/http"
	"encoding/json"
	"log"
	"chirpy/internal/database"
	"chirpy/internal/auth"
)

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var params parameters
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "could not decode json body")
		return
	}

	log.Println(params.Email)
	hash, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 500, "could not hash passowrd")
		return
	}
	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{Email: params.Email, HashedPassword: hash})
	if err != nil {
		log.Println(err)
		respondWithError(w, 500, "could not create user")
		return
	}

	respondWithJSON(w, 201, User{
		Id: user.ID,
		Created_at: user.CreatedAt,
		Updated_at: user.UpdatedAt,
		Email: user.Email,
	})
}