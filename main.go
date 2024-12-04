package main

import (
	"fmt"
	"net/http"
)

func main() {
	// multiplexer: connects one input line to an output line
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}

	// connects a pattern in the url to a file path
	mux.Handle("/", http.FileServer(http.Dir(".")))

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}

}
