package main

import (
	"flag"
	"html/template"
	"net/http"
	"strings"
)

type config struct {
	proxyIP     string
	traefikMode bool
	noProxy     bool
}

var conf config

func main() {
	if err := initDockerCli(); err != nil {
		panic(err)
	}

	flag.StringVar(&conf.proxyIP, "proxy-ip", "127.0.0.1", "Specific proxy IP for hosts entries")
	flag.BoolVar(&conf.traefikMode, "traefik-mode", false, "Enable integration with Traefik proxy")
	flag.BoolVar(&conf.noProxy, "no-proxy", false, "Disable proxy and reach containers by theirs IP")

	//hosts()

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
