package image

import (
	cont "context"
	"fmt"

	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/tasks/context"
)

// RunTag builds or pulls an image if it is out of date
func RunTag(ctx *context.ExecuteContext, t *Task, _ bool) (bool, error) {
	tag := func(tag string) error {
		return tagImage(ctx, t.config, tag)
	}
	if err := t.ForEachTag(ctx, tag); err != nil {
		return false, err
	}
	t.logger().Info("Tagged")
	return true, nil
}

func tagImage(ctx *context.ExecuteContext, config *config.ImageConfig, imageTag string) error {
	canonicalImageTag := GetImageName(ctx, config)
	if imageTag == canonicalImageTag {
		return nil
	}

	err := ctx.Client.ImageTag(cont.Background(), canonicalImageTag, imageTag)
	if err != nil {
		return fmt.Errorf("failed to add tag %q: %s", imageTag, err)
	}
	return nil
}
