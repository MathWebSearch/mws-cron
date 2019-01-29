package docker

import (
	"context"
	"errors"
	"io/ioutil"
	"strings"

	"github.com/docker/docker/api/types"
)

// GetContainerInfo gets information about the current container
func GetContainerInfo() (info types.ContainerJSON, err error) {

	// get own docker id
	myID, err := GetDockerID()
	if err != nil {
		return
	}

	// and call inspect on it
	info, err = dockerClient.ContainerInspect(context.Background(), myID)
	return
}

// GetDockerID gets the ID of the currently running docker container
func GetDockerID() (name string, err error) {
	// read the name of the first cpu
	contents, err := ioutil.ReadFile("/proc/1/cpuset")
	if err != nil {
		return
	}

	// turn it into a string and make sure it starts with '/docker/'
	strcontents := string(contents)
	if !strings.HasPrefix(strcontents, "/docker/") {
		err = errors.New("Not inside a docker container. ")
		return
	}

	// and remove everything after /docker/
	name = strings.TrimSpace(strings.TrimPrefix(strcontents, "/docker/"))
	return
}
