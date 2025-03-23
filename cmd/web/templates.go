package main

import "github.com/andruixxd31/snippet-box/internal/models"

type templateData struct {
	Snippet  models.Snippet
	Snippets []models.Snippet
}
