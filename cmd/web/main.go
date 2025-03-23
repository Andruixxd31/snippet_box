package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/andruixxd31/snippet-box/internal/models"
	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	logger        *slog.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "andres:xiuxiu@/snippetbox?parseTime=true", "Mysql data source name")
	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	mux := http.NewServeMux()

	app := &application{
		logger:        logger,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static")})
	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
	mux.HandleFunc("GET /snippet/create", app.snippetCreate)
	mux.HandleFunc("POST /snippet/createPost", app.snippetCreatePost)

	logger.Info("starting server", "addr", *addr)

	err = http.ListenAndServe(*addr, mux)
	logger.Error(err.Error())
	os.Exit(1)
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

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
