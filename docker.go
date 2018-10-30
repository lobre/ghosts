package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

type docker struct {
	*client.Client
}

// Init docker cli
func newDockerCli() (docker, error) {
	err := setDockerAPIVersion()
	if err != nil {
		return docker{}, err
	}

	dockerCli, err := client.NewEnvClient()
	if err != nil {
		return docker{}, err
	}
	return docker{dockerCli}, nil
}

// Get the list of running containers
func (docker docker) getContainers() (containers []types.Container, err error) {
	containers, err = docker.ContainerList(context.Background(), types.ContainerListOptions{
		All: false,
	})
	if err != nil {
		return nil, err
	}
	return containers, nil
}

// Listen for docker events
func (docker docker) listenContainers() (<-chan events.Message, <-chan error) {
	filter := filters.NewArgs()
	filter.Add("type", "container")
	filter.Add("event", "start")
	filter.Add("event", "die")

	return docker.Events(context.Background(), types.EventsOptions{
		Filters: filter,
	})
}

func setDockerAPIVersion() error {
	cmd := exec.Command("docker", "version", "--format", "{{.Server.APIVersion}}")
	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	cmdErr := &bytes.Buffer{}
	cmd.Stderr = cmdErr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf(cmdErr.String())
	}
	apiVersion := strings.TrimSpace(string(cmdOutput.Bytes()))
	os.Setenv("DOCKER_API_VERSION", apiVersion)
	return nil
}
