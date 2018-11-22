package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"sync"
)

type config struct {
	Addr                string
	Help                string
	ProxyIP             string
	Hosts               string
	TraefikMode         bool
	ProxyMode           bool
	ProxyNetAutoConnect bool
	ProxyContainerName  string
	NoHosts             bool
	NoWeb               bool
	NoHelp              bool
	HostsForceCRLF      bool
}

type processor interface {
	init() error
	startEvent(id string) error
	dieEvent(id string) error
}

func main() {
	docker, err := newDocker()
	if err != nil {
		log.Fatal(err)
	}

	var config config
	flag.StringVar(&config.Addr, "addr", ":8080", "Web app address and port")
	flag.StringVar(&config.Help, "help", "https://github.com/lobre/ghosts/blob/master/README.md", "Change the Web help link")
	flag.StringVar(&config.ProxyIP, "proxyip", "127.0.0.1", "Specific proxy IP for hosts entries")
	flag.StringVar(&config.Hosts, "hosts", "", "Custom location for hosts file")
	flag.BoolVar(&config.TraefikMode, "traefikmode", false, "Enable integration with Traefik proxy")
	flag.BoolVar(&config.ProxyMode, "proxymode", false, "Enable proxy")
	flag.BoolVar(&config.ProxyNetAutoConnect, "proxynetautoconnect", false, "Enable automatic network connection between proxy and containers")
	flag.StringVar(&config.ProxyContainerName, "proxycontainername", "", "Name of proxy container")
	flag.BoolVar(&config.NoHosts, "nohosts", false, "Don't generate hosts file")
	flag.BoolVar(&config.NoWeb, "noweb", false, "Don't start web server")
	flag.BoolVar(&config.NoHelp, "nohelp", false, "Disable help on web interface")
	flag.BoolVar(&config.HostsForceCRLF, "hostsforcecrlf", false, "Force CRLF end of lines when generating hosts entries")
	flag.Parse()

	if config.Hosts != "" {
		os.Setenv("HOSTS_PATH", config.Hosts)
	}

	listener := newListener(docker)
	em := newEntriesManager(docker, config)

	// Network
	if (config.ProxyMode || config.TraefikMode) && config.ProxyNetAutoConnect {
		proxyName := config.ProxyContainerName
		if proxyName == "" && config.TraefikMode {
			proxyName = "traefik"
		}
		if proxyName != "" {
			np, err := newNetworksProcessor(docker, config, em, proxyName)
			if err != nil {
				log.Fatal(err)
			}

			listener.addProcessor(np)
		}
	}

	// Hosts
	if !config.NoHosts {
		hp := newHostsProcessor(config, em)
		listener.addProcessor(hp)
	}

	var wg sync.WaitGroup
	var stop <-chan int

	// Docker listener routine
	wg.Add(1)
	go func() {
		if err := listener.start(stop); err != nil {
			log.Fatal(err)
		}
		wg.Done()
	}()

	// Web routine
	if !config.NoWeb {
		wg.Add(1)
		go func() {
			http.Handle("/", &appHandler{config, em})
			http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
			log.Fatal(http.ListenAndServe(config.Addr, nil))
			wg.Done()
		}()
	}

	wg.Wait()
}
