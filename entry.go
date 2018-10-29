package main

import (
	"fmt"
	"strings"
)

const labelPrefix string = "ghosts"
const defaultCategory string = "apps"

type entry struct {
	Host  string
	IP    string
	Port  string
	Proto string

	Name        string
	Category    string
	Description string
	Logo        string
	Auth        bool

	Hide    bool
	OnlyWeb bool
	Direct  bool
}

func getEntries(cli cli, conf config) ([]entry, error) {
	var entries []entry

	containers, err := cli.getContainers()
	if err != nil {
		return entries, err
	}

	for _, container := range containers {
		entry := entry{}

		// Check if enabled
		if val, ok := container.Labels[fmt.Sprintf("%s.enabled", labelPrefix)]; ok && val != "true" {
			continue
		} else if !ok && !conf.autoEnabled {
			continue
		}

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

		// Name
		entry.Name = "unknown"
		if val, ok := container.Labels[fmt.Sprintf("%s.name", labelPrefix)]; ok {
			entry.Name = val
		} else if len(container.Names) > 0 {
			entry.Name = strings.TrimPrefix(container.Names[0], "/")
		}

		// Host
		entry.Host = fmt.Sprintf("%s%s", entry.Name, ".dev")
		if val, ok := container.Labels[fmt.Sprintf("%s.host", labelPrefix)]; ok {
			entry.Host = val
		} else if val, ok := container.Labels["traefik.frontend.rule"]; ok && conf.traefikMode {
			val = strings.TrimPrefix(val, "Host:")
			array := strings.Split(val, ",")
			if len(array) > 0 {
				entry.Host = array[0]
			}
		}

		// Protocol
		entry.Proto = "http"
		if val, ok := container.Labels[fmt.Sprintf("%s.proto", labelPrefix)]; ok {
			entry.Proto = val
		} else if val, ok := container.Labels["traefik.frontend.entryPoints"]; ok && conf.traefikMode {
			array := strings.Split(val, ",")
			if len(array) > 0 {
				entry.Proto = array[0]
			}
		}

		// Auth
		entry.Auth = false
		if val, ok := container.Labels[fmt.Sprintf("%s.auth", labelPrefix)]; ok && val == "true" {
			entry.Auth = true
		} else if val, ok := container.Labels["traefik.frontend.auth.basic"]; ok && val != "" && conf.traefikMode {
			entry.Auth = true
		}

		// Category
		entry.Category = defaultCategory
		if val, ok := container.Labels[fmt.Sprintf("%s.category", labelPrefix)]; ok {
			entry.Category = strings.ToLower(val)
		}

		// Logo
		if val, ok := container.Labels[fmt.Sprintf("%s.logo", labelPrefix)]; ok {
			entry.Logo = val
		}

		// Description
		if val, ok := container.Labels[fmt.Sprintf("%s.description", labelPrefix)]; ok {
			entry.Description = val
		}

		// Hide
		entry.Hide = false
		if val, ok := container.Labels[fmt.Sprintf("%s.hide", labelPrefix)]; ok && val == "true" {
			entry.Hide = true
		}

		// Only web
		entry.OnlyWeb = false
		if val, ok := container.Labels[fmt.Sprintf("%s.onlyweb", labelPrefix)]; ok && val == "true" {
			entry.OnlyWeb = true
		}

		// Direct
		entry.Direct = false
		if val, ok := container.Labels[fmt.Sprintf("%s.direct", labelPrefix)]; ok && val == "true" {
			entry.Direct = true
		}

		entries = append(entries, entry)
	}
	return entries, nil
}
