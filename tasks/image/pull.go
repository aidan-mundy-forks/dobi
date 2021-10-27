package image

import (
	cont "context"
	"io"
	"os"
	"time"

	"github.com/dnephin/dobi/tasks/context"
	docker_types "github.com/docker/docker/api/types"
)

// RunPull builds or pulls an image if it is out of date
func RunPull(ctx *context.ExecuteContext, t *Task, _ bool) (bool, error) {
	record, err := getImageRecord(recordPath(ctx, t.config))
	switch {
	case !t.config.Pull.Required(record.LastPull):
		t.logger().Debugf("Pull not required")
		return false, nil
	case err != nil:
		t.logger().Warnf("Failed to get image record: %s", err)
	}

	pullTag := func(tag string) error {
		return pullImage(ctx, t, tag)
	}
	if err := t.ForEachRemoteTag(ctx, pullTag); err != nil {
		return false, err
	}

	image, err := GetImage(ctx, t.config)
	if err != nil {
		return false, err
	}
	record = imageModifiedRecord{LastPull: now(), ImageID: image.ID}

	if err := updateImageRecord(recordPath(ctx, t.config), record); err != nil {
		t.logger().Warnf("Failed to update image record: %s", err)
	}

	t.logger().Info("Pulled")
	return true, nil
}

func now() *time.Time {
	now := time.Now()
	return &now
}

func pullImage(ctx *context.ExecuteContext, t *Task, imageTag string) error {
	return Stream(os.Stdout, func(out io.Writer) error {
		_, err := ctx.Client.ImagePull(cont.Background(), imageTag, docker_types.ImagePullOptions{})
		return err
	})
}
