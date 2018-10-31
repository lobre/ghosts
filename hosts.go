package main

import (
	"github.com/lextoumbourou/goodhosts"
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

		if (entry.Direct || (!config.ProxyMode && !config.TraefikMode)) && !hosts.Has(entry.IP, entry.Host) {
			hosts.Add(entry.IP, entry.Host)
		} else if !hosts.Has(config.ProxyIP, entry.Host) {
			hosts.Add(config.ProxyIP, entry.Host)
		}
	}

	if err := hosts.Flush(); err != nil {
		return err
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
		if (entry.Direct || (!config.ProxyMode && !config.TraefikMode)) && hosts.Has(entry.IP, entry.Host) {
			hosts.Remove(entry.IP, entry.Host)
		} else if hosts.Has(config.ProxyIP, entry.Host) {
			hosts.Remove(config.ProxyIP, entry.Host)
		}
	}

	if err := hosts.Flush(); err != nil {
		return err
	}

	return nil
}
