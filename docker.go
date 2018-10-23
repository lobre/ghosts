package main

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

type cli struct {
	*client.Client
}

// Init docker cli
func newDockerCli() (cli, error) {
	err := setDockerApiVersion()
	if err != nil {
		return cli{}, err
	}

	dockerCli, err := client.NewEnvClient()
	if err != nil {
		return cli{}, err
	}
	return cli{dockerCli}, nil
}

// Get the list of running containers
func (cli cli) getContainers() (containers []types.Container, err error) {
	containers, err = cli.ContainerList(context.Background(), types.ContainerListOptions{
		All: true,
	})
	if err != nil {
		return nil, err
	}
	return containers, nil
}

// Listen for docker events
func (cli cli) listenContainers() (<-chan events.Message, <-chan error) {
	filter := filters.NewArgs()
	filter.Add("type", "container")
	filter.Add("event", "start")
	filter.Add("event", "die")

	return cli.Events(context.Background(), types.EventsOptions{
		Filters: filter,
	})
}

func setDockerApiVersion() error {
	cmd := exec.Command("docker", "version", "--format", "{{.Server.APIVersion}}")
	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput

	err := cmd.Run()
	if err != nil {
		return err
	}
	apiVersion := strings.TrimSpace(string(cmdOutput.Bytes()))
	os.Setenv("DOCKER_API_VERSION", apiVersion)
	return nil
}
