package main

import (
	"github.com/lextoumbourou/goodhosts"
)

func hosts(docker docker, config config) error {
	// Initialize hosts
	if err := generateHosts(docker, config); err != nil {
		return err
	}

	// Listen to Docker events
	msgCh, errCh := docker.listenContainers()
	for {
		select {
		case msg := <-msgCh:
			switch msg.Action {
			case "die":
				if err := cleanHosts(docker, config, msg.ID); err != nil {
					return err
				}
			case "start":
				if err := generateHosts(docker, config, msg.ID); err != nil {
					return err
				}
			}
		case err := <-errCh:
			return err
		}
	}
}

func generateHosts(docker docker, config config, ids ...string) error {
	if config.onlyWeb {
		return nil
	}

	entries, err := getEntries(docker, config, ids...)
	if err != nil {
		return err
	}

	hosts, err := goodhosts.NewHosts()
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.OnlyWeb {
			continue
		}

		if (config.directMode || entry.Direct) && !hosts.Has(entry.IP, entry.Host) {
			hosts.Add(entry.IP, entry.Host)
		} else if !hosts.Has(config.proxyIP, entry.Host) {
			hosts.Add(config.proxyIP, entry.Host)
		}
	}

	if err := hosts.Flush(); err != nil {
		return err
	}

	return nil
}

func cleanHosts(docker docker, config config, ids ...string) error {
	if config.onlyWeb {
		return nil
	}

	entries, err := getEntries(docker, config, ids...)
	if err != nil {
		return err
	}

	hosts, err := goodhosts.NewHosts()
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.OnlyWeb {
			continue
		}

		if (config.directMode || entry.Direct) && hosts.Has(entry.IP, entry.Host) {
			hosts.Remove(entry.IP, entry.Host)
		} else if hosts.Has(config.proxyIP, entry.Host) {
			hosts.Remove(config.proxyIP, entry.Host)
		}
	}

	if err := hosts.Flush(); err != nil {
		return err
	}

	return nil
}
