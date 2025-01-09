package main

import (
	"chirpy/internal/database"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/google/uuid"
)
func (cfg *apiConfig) getChirpsHandler(w http.ResponseWriter, r *http.Request) {
	var chirps []database.Chirp
	var err error

	sortDirection := strings.ToUpper(r.URL.Query().Get("sort"))
	log.Println(sortDirection)


	authorID := r.URL.Query().Get("author_id")
	if authorID != "" {
		authorID, err := uuid.Parse(authorID)
		if err != nil {
			log.Println(err)
			respondWithError(w, 400, "could not get chirps")
			return
		}
		chirps, err = cfg.db.GetChirpsByAuthor(r.Context(), authorID)
		if err != nil {
			log.Println(err)
			respondWithError(w, 400, "could not get chirps")
			return
		}
	} else {
		chirps, err = cfg.db.GetChirps(r.Context())
		if err != nil {
			log.Println(err)
			respondWithError(w, 400, "could not get chirps")
			return
		}
	}

	if sortDirection == "DESC" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.After((chirps[j].CreatedAt))
		})
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
	chirpIDStr := strings.TrimPrefix(r.URL.Path, "/api/chirps/")
	log.Println(chirpIDStr)
	chirpID, err := uuid.Parse(chirpIDStr)
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