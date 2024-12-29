package main

import (
	"log"
	"net/http"
	"path/filepath"
)

func main() {
	mux := http.NewServeMux()

	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static")})
	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /snippet/view/{id}", snippetView)
	mux.HandleFunc("GET /snippet/create", snippetCreate)
	mux.HandleFunc("POST /snippet/createPost", snippetCreatePost)

	log.Print("Starting server on port: 4000")

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}

type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	log.Printf("Trying to open path: %s", path)
	f, err := nfs.fs.Open(path)
	if err != nil {
		log.Printf("Error opening path: %s, error: %v", path, err)
		return nil, err
	}

	s, err := f.Stat()
	if err != nil {
		log.Printf("Error getting stats for path: %s, error: %v", path, err)
		return nil, err
	}

	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		log.Printf("Path is a directory. Checking for index.html at: %s", index)
		if _, err := nfs.fs.Open(index); err != nil {
			log.Printf("index.html not found in directory: %s", path)
			closeErr := f.Close()
			if closeErr != nil {
				log.Printf("Error closing file: %v", closeErr)
				return nil, closeErr
			}
			return nil, err
		}
	}

	return f, nil
}
