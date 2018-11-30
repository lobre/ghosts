package main

import (
	"html/template"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/iancoleman/strcase"
)

type appHandler struct {
	config config
	em     entriesManager
}

func (h appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("index.html").Funcs(template.FuncMap{
		"upper":   strings.ToUpper,
		"spacify": spacify,
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
						host := entry.IP
						if segment.Port != "80" && segment.Port != "443" {
							host = net.JoinHostPort(host, segment.Port)
						}
						entries[i].Segments[name].URLS[j].Host = host

						// Quickfix to keep only one url
						tmpSeg := entries[i].Segments[name]
						tmpSeg.URLS = []url.URL{entries[i].Segments[name].URLS[j]}
						entries[i].Segments[name] = tmpSeg
						break
					} else if !h.config.ProxyMode || entry.Direct {
						// Replace Port in URL
						host := u.Host
						if strings.Contains(host, ":") {
							var err error
							host, _, err = net.SplitHostPort(u.Host)
							if err != nil {
								continue
							}
						}
						if segment.Port != "80" && segment.Port != "443" {
							host = net.JoinHostPort(host, segment.Port)
						}
						entries[i].Segments[name].URLS[j].Host = host
					}
				}
			}
			categories[entry.Category] = append(categories[entry.Category], entries[i])
		}
	}
	return categories, nil
}

func spacify(s string) string {
	return strings.Title(strcase.ToDelimited(s, ' '))
}
