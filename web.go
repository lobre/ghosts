package main

import (
	"html/template"
	"net/http"
	"strings"
)

type frontEntries map[string][]entry

type appHandler struct {
	config config
	em     entriesManager
}

func (h *appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("index.html").Funcs(template.FuncMap{
		"capitalize": strings.Title,
		"upper":      strings.ToUpper,
	}).ParseFiles("index.html")

	entries, err := h.em.get()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = tmpl.Execute(w, struct {
		Config         config
		EntriesManager entriesManager
		Entries        frontEntries
	}{
		h.config,
		h.em,
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
		if !entry.NoWeb {
			categories[entry.Category] = append(categories[entry.Category], entry)
		}
	}
	return categories
}
