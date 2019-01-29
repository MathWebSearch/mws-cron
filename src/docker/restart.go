package docker

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

// RestartContainers restarts the containers matched by the filter
func RestartContainers(filters filters.Args, timeout *time.Duration) (err error) {
	ctx := context.Background()

	// find all the matching containers
	var containers []types.Container
	if containers, err = dockerClient.ContainerList(ctx, types.ContainerListOptions{Filters: filters}); err != nil {
		return
	}

	// restart all the docker containers
	for _, container := range containers {
		fmt.Printf("Restarting container %q\n", container.Names[0])

		if err = dockerClient.ContainerRestart(ctx, container.ID, timeout); err != nil {
			return
		}
	}

	return
}
