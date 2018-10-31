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

	NoWeb     bool
	NoHosts   bool
	Direct    bool
	WebDirect bool
}

func getEntries(docker docker, config config, ids ...string) ([]entry, error) {
	var entries []entry

	containers, err := docker.getContainers(ids...)
	if err != nil {
		return entries, err
	}

	for _, container := range containers {
		entry := entry{}

		// Check if enabled
		if val, ok := container.Labels[fmt.Sprintf("%s.enabled", labelPrefix)]; ok && val != "true" {
			continue
		}

		// Host
		if val, ok := container.Labels[fmt.Sprintf("%s.host", labelPrefix)]; ok {
			entry.Host = val
		} else if val, ok := container.Labels["traefik.frontend.rule"]; ok && config.traefikMode {
			val = strings.TrimPrefix(val, "Host:")
			array := strings.Split(val, ",")
			if len(array) > 0 {
				entry.Host = array[0]
			}
		} else {
			continue
		}

		// Take the IP of the first network
		for _, n := range container.NetworkSettings.Networks {
			entry.IP = n.IPAddress
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
		} else if val, ok := container.Labels["traefik.frontend.entryPoints"]; ok && config.traefikMode {
			array := strings.Split(val, ",")
			if len(array) > 0 {
				entry.Proto = array[0]
			}
		}

		// Auth
		entry.Auth = false
		if val, ok := container.Labels[fmt.Sprintf("%s.auth", labelPrefix)]; ok && val == "true" {
			entry.Auth = true
		} else if val, ok := container.Labels["traefik.frontend.auth.basic"]; ok && val != "" && config.traefikMode {
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

	if e.Direct || e.WebDirect || (!config.proxyMode && !config.traefikMode) {
		// Direct mode

		// Use container IP if hosts are not generated and in direct mode
		if e.WebDirect || config.noHosts || e.NoHosts {
			host = e.IP
			port = e.Port
		} else {
			host = e.Host
			port = e.Port
		}
	} else {
		// Proxy mode

		if e.Proto == "http" {
			host = e.Host
			port = "80"
		} else {
			host = e.Host
			port = "443"
		}

	}

	return fmt.Sprintf("%s://%s:%s", e.Proto, host, port)
}
