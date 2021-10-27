package image

import (
	cont "context"
	"fmt"
	"strings"

	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/logging"
	"github.com/dnephin/dobi/tasks/context"
	docker_types "github.com/docker/docker/api/types"
)

const (
	defaultRepo = "https://index.docker.io/v1/"
)

// GetImage returns the image created by an image config
func GetImage(ctx *context.ExecuteContext, conf *config.ImageConfig) (docker_types.ImageInspect, error) {
	image, _, err := ctx.Client.ImageInspectWithRaw(cont.Background(), GetImageName(ctx, conf))
	return image, err
}

// GetImageName returns the image name for an image config
func GetImageName(ctx *context.ExecuteContext, conf *config.ImageConfig) string {
	return fmt.Sprintf("%s:%s", conf.Image, GetCanonicalTag(ctx, conf))
}

// GetCanonicalTag returns the canonical tag for an image config
func GetCanonicalTag(ctx *context.ExecuteContext, conf *config.ImageConfig) string {
	if len(conf.Tags) > 0 {
		return conf.Tags[0]
	}
	return ctx.Env.Unique()
}

func parseAuthRepo(image string) string {
	return splitHostname(image)
}

// Copied from github.com/docker/docker/reference/reference.go
// That package is conflicting with other dependencies, so it can't be imported
// at this time.
func splitHostname(name string) string {
	i := strings.IndexRune(name, '/')
	if i == -1 || (!strings.ContainsAny(name[:i], ".:") && name[:i] != "localhost") {
		logging.Log.Debugf("Using default registry %q", defaultRepo)
		return defaultRepo
	}
	return name[:i]
}
