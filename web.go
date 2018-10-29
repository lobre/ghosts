package main

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

type frontEntries map[string][]entry

type appHandler struct {
	config config
	docker docker
}

func (h *appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("index.html").Funcs(template.FuncMap{
		"capitalize": strings.Title,
		"upper":      strings.ToUpper,
	}).ParseFiles("index.html")

	entries, err := getEntries(h.docker, h.config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	spew.Dump(entries)

	err = tmpl.Execute(w, struct {
		Config  config
		Entries frontEntries
	}{
		h.config,
		prepare(entries),
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Separate entries by categories
func prepare(entries []entry) frontEntries {
	categories := make(frontEntries)
	for _, entry := range entries {
		categories[entry.Category] = append(categories[entry.Category], entry)
	}
	return categories
}
