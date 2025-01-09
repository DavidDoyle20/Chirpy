package main

import (
	"net/http"
	"encoding/json"
	"log"
	"chirpy/internal/auth"
	"strings"
	"chirpy/internal/database"
)

func (cfg *apiConfig) createChirpHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	var params parameters

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		log.Println("Could not decode json body", err)
		respondWithError(w, 500, err.Error())
		return
	}
	tokenString, err := auth.GetAuthorizationHeader("Bearer", r.Header)
	if err != nil {
		log.Println("There was an error while getting the bearer token")
		respondWithError(w, 401, "Unauthorized")
		return
	}
	userID, err := auth.ValidateJWT(tokenString, cfg.secretKey)
	if err != nil {
		log.Println("Could not validate jwt", err)
		respondWithError(w, 401, "Unauthorized")
		return
	}

	if !validateLength(params.Body) {
		log.Println("Body is longer than 140 characters")
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	chrp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{UserID: userID, Body: params.Body})
	if err != nil {
		log.Println(err)
		respondWithJSON(w, 500, "Could not create the chirp")
		return
	}
	respondWithJSON(w, 201, Chirp{
		Id: chrp.ID,
		Created_at: chrp.CreatedAt,
		Updated_at: chrp.UpdatedAt,
		Body: chrp.Body,
		User_id: chrp.UserID,
	})
}

func validateLength(s string) bool {
	return len(s) <= 140
}

func replaceBadWords(s string) string {
	words := strings.Split(s, " ")
	for i, w := range words {
		w = strings.ToLower(w)
		if _, in := bannedWords[w]; in {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}