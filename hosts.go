package main

import (
	"github.com/lextoumbourou/goodhosts"
)

const proxyIp = "127.0.0.1"

func hosts() error {
	hosts, err := goodhosts.NewHosts()
	if err != nil {
		return err
	}

	// entriesMap, err := getEntries(nil, nil)
	// if err != nil {
	// 	return nil
	// }

	// for _, entries := range entriesMap {
	// 	for _, entry := range entries {
	// 		if hosts.Has(proxyIp, entry.Host) {
	// 			fmt.Println("Entry exists")
	// 		} else {
	// 			fmt.Println("Entry does not exist")
	// 		}
	// 	}
	// }

	hosts.Add("127.1.1.1", "facebook.com", "twitter.com")

	if err := hosts.Flush(); err != nil {
		panic(err)
	}

	return nil
}
