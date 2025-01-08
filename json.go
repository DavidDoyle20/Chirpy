package main

import (
	"encoding/json"
	"log"
	"net/http"
)
func respondWithError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	type Error struct {
		Error string `json:"error"`
	}
	e := Error{
		Error: msg,
	}
	data, err := json.Marshal(e)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(code)
	w.Write([]byte(data))
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(payload)
	if err != nil {
		respondWithError(w, 500, "Server failed to marshal JSON")
		return
	}

	w.WriteHeader(code)
	w.Write(data)
}