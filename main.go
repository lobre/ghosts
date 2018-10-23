package main

import (
	"flag"
	"net/http"
	"os"
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
	cli, err := newDockerCli()
	if err != nil {
		panic(err)
	}

	var conf config
	flag.StringVar(&conf.addr, "addr", ":8080", "Web app address and port")
	flag.StringVar(&conf.proxyIP, "proxy-ip", "127.0.0.1", "Specific proxy IP for hosts entries")
	flag.StringVar(&conf.hosts, "hosts", "", "Custom location for hosts file")
	flag.BoolVar(&conf.autoEnabled, "auto-enabled", true, "Automatically enable new containers without the enabled label")
	flag.BoolVar(&conf.traefikMode, "traefik-mode", false, "Enable integration with Traefik proxy")
	flag.BoolVar(&conf.directMode, "direct-mode", false, "Disable proxy and reach containers by theirs IP")
	flag.BoolVar(&conf.onlyWeb, "only-web", false, "Don't generate hosts file")
	flag.Parse()

	if conf.hosts != "" {
		os.Setenv("HOSTS_PATH", conf.hosts)
	}

	// Start web server
	http.Handle("/", &appHandler{conf, cli})
	http.ListenAndServe(conf.addr, nil)
}
