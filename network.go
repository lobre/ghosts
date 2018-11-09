package main

import (
	"context"

	"github.com/docker/docker/api/types/network"
)

func networkConnect(proxy string, docker docker, config config) error {
	proxyID, err := docker.containerIDFromName(proxy)
	if err != nil {
		return err
	}

	// Initialize connection
	if err := connect(proxyID, docker, config); err != nil {
		return err
	}

	// Listen to Docker events
	msgCh, errCh := docker.listenContainers()
	for {
		select {
		case msg := <-msgCh:
			switch msg.Action {
			case "start":
				if err := connect(proxyID, docker, config, msg.ID); err != nil {
					return err
				}
			case "die":
				if err := disconnect(proxyID, docker, config, msg.ID); err != nil {
					return err
				}
			}
		case err := <-errCh:
			return err
		}
	}
}

func connect(proxyID string, docker docker, config config, ids ...string) error {
	entries, err := getEntries(docker, config, ids...)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.NoNetAutoConnect {
			continue
		}

		hasNetwork, err := docker.hasNetwork(proxyID, entry.NetworkID)
		if err != nil {
			return err
		}

		if !hasNetwork {
			err := docker.NetworkConnect(context.Background(), entry.NetworkID, proxyID, &network.EndpointSettings{})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func disconnect(proxyID string, docker docker, config config, ids ...string) error {
	entries, err := getEntries(docker, config, ids...)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		hasNetwork, err := docker.hasNetwork(proxyID, entry.NetworkID)
		if err != nil {
			return err
		}

		if hasNetwork {
			err := docker.NetworkDisconnect(context.Background(), entry.NetworkID, proxyID, false)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
