package main

import (
	"fmt"
	"strings"
)

const prefix = "ghosts"
const defaultCategory = "others"

type entry struct {
	Name     string
	Protocol string
	Host     string
}

func entries() (map[string][]entry, error) {
	entries := make(map[string][]entry)

	containers, err := getContainers()
	if err != nil {
		return entries, err
	}

	for _, container := range containers {
		category := defaultCategory
		entry := entry{}

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
			entry.Protocol = val
		}

		// Host
		if val, ok := container.Labels[fmt.Sprintf("%s.host", prefix)]; ok {
			entry.Host = val
		}

		entries[category] = append(entries[category], entry)
	}
	return entries, nil
}
