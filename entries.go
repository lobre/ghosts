package main

import (
	"fmt"
	"strings"
)

const labelPrefix string = "ghosts"
const defaultCategory string = "apps"

type entriesManager struct {
	docker docker
	config config
}

type entry struct {
	Hosts     []string
	IP        string
	NetworkID string
	Port      string
	Proto     string

	Name        string
	Category    string
	Description string
	Logo        string
	Auth        bool

	NoWeb            bool
	NoHosts          bool
	NoNetAutoConnect bool
	Direct           bool
	WebDirect        bool
}

func newEntriesManager(docker docker, config config) entriesManager {
	return entriesManager{docker, config}
}

func (em entriesManager) get(ids ...string) ([]entry, error) {
	var entries []entry

	containers, err := em.docker.getContainers(ids...)
	if err != nil {
		return entries, err
	}

	for _, container := range containers {
		entry := entry{}

		// Host
		if val, ok := container.Labels[fmt.Sprintf("%s.hosts", labelPrefix)]; ok {
			array := strings.Split(val, ",")
			if len(array) > 0 {
				entry.Hosts = array
			}
		} else if val, ok := container.Labels["traefik.frontend.rule"]; ok && em.config.TraefikMode {
			val = strings.TrimPrefix(val, "Host:")
			array := strings.Split(val, ",")
			if len(array) > 0 {
				entry.Hosts = array
			}
		}

		// Skip if no hosts
		if len(entry.Hosts) == 0 {
			continue
		}

		// Take the IP of the first network and the network ID
		for _, n := range container.NetworkSettings.Networks {
			entry.IP = n.IPAddress
			entry.NetworkID = n.NetworkID
			break
		}

		// Port
		if val, ok := container.Labels[fmt.Sprintf("%s.port", labelPrefix)]; ok {
			entry.Port = val
		} else {
			// Take the first port exposed
			for _, p := range container.Ports {
				entry.Port = fmt.Sprint(p.PrivatePort)
				break
			}
		}

		// Name
		entry.Name = "unknown"
		if val, ok := container.Labels[fmt.Sprintf("%s.name", labelPrefix)]; ok {
			entry.Name = val
		} else if len(container.Names) > 0 {
			entry.Name = strings.TrimPrefix(container.Names[0], "/")
		}

		// Protocol
		entry.Proto = "http"
		if val, ok := container.Labels[fmt.Sprintf("%s.proto", labelPrefix)]; ok {
			entry.Proto = val
		} else if val, ok := container.Labels["traefik.frontend.entryPoints"]; ok && em.config.TraefikMode {
			array := strings.Split(val, ",")
			if len(array) > 0 {
				entry.Proto = array[0]
			}
		}

		// Auth
		entry.Auth = false
		if val, ok := container.Labels[fmt.Sprintf("%s.auth", labelPrefix)]; ok && val == "true" {
			entry.Auth = true
		} else if val, ok := container.Labels["traefik.frontend.auth.basic"]; ok && val != "" && em.config.TraefikMode {
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

		// No Web
		entry.NoWeb = false
		if val, ok := container.Labels[fmt.Sprintf("%s.noweb", labelPrefix)]; ok && val == "true" {
			entry.NoWeb = true
		}

		// No Hosts
		entry.NoHosts = false
		if val, ok := container.Labels[fmt.Sprintf("%s.nohosts", labelPrefix)]; ok && val == "true" {
			entry.NoHosts = true
		}

		// No Net Auto Connect
		entry.NoNetAutoConnect = false
		if val, ok := container.Labels[fmt.Sprintf("%s.nonetautoconnect", labelPrefix)]; ok && val == "true" {
			entry.NoNetAutoConnect = true
		}

		// Direct
		entry.Direct = false
		if val, ok := container.Labels[fmt.Sprintf("%s.direct", labelPrefix)]; ok && val == "true" {
			entry.Direct = true
		}

		// Web Direct
		entry.WebDirect = false
		if val, ok := container.Labels[fmt.Sprintf("%s.webdirect", labelPrefix)]; ok && val == "true" {
			entry.WebDirect = true
		}

		entries = append(entries, entry)
	}
	return entries, nil
}

func (em entriesManager) URLS(e entry) []string {
	var urls []string
	var port string

	// Check specific port if direct mode
	if e.Direct || e.WebDirect || (!em.config.ProxyMode && !em.config.TraefikMode) {
		port = fmt.Sprintf(":%s", e.Port)

		// Use container IP if hosts are not generated and in direct mode
		if e.WebDirect || em.config.NoHosts || e.NoHosts {
			return []string{fmt.Sprintf("%s://%s%s", e.Proto, e.IP, e.Port)}
		}
	}

	for _, host := range e.Hosts {
		urls = append(urls, fmt.Sprintf("%s://%s%s", e.Proto, host, port))
	}

	return urls
}
