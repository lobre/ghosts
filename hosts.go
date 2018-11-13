package main

import (
	"bytes"
	"io/ioutil"

	"github.com/lobre/goodhosts"
)

func hosts(docker docker, config config) error {
	// Initialize hosts
	if err := addHosts(docker, config); err != nil {
		return err
	}

	// Listen to Docker events
	msgCh, errCh := docker.listenContainers()
	for {
		select {
		case msg := <-msgCh:
			switch msg.Action {
			case "start":
				if err := addHosts(docker, config, msg.ID); err != nil {
					return err
				}
			case "die":
				if err := removeHosts(docker, config, msg.ID); err != nil {
					return err
				}
			}
		case err := <-errCh:
			return err
		}
	}
}

func addHosts(docker docker, config config, ids ...string) error {
	entries, err := getEntries(docker, config, ids...)
	if err != nil {
		return err
	}

	hosts, err := goodhosts.NewHosts()
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.NoHosts || entry.WebDirect {
			continue
		}

		ip := config.ProxyIP
		if entry.Direct || (!config.ProxyMode && !config.TraefikMode) {
			ip = entry.IP
		}

		for _, host := range entry.Hosts {
			if !hosts.Has(ip, host) {
				hosts.Add(ip, host)
			}
		}
	}

	if err := hosts.Flush(); err != nil {
		return err
	}

	if config.HostsForceCRLF {
		if err := forceCrlfEOL(); err != nil {
			return err
		}
	}

	return nil
}

func removeHosts(docker docker, config config, ids ...string) error {
	entries, err := getEntries(docker, config, ids...)
	if err != nil {
		return err
	}

	hosts, err := goodhosts.NewHosts()
	if err != nil {
		return err
	}

	for _, entry := range entries {
		ip := config.ProxyIP
		if entry.Direct || (!config.ProxyMode && !config.TraefikMode) {
			ip = entry.IP
		}

		for _, host := range entry.Hosts {
			if hosts.Has(ip, host) {
				hosts.Remove(ip, host)
			}
		}
	}

	if err := hosts.Flush(); err != nil {
		return err
	}

	if config.HostsForceCRLF {
		if err := forceCrlfEOL(); err != nil {
			return err
		}
	}

	return nil
}

func forceCrlfEOL() error {
	hosts, err := goodhosts.NewHosts()
	if err != nil {
		return err
	}

	read, err := ioutil.ReadFile(hosts.Path)
	if err != nil {
		return err
	}

	// Replace lf by crlf
	normalized := bytes.Replace(read, []byte{10}, []byte{13, 10}, -1)

	err = ioutil.WriteFile(hosts.Path, normalized, 0)
	if err != nil {
		return err
	}

	return nil
}
