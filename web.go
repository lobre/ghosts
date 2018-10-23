package main

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

type frontEntries map[string][][]entry

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
	spew.Dump(entries)

	err = tmpl.Execute(w, struct {
		DefaultCategory string
		Entries         frontEntries
	}{
		defaultCategory,
		paginate(entries),
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func capitalize(s string) string {
	return strings.Title(s)
}

// TODO
func paginate([]entry) frontEntries {
	return nil
}
