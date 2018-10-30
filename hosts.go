package main

import (
	"bufio"
	"io/ioutil"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
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
			id := msg.ID
			action := msg.Action
			spew.Dump(msg)
			if err := generateHosts(docker, config); err != nil {
				return err
			}
		case err := <-errCh:
			return err
		}
	}
}

func generateHosts(docker docker, config config) error {
	if config.onlyWeb {
		return nil
	}

	entries, err := getEntries(docker, config)
	if err != nil {
		return err
	}
	spew.Dump(entries)

	if err := cleanEntries(entries); err != nil {
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

func cleanEntries(entries []entry) error {
	hosts, err := goodhosts.NewHosts()
	if err != nil {
		return err
	}

	f, err := os.Open(hosts.Path)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(f)
	lines := []string{}

	for scanner.Scan() {
		line := scanner.Text()
		keep := true
		for _, entry := range entries {
			if strings.Contains(line, entry.Host) {
				keep = false
				break
			}
		}
		if keep {
			lines = append(lines, line)
		}
	}

	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(hosts.Path, []byte(output), 0644)
	if err != nil {
		return err
	}

	return nil
}
