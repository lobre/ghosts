package main

import (
	"fmt"
	"strings"
)

const prefix = "ghosts"
const defaultCategory = "others"
const noName = "unknown"

type entry struct {
	Host  string
	IP    string
	Port  string
	Proto string // TODO

	Name        string
	Category    string
	Description string
	Logo        string // TODO
	Auth        string // TODO

	Hide   bool
	Direct bool
}

func getEntries(cli cli, conf config) (map[string][]entry, error) {
	entries := make(map[string][]entry)

	containers, err := cli.getContainers()
	if err != nil {
		return entries, err
	}

	for _, container := range containers {
		entry := entry{}

		// Take the IP of the first network
		for _, n := range container.NetworkSettings.Networks {
			entry.IP = n.IPAddress
			break
		}

		// Take the first port exposed
		for _, p := range container.Ports {
			entry.Port = fmt.Sprint(p.PrivatePort)
			break
		}

		// Host
		if val, ok := container.Labels[fmt.Sprintf("%s.host", prefix)]; ok {
			entry.Host = val
		} else if val, ok := container.Labels["traefik.frontend.rule"]; ok {
			entry.Host = val
		}

		// Protocol
		if val, ok := container.Labels[fmt.Sprintf("%s.protocol", prefix)]; ok {
			entry.Proto = val
		}

		// Name
		if val, ok := container.Labels[fmt.Sprintf("%s.name", prefix)]; ok {
			entry.Name = val
		} else if len(container.Names) > 0 {
			entry.Name = strings.TrimPrefix(container.Names[0], "/")
		} else {
			entry.Name = noName
		}

		// Category
		category := defaultCategory
		if val, ok := container.Labels[fmt.Sprintf("%s.category", prefix)]; ok {
			category = val
		}

		// Description
		if val, ok := container.Labels[fmt.Sprintf("%s.description", prefix)]; ok {
			entry.Description = val
		}

		// Hide
		if val, ok := container.Labels[fmt.Sprintf("%s.hide", prefix)]; ok && val == "true" {
			entry.Hide = true
		}

		// Direct
		if val, ok := container.Labels[fmt.Sprintf("%s.direct", prefix)]; ok && val == "true" {
			entry.Direct = true
		}

		entries[strings.ToLower(category)] = append(entries[category], entry)
	}
	return entries, nil
}
