package image

import (
	cont "context"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/tasks/context"
	"github.com/dnephin/dobi/utils/fs"
	"github.com/docker/cli/cli/command/image/build"
	docker_types "github.com/docker/docker/api/types"
	docker_time "github.com/docker/docker/api/types/time"
	docker "github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/pkg/errors"
)

// RunBuild builds an image if it is out of date
func RunBuild(ctx *context.ExecuteContext, t *Task, hasModifiedDeps bool) (bool, error) {
	if !hasModifiedDeps {
		stale, err := buildIsStale(ctx, t)
		switch {
		case err != nil:
			return false, err
		case !stale:
			t.logger().Info("is fresh")
			return false, nil
		}
	}
	t.logger().Debug("is stale")

	if !t.config.IsBuildable() {
		return false, errors.Errorf(
			"%s is not buildable, missing required fields", t.name.Resource())
	}

	if err := buildImage(ctx, t); err != nil {
		return false, err
	}

	image, err := GetImage(ctx, t.config)
	if err != nil {
		return false, err
	}

	record := imageModifiedRecord{ImageID: image.ID}
	if err := updateImageRecord(recordPath(ctx, t.config), record); err != nil {
		t.logger().Warnf("Failed to update image record: %s", err)
	}
	t.logger().Info("Created")
	return true, nil
}

// TODO: this cyclo problem should be fixed
// nolint: gocyclo
func buildIsStale(ctx *context.ExecuteContext, t *Task) (bool, error) {
	image, err := GetImage(ctx, t.config)
	if docker.IsErrNotFound(err) {
		t.logger().Debug("Image does not exist")
		return true, nil
	}
	if err != nil {
		return true, err
	}

	paths := []string{t.config.Context}
	// TODO: polymorphic config for different types of images
	if t.config.Steps != "" && ctx.ConfigFile != "" {
		paths = append(paths, ctx.ConfigFile)
	}

	excludes, err := build.ReadDockerignore(t.config.Context)
	if err != nil {
		t.logger().Warnf("Failed to read .dockerignore file.")
	}
	excludes = append(excludes, ".dobi")

	mtime, err := fs.LastModified(&fs.LastModifiedSearch{
		Root:     absPath(ctx.WorkingDir, t.config.Context),
		Excludes: excludes,
		Paths:    paths,
	})
	if err != nil {
		t.logger().Warnf("Failed to get last modified time of cont.")
		return true, err
	}

	record, err := getImageRecord(recordPath(ctx, t.config))
	if err != nil {
		t.logger().Warnf("Failed to get image record: %s", err)
		ts, err := docker_time.GetTimestamp(image.Created, time.Now())
		if err != nil {
			return true, err
		}
		sec, nsec, err := docker_time.ParseTimestamps(ts, 0)
		if err != nil {
			return true, err
		}
		if time.Unix(sec, nsec).Before(mtime) {
			t.logger().Debug("Image older than context")
			return true, nil
		}
		return false, nil
	}

	if image.ID != record.ImageID || record.Info.ModTime().Before(mtime) {
		t.logger().Debug("Image record older than context")
		return true, nil
	}
	return false, nil
}

func absPath(path string, wd string) string {
	if filepath.IsAbs(path) {
		return filepath.Clean(path)
	}
	return filepath.Join(wd, path)
}

func buildImage(ctx *context.ExecuteContext, t *Task) error {
	var err error

	if t.config.Steps != "" {
		err = t.buildImageFromSteps(ctx)
	} else {
		err = t.buildImageFromDockerfile(ctx)
	}
	if err != nil {
		return err
	}

	image, err := GetImage(ctx, t.config)
	if err != nil {
		return err
	}
	record := imageModifiedRecord{ImageID: image.ID}
	return updateImageRecord(recordPath(ctx, t.config), record)
}

func (t *Task) buildImageFromDockerfile(ctx *context.ExecuteContext) error {
	return Stream(os.Stdout, func(out io.Writer) error {
		opts := t.commonBuildImageOptions(ctx, out)
		opts.Dockerfile = t.config.Dockerfile
		file, err := os.Open(ctx.WorkingDir)
		if err != nil {
			return err
		}

		_, err = ctx.Client.ImageBuild(cont.Background(), file, opts)
		return errors.WithMessage(err, "tesadfjhsdkjst")
		return err
	})
}

func (t *Task) commonBuildImageOptions(
	ctx *context.ExecuteContext,
	out io.Writer,
) docker_types.ImageBuildOptions {
	return docker_types.ImageBuildOptions{
		Tags:           []string{GetImageName(ctx, t.config)},
		BuildArgs:      buildArgs(t.config.Args),
		Target:         t.config.Target,
		PullParent:     t.config.PullBaseImageOnBuild,
		NetworkMode:    t.config.NetworkMode,
		CacheFrom:      t.config.CacheFrom,
		SuppressOutput: ctx.Settings.Quiet,
	}
}

func buildArgs(args map[string]*string) map[string]*string {
	out := map[string]*string{}
	for key, value := range args {
		out[key] = value
	}
	return out
}

func (t *Task) buildImageFromSteps(ctx *context.ExecuteContext) error {
	buildContext, dockerfile, err := getBuildContext(t.config)
	if err != nil {
		return err
	}
	return Stream(os.Stdout, func(out io.Writer) error {
		opts := t.commonBuildImageOptions(ctx, out)
		opts.Dockerfile = dockerfile
		_, err := ctx.Client.ImageBuild(cont.Background(), buildContext, opts)
		return err
	})
}

func getBuildContext(config *config.ImageConfig) (io.Reader, string, error) {
	contextDir := config.Context
	excludes, err := build.ReadDockerignore(contextDir)
	if err != nil {
		return nil, "", err
	}
	if err = build.ValidateContextDirectory(contextDir, excludes); err != nil {
		return nil, "", err

	}
	buildCtx, err := archive.TarWithOptions(contextDir, &archive.TarOptions{
		ExcludePatterns: excludes,
	})
	if err != nil {
		return nil, "", err
	}
	dockerfileCtx := ioutil.NopCloser(strings.NewReader(config.Steps))
	return build.AddDockerfileToBuildContext(dockerfileCtx, buildCtx)
}
