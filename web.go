package main

import (
	"html/template"
	"net"
	"net/http"
	"strings"
)

type appHandler struct {
	config config
	em     entriesManager
}

func (h appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("index.html").Funcs(template.FuncMap{
		"capitalize": strings.Title,
		"upper":      strings.ToUpper,
	}).ParseFiles("index.html")

	entries, err := h.getPreparedEntries()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = tmpl.Execute(w, struct {
		Config  config
		Entries map[string][]entry
	}{
		h.config,
		entries,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Separate entries by categories and replace urls/port if needed
func (h appHandler) getPreparedEntries() (map[string][]entry, error) {
	categories := make(map[string][]entry)

	entries, err := h.em.get()
	if err != nil {
		return categories, err
	}

	for i, entry := range entries {
		if !entry.NoWeb {
			for name, segment := range entry.Segments {
				for j, u := range segment.URLS {
					if entry.WebDirect ||
						((!h.config.ProxyMode || entry.Direct) && (h.config.NoHosts || entry.NoHosts)) {
						// Replace IP and Port in URL
						// TODO put out of for to have only 1 url and 1 segment per entry
						host := net.JoinHostPort(entry.IP, segment.Port)
						entries[i].Segments[name].URLS[j].Host = host
					} else if !h.config.ProxyMode || entry.Direct {
						// Replace Port in URL
						host := net.JoinHostPort(u.Host, segment.Port)
						entries[i].Segments[name].URLS[j].Host = host
					}
				}
			}
			categories[entry.Category] = append(categories[entry.Category], entries[i])
		}
	}
	return categories, nil
}
