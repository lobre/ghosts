package main

import (
	"flag"
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

type config struct {
	proxyIP     string
	hosts       string
	traefikMode bool
	directMode  bool
	onlyWeb     bool
}

func main() {
	cli, err := newDockerCli()
	if err != nil {
		panic(err)
	}

	var conf config
	flag.StringVar(&conf.proxyIP, "proxy-ip", "127.0.0.1", "Specific proxy IP for hosts entries")
	flag.StringVar(&conf.hosts, "hosts", "", "Custom location for hosts file")
	flag.BoolVar(&conf.traefikMode, "traefik-mode", false, "Enable integration with Traefik proxy")
	flag.BoolVar(&conf.directMode, "direct-mode", false, "Disable proxy and reach containers by theirs IP")
	flag.BoolVar(&conf.onlyWeb, "only-web", false, "Don't generate hosts file")
	flag.Parse()

	if conf.hosts != "" {
		os.Setenv("HOSTS_PATH", conf.hosts)
	}

	entries, err := getEntries(cli, conf)
	if err != nil {
		panic(err)
	}
	spew.Dump(entries)

	//hosts()

	// http.HandleFunc("/", index)
	// http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("index.html").Funcs(template.FuncMap{
		"capitalize": capitalize,
	}).ParseFiles("index.html")

	// entries, err := entries()
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// }

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func capitalize(s string) string {
	return strings.Title(s)
}
