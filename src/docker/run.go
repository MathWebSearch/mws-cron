package docker

import (
	"context"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/pkg/stdcopy"
)

// StartForeground starts a container in the foreground
func StartForeground(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, containerName string, dstStd io.Writer, dstErr io.Writer) (waitBody container.ContainerWaitOKBody, err error) {

	// create the container with the given options
	var resp container.ContainerCreateCreatedBody
	if resp, err = dockerClient.ContainerCreate(ctx, config, hostConfig, networkingConfig, containerName); err != nil {
		return
	}

	// start the container
	if err = dockerClient.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return
	}

	// attach to the container
	var hijack types.HijackedResponse
	if hijack, err = dockerClient.ContainerAttach(ctx, resp.ID, types.ContainerAttachOptions{
		Stream: true,
		Stdin:  false,
		Stdout: true,
		Stderr: true,
	}); err != nil {
		return
	}

	// pass through all the input and output
	go func() {
		defer hijack.Close()
		if config.Tty {
			io.Copy(dstStd, hijack.Reader)
		} else {
			stdcopy.StdCopy(dstStd, dstErr, hijack.Reader)
		}
	}()

	// wait for it
	statusCh, errCh := dockerClient.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err = <-errCh:
	case waitBody = <-statusCh:
	}

	// and return
	return
}
