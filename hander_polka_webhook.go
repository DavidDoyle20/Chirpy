package main

import (
	"chirpy/internal/auth"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) polkaWebhookHander(w http.ResponseWriter, r *http.Request) {
	type EventData struct {
		UserID string `json:"user_id"`
	}
	type body struct {
		Event string `json:"event"`
		Data EventData `json:"data"`
	}

	var b body
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&b)
	if err != nil {
		log.Println("Could not decode the json body")
		respondWithJSON(w, http.StatusNoContent, nil)
		return
	}

	if b.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusNoContent, nil)
		return
	}

	apiKey, err := auth.GetAuthorizationHeader("ApiKey", r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Api key not present")
		return
	}
	if apiKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "Invalid api key")
		return
	}

	userID, err := uuid.Parse(b.Data.UserID)
	if err != nil {
		log.Panicln("Could not parse user id")
		respondWithJSON(w, http.StatusNoContent, nil)
	}

	err = cfg.db.UpgradeUserToChirpyRed(r.Context(), userID)
	if err != nil {
		log.Panicln("Could not upgrade user to chirpy red")
		respondWithError(w, http.StatusNotFound, "User doesnt exist")
	}
	respondWithJSON(w, http.StatusNoContent, nil)
}