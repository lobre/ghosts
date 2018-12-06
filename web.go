package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/iancoleman/strcase"
)

type webEntry struct {
	Entry entry
	URLS  map[string][]string
}

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
		Config     config
		WebEntries map[string][]webEntry
	}{
		h.config,
		entries,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Separate entries by categories and replace urls/port if needed
func (h appHandler) getPreparedEntries() (map[string][]webEntry, error) {
	categories := make(map[string][]webEntry)

	entries, err := h.em.get()
	if err != nil {
		return categories, err
	}

	for _, entry := range entries {
		if !entry.NoWeb {

			webEntry := webEntry{entry, make(map[string][]string)}

			for name, segment := range entry.Segments {

				stop := false
				for _, host := range segment.Hosts {
					if stop {
						break
					}
					for _, path := range segment.Paths {
						if stop {
							break
						}

						port := false
						if entry.WebDirect ||
							((!h.config.ProxyMode || entry.Direct) && (h.config.NoHosts || entry.NoHosts)) {

							host = entry.IP
							port = true
							stop = true
						} else if !h.config.ProxyMode || entry.Direct {
							port = true
						}

						url := fmt.Sprintf("%s://%s", segment.Proto, host)
						if port && segment.Port != "80" && segment.Port != "443" {
							url = fmt.Sprintf("%s:%s", url, segment.Port)
						}
						url = fmt.Sprint(url, path)
						webEntry.URLS[name] = append(webEntry.URLS[name], url)
					}
				}
			}
			for _, cat := range entry.Category {
				categories[cat] = append(categories[cat], webEntry)
			}
		}
	}
	return categories, nil
}

func spacify(s string) string {
	return strings.Title(strcase.ToDelimited(s, ' '))
}
