package docker

import (
	"context"
	"os"

	"github.com/docker/docker/api/types/container"
)

// SpawnWithOwnVolumes spawns a new docker container with the same volumes as the current container
func SpawnWithOwnVolumes(image string, cmd []string) (waitBody container.ContainerWaitOKBody, err error) {
	// get information about the current container
	info, err := GetContainerInfo()
	if err != nil {
		return
	}

	// create context and configuration
	ctx := context.Background()
	createConfig := container.Config{
		Image: image,
		Cmd:   cmd,
	}
	hostConfig := container.HostConfig{
		AutoRemove:  true,
		VolumesFrom: []string{info.ID},
	}

	// and run in the foreground
	return StartForeground(ctx, &createConfig, &hostConfig, nil, "", os.Stdout, os.Stderr)
}
