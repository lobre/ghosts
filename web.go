package main

import (
	"html/template"
	"net/http"
	"strings"
)

type frontEntries map[string][]entry

type appHandler struct {
	conf config
	cli  cli
}

func (h *appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("index.html").Funcs(template.FuncMap{
		"capitalize": capitalize,
	}).ParseFiles("index.html")

	entries, err := getEntries(h.cli, h.conf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = tmpl.Execute(w, struct {
		Entries frontEntries
	}{
		prepare(entries),
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func capitalize(s string) string {
	return strings.Title(s)
}

// Separate entries by categories
func prepare(entries []entry) frontEntries {
	categories := make(frontEntries)
	for _, entry := range entries {
		categories[entry.Category] = append(categories[entry.Category], entry)
	}
	return categories
}
