package main

import (
	"html/template"
	"net/http"
	"strings"
)

func main() {
	if err := initDockerCli(); err != nil {
		panic(err)
	}

	hosts()

	http.HandleFunc("/", index)
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("index.html").Funcs(template.FuncMap{
		"capitalize": capitalize,
	}).ParseFiles("index.html")

	entries, err := entries()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = tmpl.Execute(w, entries)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func capitalize(s string) string {
	return strings.Title(s)
}
