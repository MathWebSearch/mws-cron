package docker

import "github.com/docker/docker/client"

// dockerClient is a client to the Docker API
var dockerClient *client.Client

func init() {
	var err error
	dockerClient, err = client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.26"))
	if err != nil {
		panic(err)
	}
}
