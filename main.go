package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"sync"
)

type config struct {
	Addr        string
	Help        string
	ProxyIP     string
	Hosts       string
	TraefikMode bool
	ProxyMode   bool
	NoHosts     bool
	NoWeb       bool
	NoHelp      bool
	ForceCRLF   bool
}

func main() {
	docker, err := newDockerCli()
	if err != nil {
		panic(err)
	}

	var config config
	flag.StringVar(&config.Addr, "addr", ":8080", "Web app address and port")
	flag.StringVar(&config.Help, "help", "https://github.com/lobre/ghosts/blob/master/README.md", "Change the Web help link")
	flag.StringVar(&config.ProxyIP, "proxyip", "127.0.0.1", "Specific proxy IP for hosts entries")
	flag.StringVar(&config.Hosts, "hosts", "", "Custom location for hosts file")
	flag.BoolVar(&config.TraefikMode, "traefikmode", false, "Enable integration with Traefik proxy")
	flag.BoolVar(&config.ProxyMode, "proxymode", false, "Enable proxy")
	flag.BoolVar(&config.NoHosts, "nohosts", false, "Don't generate hosts file")
	flag.BoolVar(&config.NoWeb, "noweb", false, "Don't start web server")
	flag.BoolVar(&config.NoHelp, "nohelp", false, "Disable help on web interface")
	flag.BoolVar(&config.ForceCRLF, "forcecrlf", false, "Force CRLF end of lines")
	flag.Parse()

	if config.Hosts != "" {
		os.Setenv("HOSTS_PATH", config.Hosts)
	}

	var wg sync.WaitGroup

	// Hosts
	if !config.NoHosts {
		wg.Add(1)
		go func() {
			if err := hosts(docker, config); err != nil {
				log.Fatal(err)
			}
			wg.Done()
		}()
	}

	// Web
	if !config.NoWeb {
		wg.Add(1)
		go func() {
			http.Handle("/", &appHandler{config, docker})
			http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
			log.Fatal(http.ListenAndServe(config.Addr, nil))
			wg.Done()
		}()
	}

	wg.Wait()
}
