package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
)
func (cfg *apiConfig) getChirpsHandler(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		log.Println(err)
		respondWithError(w, 400, "could not get chirps")
		return
	}

	var chirpResp []Chirp
	for _, c := range chirps {
		chirpResp = append(chirpResp, Chirp{
			Id: c.ID,
			Created_at: c.CreatedAt,
			Updated_at: c.UpdatedAt,
			Body: c.Body,
			User_id: c.UserID,
		})
	}
	respondWithJSON(w, 200, chirpResp)
}

func (cfg *apiConfig) getChirpByIDHandler(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		log.Println(err)
		respondWithError(w, 400, "uuid format incorrect")
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	respondWithJSON(w, 200, Chirp{
		Id: chirp.ID,
		Created_at: chirp.CreatedAt,
		Updated_at: chirp.UpdatedAt,
		Body: chirp.Body,
		User_id: chirp.UserID,
	})
}