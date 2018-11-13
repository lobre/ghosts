package main

import (
	"fmt"
	"strings"
)

const labelPrefix string = "ghosts"
const defaultCategory string = "apps"

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

func getEntries(docker docker, config config, ids ...string) ([]entry, error) {
	var entries []entry

	containers, err := docker.getContainers(ids...)
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
		} else if val, ok := container.Labels["traefik.frontend.rule"]; ok && config.TraefikMode {
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
		} else if val, ok := container.Labels["traefik.frontend.entryPoints"]; ok && config.TraefikMode {
			array := strings.Split(val, ",")
			if len(array) > 0 {
				entry.Proto = array[0]
			}
		}

		// Auth
		entry.Auth = false
		if val, ok := container.Labels[fmt.Sprintf("%s.auth", labelPrefix)]; ok && val == "true" {
			entry.Auth = true
		} else if val, ok := container.Labels["traefik.frontend.auth.basic"]; ok && val != "" && config.TraefikMode {
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

func (e entry) URL(config config) string {
	var host, port string

	if e.Direct || e.WebDirect || (!config.ProxyMode && !config.TraefikMode) {
		// Direct mode

		// Use container IP if hosts are not generated and in direct mode
		if e.WebDirect || config.NoHosts || e.NoHosts {
			host = e.IP
			port = e.Port
		} else {
			host = e.Hosts[0]
			port = e.Port
		}
	} else {
		// Proxy mode

		if e.Proto == "http" {
			host = e.Hosts[0]
			port = "80"
		} else {
			host = e.Hosts[0]
			port = "443"
		}

	}

	return fmt.Sprintf("%s://%s:%s", e.Proto, host, port)
}
