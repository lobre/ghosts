package main

import (
	"fmt"

	"github.com/lextoumbourou/goodhosts"
)

const proxyIp = "127.0.0.1"

func hosts() error {
	hosts, err := goodhosts.NewHosts()
	if err != nil {
		return err
	}

	entriesMap, err := entries()
	if err != nil {
		return nil
	}

	for _, entries := range entriesMap {
		for _, entry := range entries {
			if hosts.Has(proxyIp, entry.Host) {
				fmt.Println("Entry exists")
			} else {
				fmt.Println("Entry does not exist")
			}
		}
	}

	return nil
}
