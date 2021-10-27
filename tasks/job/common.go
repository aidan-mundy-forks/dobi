package job

import (
	cont "context"
	"fmt"

	"github.com/dnephin/dobi/tasks/client"
	"github.com/dnephin/dobi/tasks/context"
	docker_types "github.com/docker/docker/api/types"
	docker "github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
)

// containerName returns the name of the container
func containerName(ctx *context.ExecuteContext, name string) string {
	return fmt.Sprintf("%s-%s", ctx.Env.Unique(), name)
}

// removeContainer removes a container by ID, and logs a warning if the remove
// fails.
func removeContainer(
	logger *log.Entry,
	client client.DockerClient,
	containerID string,
) (bool, error) {
	logger.Debug("Removing container")
	err := client.ContainerRemove(cont.Background(), containerID, docker_types.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	})
	if docker.IsErrNotFound(err) {
		return false, nil
	}
	if err == nil {
		return true, nil
	}

	logger.WithFields(log.Fields{"container": containerID}).Warnf(
		"Failed to remove container: %s", err)
	return false, err
}
