package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"sync"
)

type config struct {
	addr        string
	proxyIP     string
	hosts       string
	autoEnabled bool
	traefikMode bool
	directMode  bool
	onlyWeb     bool
}

func main() {
	docker, err := newDockerCli()
	if err != nil {
		panic(err)
	}

	var config config
	flag.StringVar(&config.addr, "addr", ":8080", "Web app address and port")
	flag.StringVar(&config.proxyIP, "proxy-ip", "127.0.0.1", "Specific proxy IP for hosts entries")
	flag.StringVar(&config.hosts, "hosts", "", "Custom location for hosts file")
	flag.BoolVar(&config.autoEnabled, "auto-enabled", true, "Automatically enable new containers without the enabled label")
	flag.BoolVar(&config.traefikMode, "traefik-mode", false, "Enable integration with Traefik proxy")
	flag.BoolVar(&config.directMode, "direct-mode", false, "Disable proxy and reach containers by theirs IP")
	flag.BoolVar(&config.onlyWeb, "only-web", false, "Don't generate hosts file")
	flag.Parse()

	if config.hosts != "" {
		os.Setenv("HOSTS_PATH", config.hosts)
	}

	var wg sync.WaitGroup

	// Hosts
	wg.Add(1)
	go func() {
		if err := hosts(docker, config); err != nil {
			log.Fatal(err)
		}
		wg.Done()
	}()

	// Web
	wg.Add(1)
	go func() {
		http.Handle("/", &appHandler{config, docker})
		http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
		log.Fatal(http.ListenAndServe(config.addr, nil))
		wg.Done()
	}()

	wg.Wait()
}
