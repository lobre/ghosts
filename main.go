package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type config struct {
	Addr                string
	Help                string
	ProxyIP             string
	Hosts               string
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
	if config.ProxyMode && config.ProxyNetAutoConnect && config.ProxyContainerName != "" {
		np, err := newNetworksProcessor(docker, config, em, config.ProxyContainerName)
		if err != nil {
			log.Fatal(err)
		}

		listener.addProcessor(np)
	}

	// Hosts
	if !config.NoHosts {
		hp := newHostsProcessor(config, em)
		listener.addProcessor(hp)
	}

	var wg sync.WaitGroup
	sigstop := make(chan os.Signal)
	listenerStop := make(chan int)
	signal.Notify(sigstop, syscall.SIGINT, syscall.SIGTERM)

	// Docker listener routine
	wg.Add(1)
	go func() {
		if err := listener.start(listenerStop); err != nil {
			log.Fatal(err)
		}
		log.Print("Listener stopped")
		wg.Done()
	}()

	// Web routine
	server := &http.Server{Addr: config.Addr}
	http.Handle("/", &appHandler{config, em})
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	if !config.NoWeb {
		wg.Add(1)
		go func() {
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatal(err)
			}
			log.Print("Web server stopped")
			wg.Done()
		}()
	}

	// Sigstop signal received
	<-sigstop

	// Gracefully stop web server
	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatal(err)
	}

	// Gracefully stop listener
	listenerStop <- 1

	wg.Wait()
}
