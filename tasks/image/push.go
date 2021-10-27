package image

import (
	cont "context"
	"io"
	"os"

	"github.com/dnephin/dobi/tasks/context"
	docker_types "github.com/docker/docker/api/types"
)

// RunPush pushes an image to the registry
func RunPush(ctx *context.ExecuteContext, t *Task, _ bool) (bool, error) {
	pushTag := func(tag string) error {
		return pushImage(ctx, tag)
	}
	if err := t.ForEachRemoteTag(ctx, pushTag); err != nil {
		return false, err
	}
	t.logger().Info("Pushed")
	return true, nil
}

func pushImage(ctx *context.ExecuteContext, tag string) error {
	return Stream(os.Stdout, func(out io.Writer) error {
		_, err := ctx.Client.ImagePush(cont.Background(), tag, docker_types.ImagePushOptions{})
		return err
	})
}
