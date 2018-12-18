package main

import (
	"github.com/lobre/goodhosts"
)

type hostsProcessor struct {
	config config
	em     entriesManager
}

func newHostsProcessor(config config, em entriesManager) hostsProcessor {
	return hostsProcessor{config, em}
}

func (h hostsProcessor) init() error {
	return h.add()
}

func (h hostsProcessor) startEvent(id string) error {
	return h.add(id)
}

func (h hostsProcessor) dieEvent(id string) error {
	return h.remove(id)
}

func (h hostsProcessor) add(ids ...string) error {
	entries, err := h.em.get(ids...)
	if err != nil {
		return err
	}

	hosts, err := goodhosts.NewHosts()
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.NoHosts || entry.WebDirect {
			continue
		}

		ip := h.config.ProxyIP
		if entry.Direct || !h.config.ProxyMode {
			ip = entry.IP
		}

		for _, segment := range entry.Segments {
			for _, host := range segment.Hosts {
				hosts.Remove("", host)
				hosts.Add(ip, host)
			}
		}
	}

	if err := hosts.Flush(h.config.HostsForceCRLF); err != nil {
		return err
	}

	return nil
}

func (h hostsProcessor) remove(ids ...string) error {
	entries, err := h.em.get(ids...)
	if err != nil {
		return err
	}

	hosts, err := goodhosts.NewHosts()
	if err != nil {
		return err
	}

	for _, entry := range entries {
		ip := h.config.ProxyIP
		if entry.Direct || !h.config.ProxyMode {
			ip = entry.IP
		}

		for _, segment := range entry.Segments {
			for _, host := range segment.Hosts {
				if hosts.Has(ip, host) {
					isUnique, err := h.isUnique(host, entry)
					if err != nil {
						return err
					}
					if h.config.ProxyMode && !entry.Direct && !isUnique {
						continue
					}
					hosts.Remove(ip, host)
				}
			}
		}
	}

	if err := hosts.Flush(h.config.HostsForceCRLF); err != nil {
		return err
	}

	return nil
}

// Check if host exist in another container
func (h hostsProcessor) isUnique(host string, entry entry) (bool, error) {
	entries, err := h.em.get()
	if err != nil {
		return false, err
	}
	for _, e := range entries {
		if e.Name != entry.Name {
			for _, s := range e.Segments {
				for _, h := range s.Hosts {
					if h == host {
						return false, nil
					}
				}
			}
		}
	}
	return true, nil
}
