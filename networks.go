package main

import (
	"context"

	"github.com/docker/docker/api/types/network"
)

type networksProcessor struct {
	docker  docker
	config  config
	em      entriesManager
	proxyID string
}

func newNetworksProcessor(docker docker, config config, em entriesManager, proxy string) (networksProcessor, error) {
	var np networksProcessor
	np.docker = docker
	np.config = config
	np.em = em

	var err error
	np.proxyID, err = docker.containerIDFromName(proxy)
	if err != nil {
		return np, err
	}

	return np, nil
}

func (n networksProcessor) init() error {
	return n.connect()
}

func (n networksProcessor) startEvent(id string) error {
	return n.connect(id)
}

func (n networksProcessor) dieEvent(id string) error {
	return n.disconnect(id)
}

func (n networksProcessor) connect(ids ...string) error {
	entries, err := n.em.get(ids...)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.NoNetAutoConnect {
			continue
		}

		hasNetwork, err := n.docker.hasNetwork(n.proxyID, entry.NetworkID)
		if err != nil {
			return err
		}

		if !hasNetwork {
			err := n.docker.NetworkConnect(context.Background(), entry.NetworkID, n.proxyID, &network.EndpointSettings{})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (n networksProcessor) disconnect(ids ...string) error {
	entries, err := n.em.get(ids...)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		hasNetwork, err := n.docker.hasNetwork(n.proxyID, entry.NetworkID)
		if err != nil {
			return err
		}

		if hasNetwork {
			err := n.docker.NetworkDisconnect(context.Background(), entry.NetworkID, n.proxyID, false)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
