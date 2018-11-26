package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/docker/docker/api/types"
)

const (
	labelPrefix     string = "ghosts"
	defaultCategory string = "apps"
)

type entriesManager struct {
	docker docker
	config config
}

type segment struct {
	URLS []string
	Port string
}

type entry struct {
	Segments  map[string]segment
	IP        string
	NetworkID string

	Name        string
	Category    string
	Description string
	Logo        string
	Auth        bool

	NoWeb            bool
	NoHosts          bool
	NoNetAutoConnect bool

	Direct    bool
	WebDirect bool
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

		// Skip if no segments
		entry.Segments = parseSegments(container)
		if len(entry.Segments) == 0 {
			continue
		}

		// Take the IP of the first network and the network ID
		for _, n := range container.NetworkSettings.Networks {
			entry.IP = n.IPAddress
			entry.NetworkID = n.NetworkID
			break
		}

		// Name
		entry.Name = "unknown"
		if val, ok := container.Labels[fmt.Sprintf("%s.name", labelPrefix)]; ok {
			entry.Name = val
		} else if len(container.Names) > 0 {
			entry.Name = strings.TrimPrefix(container.Names[0], "/")
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

func parseSegments(container types.Container) map[string]segment {
	var segments map[string]segment

	rURLS := regexp.MustCompile(fmt.Sprintf("%s\\.([a-zA-Z]+)\\.urls", labelPrefix))
	rPort := regexp.MustCompile(fmt.Sprintf("%s\\.([a-zA-Z]+)\\.port", labelPrefix))

	var urlsMap map[string][]string
	var portMap map[string]string

	for key, value := range container.Labels {
		// Segment URLS
		if match := rURLS.FindStringSubmatch(key); match != nil {
			name := match[1]
			urlsMap[name] = strings.Split(value, ",")
		}
		// Segment port
		if match := rPort.FindStringSubmatch(key); match != nil {
			name := match[1]
			portMap[name] = value
		}
		// Default URLS
		if key == fmt.Sprintf("%s.urls", labelPrefix) {
			urlsMap[""] = strings.Split(value, ",")
		}
		// Default Port
		if key == fmt.Sprintf("%s.port", labelPrefix) {
			portMap[""] = value
		}
	}

	// Take the first port exposed
	defaultPort := "80"
	for _, p := range container.Ports {
		defaultPort = fmt.Sprint(p.PrivatePort)
		break
	}

	// Bind urls and port
	for name, urls := range urlsMap {
		s := segment{URLS: urls}
		if port, ok := portMap[name]; ok {
			s.Port = port
		} else {
			s.Port = defaultPort
		}
		segments[name] = s
	}

	return segments
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
