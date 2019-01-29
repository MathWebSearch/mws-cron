package docker

import (
	"fmt"
	"time"

	"github.com/docker/docker/api/types/filters"
)

// UpdateMWS updates the MWS job
func UpdateMWS(containerLabel string, retryFunc func(), retryDelay *time.Duration) {
	fmt.Println("Running 'mathwebsearch/mws-indexer' container to re-index harvests. ")

	wait, err := SpawnWithOwnVolumes("mathwebsearch/mws-indexer", nil)

	if err != nil || wait.StatusCode != 0 {
		if retryFunc == nil {
			fmt.Printf("Running indexer failed. ")
			return
		}

		// schedule the code to run again
		go (func() {
			fmt.Printf("Running indexer failed, will try again in %s. \n", (*retryDelay).String())
			time.Sleep(*retryDelay)
			retryFunc()
		})()

		return
	}

	fmt.Printf("Restarting containers with label %s\n", containerLabel)
	filters := filters.NewArgs()
	filters.Add("label", containerLabel)
	timeout := 15 * time.Second
	RestartContainers(filters, &timeout)
}
