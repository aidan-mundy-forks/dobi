package client

import (
	cont "context"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/docker/docker/api/types"
	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	networktypes "github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/registry"
	docker "github.com/docker/docker/client"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
)

//go:generate mockgen -source iface.go -destination mock_iface.go -package client

// DockerClient is the Docker API Client interface used by tasks
type DockerClient interface {
	docker.ConfigAPIClient
	ContainerAPIClient
	DistributionAPIClient
	ImageAPIClient
	docker.NodeAPIClient
	NetworkAPIClient
	docker.PluginAPIClient
	docker.ServiceAPIClient
	docker.SwarmAPIClient
	docker.SecretAPIClient
	SystemAPIClient
	docker.VolumeAPIClient
	ClientVersion() string
	DaemonHost() string
	HTTPClient() *http.Client
	ServerVersion(ctx cont.Context) (types.Version, error)
	NegotiateAPIVersion(ctx cont.Context)
	NegotiateAPIVersionPing(types.Ping)
	DialHijack(ctx cont.Context, url, proto string, meta map[string][]string) (net.Conn, error)
	Dialer() func(cont.Context) (net.Conn, error)
	Close() error
}

// ContainerAPIClient defines API client methods for the containers
type ContainerAPIClient interface {
	ContainerAttach(ctx cont.Context, containerString string, options types.ContainerAttachOptions) (types.HijackedResponse, error)
	ContainerCommit(ctx cont.Context, containerString string, options types.ContainerCommitOptions) (types.IDResponse, error)
	ContainerCreate(ctx cont.Context, config *containertypes.Config, hostConfig *containertypes.HostConfig, networkingConfig *networktypes.NetworkingConfig, platform *specs.Platform, containerName string) (containertypes.ContainerCreateCreatedBody, error)
	ContainerDiff(ctx cont.Context, containerString string) ([]containertypes.ContainerChangeResponseItem, error)
	ContainerExecAttach(ctx cont.Context, execID string, config types.ExecStartCheck) (types.HijackedResponse, error)
	ContainerExecCreate(ctx cont.Context, containerString string, config types.ExecConfig) (types.IDResponse, error)
	ContainerExecInspect(ctx cont.Context, execID string) (types.ContainerExecInspect, error)
	ContainerExecResize(ctx cont.Context, execID string, options types.ResizeOptions) error
	ContainerExecStart(ctx cont.Context, execID string, config types.ExecStartCheck) error
	ContainerExport(ctx cont.Context, containerString string) (io.ReadCloser, error)
	ContainerInspect(ctx cont.Context, containerString string) (types.ContainerJSON, error)
	ContainerInspectWithRaw(ctx cont.Context, containerString string, getSize bool) (types.ContainerJSON, []byte, error)
	ContainerKill(ctx cont.Context, container, signal string) error
	ContainerList(ctx cont.Context, options types.ContainerListOptions) ([]types.Container, error)
	ContainerLogs(ctx cont.Context, containerString string, options types.ContainerLogsOptions) (io.ReadCloser, error)
	ContainerPause(ctx cont.Context, containerString string) error
	ContainerRemove(ctx cont.Context, containerString string, options types.ContainerRemoveOptions) error
	ContainerRename(ctx cont.Context, container, newContainerName string) error
	ContainerResize(ctx cont.Context, containerString string, options types.ResizeOptions) error
	ContainerRestart(ctx cont.Context, containerString string, timeout *time.Duration) error
	ContainerStatPath(ctx cont.Context, container, path string) (types.ContainerPathStat, error)
	ContainerStats(ctx cont.Context, containerString string, stream bool) (types.ContainerStats, error)
	ContainerStatsOneShot(ctx cont.Context, containerString string) (types.ContainerStats, error)
	ContainerStart(ctx cont.Context, containerString string, options types.ContainerStartOptions) error
	ContainerStop(ctx cont.Context, containerString string, timeout *time.Duration) error
	ContainerTop(ctx cont.Context, containerString string, arguments []string) (containertypes.ContainerTopOKBody, error)
	ContainerUnpause(ctx cont.Context, containerString string) error
	ContainerUpdate(ctx cont.Context, containerString string, updateConfig containertypes.UpdateConfig) (containertypes.ContainerUpdateOKBody, error)
	ContainerWait(ctx cont.Context, containerString string, condition containertypes.WaitCondition) (<-chan containertypes.ContainerWaitOKBody, <-chan error)
	CopyFromContainer(ctx cont.Context, container, srcPath string) (io.ReadCloser, types.ContainerPathStat, error)
	CopyToContainer(ctx cont.Context, container, path string, content io.Reader, options types.CopyToContainerOptions) error
	ContainersPrune(ctx cont.Context, pruneFilters filters.Args) (types.ContainersPruneReport, error)
}

