package main

import (
	"fmt"
	"html/template"
	"path/filepath"
	"runtime"
	"time"

	"github.com/andruixxd31/snippet-box/internal/models"
)

type templateData struct {
	CurrentYear     int
	Snippet         models.Snippet
	Snippets        []models.Snippet
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.UTC().Format("02 Jan 2006 at 15:04")
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Join(filepath.Dir(b), "..", "..")
	templatesPath := filepath.Join(basepath, "ui", "html", "pages", "*.tmpl")

	fmt.Println(templatesPath)
	pages, err := filepath.Glob(templatesPath)
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		baseFilePath := filepath.Join(basepath, "ui", "html", "base.tmpl")
		ts, err := template.New(name).Funcs(functions).ParseFiles(baseFilePath)
		if err != nil {
			return nil, err
		}

		partialsPath := filepath.Join(basepath, "ui", "html", "partials", "*.tmpl")
		ts, err = ts.ParseGlob(partialsPath)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
