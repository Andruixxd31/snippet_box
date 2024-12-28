package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /snippet/view/{id}", snippetView)
	mux.HandleFunc("GET /snippet/create", snippetCreate)
	mux.HandleFunc("POST /snippet/createPost", snippetCreatePost)

	log.Print("Starting server on port: 4000")

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