// DistributionAPIClient defines API client methods for the registry
type DistributionAPIClient interface {
	DistributionInspect(ctx cont.Context, image, encodedRegistryAuth string) (registry.DistributionInspect, error)
}

// ImageAPIClient defines API client methods for the images
type ImageAPIClient interface {
	ImageBuild(ctx cont.Context, context io.Reader, options types.ImageBuildOptions) (types.ImageBuildResponse, error)
	BuildCachePrune(ctx cont.Context, opts types.BuildCachePruneOptions) (*types.BuildCachePruneReport, error)
	BuildCancel(ctx cont.Context, id string) error
	ImageCreate(ctx cont.Context, parentReference string, options types.ImageCreateOptions) (io.ReadCloser, error)
	ImageHistory(ctx cont.Context, imageString string) ([]image.HistoryResponseItem, error)
	ImageImport(ctx cont.Context, source types.ImageImportSource, ref string, options types.ImageImportOptions) (io.ReadCloser, error)
	ImageInspectWithRaw(ctx cont.Context, imageString string) (types.ImageInspect, []byte, error)
	ImageList(ctx cont.Context, options types.ImageListOptions) ([]types.ImageSummary, error)
	ImageLoad(ctx cont.Context, input io.Reader, quiet bool) (types.ImageLoadResponse, error)
	ImagePull(ctx cont.Context, ref string, options types.ImagePullOptions) (io.ReadCloser, error)
	ImagePush(ctx cont.Context, ref string, options types.ImagePushOptions) (io.ReadCloser, error)
	ImageRemove(ctx cont.Context, imageString string, options types.ImageRemoveOptions) ([]types.ImageDeleteResponseItem, error)
	ImageSearch(ctx cont.Context, term string, options types.ImageSearchOptions) ([]registry.SearchResult, error)
	ImageSave(ctx cont.Context, images []string) (io.ReadCloser, error)
	ImageTag(ctx cont.Context, image, ref string) error
	ImagesPrune(ctx cont.Context, pruneFilter filters.Args) (types.ImagesPruneReport, error)
}

// NetworkAPIClient defines API client methods for the networks
type NetworkAPIClient interface {
	NetworkConnect(ctx cont.Context, network, containerString string, config *networktypes.EndpointSettings) error
	NetworkCreate(ctx cont.Context, name string, options types.NetworkCreate) (types.NetworkCreateResponse, error)
	NetworkDisconnect(ctx cont.Context, network, containerString string, force bool) error
	NetworkInspect(ctx cont.Context, network string, options types.NetworkInspectOptions) (types.NetworkResource, error)
	NetworkInspectWithRaw(ctx cont.Context, network string, options types.NetworkInspectOptions) (types.NetworkResource, []byte, error)
	NetworkList(ctx cont.Context, options types.NetworkListOptions) ([]types.NetworkResource, error)
	NetworkRemove(ctx cont.Context, network string) error
	NetworksPrune(ctx cont.Context, pruneFilter filters.Args) (types.NetworksPruneReport, error)
}

// SystemAPIClient defines API client methods for the system
type SystemAPIClient interface {
	Events(ctx cont.Context, options types.EventsOptions) (<-chan events.Message, <-chan error)
	Info(ctx cont.Context) (types.Info, error)
	RegistryLogin(ctx cont.Context, auth types.AuthConfig) (registry.AuthenticateOKBody, error)
	DiskUsage(ctx cont.Context) (types.DiskUsage, error)
	Ping(ctx cont.Context) (types.Ping, error)
}
