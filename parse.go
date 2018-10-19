package main

import (
	"fmt"
	"strings"
)

const prefix = "ghosts"
const defaultCategory = "others"

type host struct {
	IP   string
	Port string
}

type entry struct {
	Host  string
	Proto string
	Path  string // is it usefull?

	Name        string
	Category    string
	Description string
	Logo        string
	Auth        string

	Hide string
}

func parse() (map[string][]entry, []host, error) {
	entries := make(map[string][]entry)
	hosts := []host{}

	containers, err := getContainers()
	if err != nil {
		return entries, hosts, err
	}

	for _, container := range containers {
		category := defaultCategory
		entry := entry{}
		host := host{}

		// Take the IP of the first network
		for _, n := range container.NetworkSettings.Networks {
			host.IP = n.IPAddress
			break
		}

		// Take the first port exposed
		for _, p := range container.Ports {
			host.Port = fmt.Sprint(p.PrivatePort)
			break
		}

		// Host
		if val, ok := container.Labels[fmt.Sprintf("%s.host", prefix)]; ok {
			entry.Host = val
		} else if val, ok := container.Labels["traefik.frontend.rule"]; ok {
			entry.Host = val
		}

		// Category
		if val, ok := container.Labels[fmt.Sprintf("%s.category", prefix)]; ok {
			category = val
		}

		// Name
		if len(container.Names) > 0 {
			entry.Name = strings.TrimPrefix(container.Names[0], "/")
		}

		// Protocol
		if val, ok := container.Labels[fmt.Sprintf("%s.protocol", prefix)]; ok {
			entry.Proto = val
		}

		// Host
		if val, ok := container.Labels[fmt.Sprintf("%s.host", prefix)]; ok {
			entry.Host = val
		}

		entries[strings.ToLower(category)] = append(entries[category], entry)
	}
	return entries, nil
}
