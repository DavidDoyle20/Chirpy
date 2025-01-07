package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
)

var bannedWords = map[string]struct{}{
	"kerfuffle": {},
	"sharbert": {},
	"fornax": {},
}


type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	// multiplexer: connects one input line to an output line
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}
	apiCfg := apiConfig {
		fileserverHits: atomic.Int32{},
	}

	// connects a pattern in the url to a file path
	handler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))


	mux.HandleFunc("GET /admin/metrics", apiCfg.AdminMetricsHandler)
	mux.HandleFunc("GET /api/healthz", ReadinessServeHTTP)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)
	mux.HandleFunc("POST /api/validate_chirp", validateChirpHandler)

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}

}

func ReadinessServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) AdminMetricsHandler(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		htmlString, err := os.ReadFile("admin.html")
		if err != nil {
			log.Printf("Error reading from file")
			w.WriteHeader(500)
		}
		w.Write([]byte(strings.Replace(string(htmlString), "%d", strconv.Itoa(int(cfg.fileserverHits.Load())), 1)))
		
}

// Writes the number of requests that have been counted
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w,r)
	})
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Swap(0)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Reset Metrics"))
}

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type cleanedBody struct {
		Cleaned_Body string `json:"cleaned_body"`
	}

	var params parameters

	// Unmarshal the json body
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		log.Println(err)
		respondWithError(w, 500, err.Error())
		return
	}

	if !validateLength(params.Body) {
		log.Println("Body is longer than 140 characters")
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	respondWithJSON(w, 200, cleanedBody{Cleaned_Body: replaceBadWords(params.Body)})
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

// customizable
// have to marshal the json first
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
