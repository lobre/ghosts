package main

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

var cli *client.Client

// Init docker cli
func initDockerCli() error {
	var err error

	cli, err = client.NewClientWithOpts(client.WithVersion("1.38"))
	if err != nil {
		return err
	}
	return nil
}

// Get the list of running containers
func getContainers() (containers []types.Container, err error) {
	containers, err = cli.ContainerList(context.Background(), types.ContainerListOptions{
		All: false,
	})
	return
}

// Listen for docker events
func listenContainers() (<-chan events.Message, <-chan error) {
	filter := filters.NewArgs()
	filter.Add("type", "container")
	filter.Add("event", "start")
	filter.Add("event", "die")

	return cli.Events(context.Background(), types.EventsOptions{
		Filters: filter,
	})
}
