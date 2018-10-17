package main

import (
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
)

const prefix = "home"
const defaultCategory = "others"

type entry struct {
	Name string
	Url  string
}

func entries(containers []types.Container) map[string][]entry {
	entries := make(map[string][]entry)

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

		// Url
		if val, ok := container.Labels[fmt.Sprintf("%s.url", prefix)]; ok {
			entry.Url = val
		}

		entries[category] = append(entries[category], entry)
	}
	return entries
}
