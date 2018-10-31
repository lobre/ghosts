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
	traefikMode bool
	proxyMode   bool
	noHosts     bool
	noWeb       bool
}

func main() {
	docker, err := newDockerCli()
	if err != nil {
		panic(err)
	}

	var config config
	flag.StringVar(&config.addr, "addr", ":8080", "Web app address and port")
	flag.StringVar(&config.proxyIP, "proxyip", "127.0.0.1", "Specific proxy IP for hosts entries")
	flag.StringVar(&config.hosts, "hosts", "", "Custom location for hosts file")
	flag.BoolVar(&config.traefikMode, "traefikmode", false, "Enable integration with Traefik proxy")
	flag.BoolVar(&config.proxyMode, "proxymode", false, "Enable proxy")
	flag.BoolVar(&config.noHosts, "nohosts", false, "Don't generate hosts file")
	flag.BoolVar(&config.noWeb, "noweb", false, "Don't start web server")
	flag.Parse()

	if config.hosts != "" {
		os.Setenv("HOSTS_PATH", config.hosts)
	}

	var wg sync.WaitGroup

	// Hosts
	if !config.noHosts {
		wg.Add(1)
		go func() {
			if err := hosts(docker, config); err != nil {
				log.Fatal(err)
			}
			wg.Done()
		}()
	}

	// Web
	if !config.noWeb {
		wg.Add(1)
		go func() {
			http.Handle("/", &appHandler{config, docker})
			http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
			log.Fatal(http.ListenAndServe(config.addr, nil))
			wg.Done()
		}()
	}

	wg.Wait()
}
