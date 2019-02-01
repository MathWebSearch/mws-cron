package docker

import (
	"context"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
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

	networkConfig := network.NetworkingConfig{
		// EndpointsConfig: duplicateNetworkConfig(&info.NetworkSettings.Networks),
	}

	// and run in the foreground
	return StartForeground(ctx, &createConfig, &hostConfig, &networkConfig, "", os.Stdout, os.Stderr)
}

// duplicateNetworkConfig copies the names of networks from old to new
func duplicateNetworkConfig(old *map[string]*network.EndpointSettings) (new map[string]*network.EndpointSettings) {
	for name, ep := range *old {
		new[name] = &network.EndpointSettings{
			NetworkID: ep.NetworkID,
		}
	}
	return
}
