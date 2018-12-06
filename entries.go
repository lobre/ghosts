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
	Hosts []string
	Paths []string
	Port  string
	Proto string
}

type entry struct {
	Segments  map[string]segment
	IP        string
	NetworkID string

	Name        string
	Category    []string
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
		}

		// Category
		entry.Category = []string{defaultCategory}
		if val, ok := container.Labels[fmt.Sprintf("%s.category", labelPrefix)]; ok {
			val = strings.ToLower(val)
			entry.Category = strings.Split(val, ",")
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
	segments := make(map[string]segment)

	regex := fmt.Sprintf("%s\\.([a-zA-Z0-9_-]+)\\.", labelPrefix)

	rHosts := regexp.MustCompile(fmt.Sprint(regex, "host"))
	rPaths := regexp.MustCompile(fmt.Sprint(regex, "path"))
	rPort := regexp.MustCompile(fmt.Sprint(regex, "port"))
	rProto := regexp.MustCompile(fmt.Sprint(regex, "proto"))

	hostsMap := make(map[string][]string)
	pathsMap := make(map[string][]string)
	portMap := make(map[string]string)
	protoMap := make(map[string]string)

	for key, value := range container.Labels {
		// Segment Hosts
		if match := rHosts.FindStringSubmatch(key); match != nil {
			name := match[1]
			for _, host := range strings.Split(value, ",") {
				hostsMap[name] = append(hostsMap[name], host)
			}
		}
		// Segment Paths
		if match := rPaths.FindStringSubmatch(key); match != nil {
			name := match[1]
			for _, path := range strings.Split(value, ",") {
				pathsMap[name] = append(pathsMap[name], path)
			}
		}
		// Segment port
		if match := rPort.FindStringSubmatch(key); match != nil {
			name := match[1]
			portMap[name] = value
		}
		// Segment proto
		if match := rProto.FindStringSubmatch(key); match != nil {
			name := match[1]
			protoMap[name] = value
		}
		// Default Hosts
		if key == fmt.Sprintf("%s.host", labelPrefix) {
			for _, u := range strings.Split(value, ",") {
				hostsMap[""] = append(hostsMap[""], u)
			}
		}
		// Default Paths
		if key == fmt.Sprintf("%s.path", labelPrefix) {
			for _, u := range strings.Split(value, ",") {
				pathsMap[""] = append(pathsMap[""], u)
			}
		}
		// Default Port
		if key == fmt.Sprintf("%s.port", labelPrefix) {
			portMap[""] = value
		}
		// Default Proto
		if key == fmt.Sprintf("%s.proto", labelPrefix) {
			protoMap[""] = value
		}
	}

	// Take the first port exposed
	defaultPort := "80"
	for _, p := range container.Ports {
		defaultPort = fmt.Sprint(p.PrivatePort)
		break
	}

	// Bind to create segments
	for name, hosts := range hostsMap {
		s := segment{Hosts: hosts}

		// Bind paths
		if paths, ok := pathsMap[name]; ok {
			s.Paths = paths
		} else {
			s.Paths = []string{"/"}
		}

		// Bind port
		if port, ok := portMap[name]; ok {
			s.Port = port
		} else {
			s.Port = defaultPort
		}

		// Bind proto
		if proto, ok := protoMap[name]; ok {
			s.Proto = proto
		} else {
			s.Proto = "http"
		}

		segments[name] = s
	}

	return segments
}
