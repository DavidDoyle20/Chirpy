package main

import (
	"net/http"
	"log"
	"strconv"
	"os"
	"strings"
)
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
		next.ServeHTTP(w, r)
	})
}