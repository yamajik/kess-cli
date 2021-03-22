package runtimes

import (
	"context"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/fatih/structs"
	"github.com/pkg/errors"
	"github.com/valyala/fasttemplate"
)

var (
	DefaultDockerRuntimeNetwork = "kess"
	DefaultDockerRuntimeVolumes = []string{"kess-configs"}

	DefaultDockerRuntimeRedisName    = "kess-system-redis"
	DefaultDockerRuntimeRedisImage   = "redis:alpine"
	DefaultDockerRuntimeRedisCmd     = []string{"redis-server"}
	DefaultDockerRuntimeRedisPorts   = []string{"6379:6379"}
	DefaultDockerRuntimeRedisNetwork = DefaultDockerRuntimeNetwork

	DefaultDockerRuntimeZipkinName    = "kess-system-zipkin"
	DefaultDockerRuntimeZipkinImage   = "openzipkin/zipkin:latest"
	DefaultDockerRuntimeZipkinCmd     = []string{""}
	DefaultDockerRuntimeZipkinPorts   = []string{"9411:9411"}
	DefaultDockerRuntimeZipkinNetwork = DefaultDockerRuntimeNetwork

	DefaultDockerRuntimePlacementName    = "kess-system-placement"
	DefaultDockerRuntimePlacementImage   = "daprio/dapr"
	DefaultDockerRuntimePlacementCmd     = []string{"./placement", "--port", "50005"}
	DefaultDockerRuntimePlacementPorts   = []string{"50005:50005"}
	DefaultDockerRuntimePlacementNetwork = DefaultDockerRuntimeNetwork

	DefaultDockerRuntimeIngressName  = "kess-system-ingress"
	DefaultDockerRuntimeIngressImage = "daprio/daprd:edge"
	DefaultDockerRuntimeIngressCmd   = []string{
		"./daprd",
		"--placement-host-address", "kess-system-placement:50005",
		"--components-path", "/components",
		"--app-id", "ingress",
		"--dapr-http-port", "3500",
		"--dapr-grpc-port", "50001",
	}
	DefaultDockerRuntimeIngressPorts   = []string{"3500:3500", "50001:50001"}
	DefaultDockerRuntimeIngressNetwork = DefaultDockerRuntimeNetwork

	DefaultDockerRuntimeSidecarName  = "kess-app-{Name}-sidecar"
	DefaultDockerRuntimeSidecarImage = "daprio/daprd:edge"
	DefaultDockerRuntimeSidecarCmd   = []string{
		"./daprd",
		"--placement-host-address", "kess-system-placement:50005",
		"--components-path", "/kess-configs/components",
		"--config", "/kess-configs/config.yaml",
	}
	DefaultDockerRuntimeSidecarNetwork = "container:{Name}"
	DefaultDockerRuntimeSidecarVolumes = []string{"kess-configs:/kess-configs"}

	DefaultDockerRuntimeAppName    = "kess-app-{Name}"
	DefaultDockerRuntimeAppNetwork = DefaultDockerRuntimeNetwork
	DefaultDockerRuntimeAppVolumes = []string{}
)

type DockerRuntimeRedisConfig struct {
	Name    string
	Image   string
	Cmd     []string
	Network string
	Ports   []string
}

func (c *DockerRuntimeRedisConfig) Default() error {
	if c.Name == "" {
		c.Name = DefaultDockerRuntimeRedisName
	}
	if c.Image == "" {
		c.Image = DefaultDockerRuntimeRedisImage
	}
	if len(c.Cmd) == 0 {
		c.Cmd = DefaultDockerRuntimeRedisCmd
	}
	if len(c.Ports) == 0 {
		c.Ports = DefaultDockerRuntimeRedisPorts
	}
	if c.Network == "" {
		c.Network = DefaultDockerRuntimeRedisNetwork
	}
	return nil
}

type DockerRuntimeZipkinConfig struct {
	Name    string
	Image   string
	Cmd     []string
	Network string
	Ports   []string
}

func (c *DockerRuntimeZipkinConfig) Default() error {
	if c.Name == "" {
		c.Name = DefaultDockerRuntimeZipkinName
	}
	if c.Image == "" {
		c.Image = DefaultDockerRuntimeZipkinImage
	}
	if len(c.Cmd) == 0 {
		c.Cmd = DefaultDockerRuntimeZipkinCmd
	}
	if len(c.Ports) == 0 {
		c.Ports = DefaultDockerRuntimeZipkinPorts
	}
	if c.Network == "" {
		c.Network = DefaultDockerRuntimeZipkinNetwork
	}
	return nil
}

type DockerRuntimePlacementConfig struct {
	Name    string
	Image   string
	Cmd     []string
	Network string
	Ports   []string
}

func (c *DockerRuntimePlacementConfig) Default() error {
	if c.Name == "" {
		c.Name = DefaultDockerRuntimePlacementName
	}
	if c.Image == "" {
		c.Image = DefaultDockerRuntimePlacementImage
	}
	if len(c.Cmd) == 0 {
		c.Cmd = DefaultDockerRuntimePlacementCmd
	}
	if len(c.Ports) == 0 {
		c.Ports = DefaultDockerRuntimePlacementPorts
	}
	if c.Network == "" {
		c.Network = DefaultDockerRuntimePlacementNetwork
	}
	return nil
}

type DockerRuntimeIngressConfig struct {
	Name    string
	Image   string
	Cmd     []string
	Network string
	Ports   []string
}

func (c *DockerRuntimeIngressConfig) Default() error {
	if c.Name == "" {
		c.Name = DefaultDockerRuntimeIngressName
	}
	if c.Image == "" {
		c.Image = DefaultDockerRuntimeIngressImage
	}
	if len(c.Cmd) == 0 {
		c.Cmd = DefaultDockerRuntimeIngressCmd
	}
	if len(c.Ports) == 0 {
		c.Ports = DefaultDockerRuntimeIngressPorts
	}
	if c.Network == "" {
		c.Network = DefaultDockerRuntimeIngressNetwork
	}
	return nil
}

type DockerRuntimeSidecarConfig struct {
	Name    string
	Image   string
	Cmd     []string
	Network string
	Volumes []string
}

func (c *DockerRuntimeSidecarConfig) Default() error {
	if c.Name == "" {
		c.Name = DefaultDockerRuntimeSidecarName
	}
	if c.Image == "" {
		c.Image = DefaultDockerRuntimeSidecarImage
	}
	if len(c.Cmd) == 0 {
		c.Cmd = DefaultDockerRuntimeSidecarCmd
	}
	if c.Network == "" {
		c.Network = DefaultDockerRuntimeSidecarNetwork
	}
	if len(c.Volumes) == 0 {
		c.Volumes = DefaultDockerRuntimeSidecarVolumes
	}
	return nil
}

type DockerRuntimeAppConfig struct {
	Name    string
	Network string
	Volumes []string
}

func (c *DockerRuntimeAppConfig) Default() error {
	if c.Name == "" {
		c.Name = DefaultDockerRuntimeAppName
	}
	if c.Network == "" {
		c.Network = DefaultDockerRuntimeAppNetwork
	}
	if len(c.Volumes) == 0 {
		c.Volumes = DefaultDockerRuntimeAppVolumes
	}
	return nil
}

type DockerRuntimeConfig struct {
	Debug     bool
	Network   string
	Volumes   []string
	Redis     DockerRuntimeRedisConfig
	Zipkin    DockerRuntimeZipkinConfig
	Placement DockerRuntimePlacementConfig
	Ingress   DockerRuntimeIngressConfig
	Sidecar   DockerRuntimeSidecarConfig
	App       DockerRuntimeAppConfig
}

func (c *DockerRuntimeConfig) Default() error {
	if c.Network == "" {
		c.Network = DefaultDockerRuntimeNetwork
	}
	if len(c.Volumes) == 0 {
		c.Volumes = DefaultDockerRuntimeVolumes
	}
	if err := c.Redis.Default(); err != nil {
		return err
	}
	if err := c.Zipkin.Default(); err != nil {
		return err
	}
	if err := c.Placement.Default(); err != nil {
		return err
	}
	if err := c.Ingress.Default(); err != nil {
		return err
	}
	if err := c.Sidecar.Default(); err != nil {
		return err
	}
	if err := c.App.Default(); err != nil {
		return err
	}
	return nil
}

type DockerRuntime struct {
	client *client.Client
	config *DockerRuntimeConfig
}

func NewDockerRuntime(config DockerRuntimeConfig) (*DockerRuntime, error) {
	if err := config.Default(); err != nil {
		return nil, err
	}

	c, err := client.NewEnvClient()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	r := DockerRuntime{
		client: c,
		config: &config,
	}

	return &r, nil
}

func (r *DockerRuntime) Install(ctx context.Context) error {
	if err := r.createNetwork(ctx, r.config.Network); err != nil {
		return err
	}

	for _, volume := range r.config.Volumes {
		if err := r.createVolume(ctx, volume); err != nil {
			return err
		}
	}

	if err := r.runContainer(ctx, DockerRuntimeRunContainerOptions{
		Name:    r.config.Redis.Name,
		Image:   r.config.Redis.Image,
		Cmd:     r.config.Redis.Cmd,
		Network: r.config.Redis.Network,
		Ports:   r.config.Redis.Ports,
		Labels: r.labels(map[string]string{
			"kess-system": "redis",
		}),
	}); err != nil {
		return err
	}

	if err := r.runContainer(ctx, DockerRuntimeRunContainerOptions{
		Name:    r.config.Zipkin.Name,
		Image:   r.config.Zipkin.Image,
		Cmd:     r.config.Zipkin.Cmd,
		Network: r.config.Zipkin.Network,
		Ports:   r.config.Zipkin.Ports,
		Labels: r.labels(map[string]string{
			"kess-system": "zipkin",
		}),
	}); err != nil {
		return err
	}

	if err := r.runContainer(ctx, DockerRuntimeRunContainerOptions{
		Name:    r.config.Placement.Name,
		Image:   r.config.Placement.Image,
		Cmd:     r.config.Placement.Cmd,
		Network: r.config.Placement.Network,
		Ports:   r.config.Placement.Ports,
		Labels: r.labels(map[string]string{
			"kess-system": "placement",
		}),
	}); err != nil {
		return err
	}

	if err := r.runContainer(ctx, DockerRuntimeRunContainerOptions{
		Name:    r.config.Ingress.Name,
		Image:   r.config.Ingress.Image,
		Cmd:     r.config.Ingress.Cmd,
		Network: r.config.Ingress.Network,
		Ports:   r.config.Ingress.Ports,
		Labels: r.labels(map[string]string{
			"kess-system": "ingress",
		}),
	}); err != nil {
		return err
	}

	return nil
}

func (r *DockerRuntime) Uninstall(ctx context.Context) error {
	filters := filters.NewArgs(filters.Arg("label", "kess"))

	containers, err := r.client.ContainerList(ctx, types.ContainerListOptions{All: true, Filters: filters})
	if err != nil {
		return err
	}

	for _, container := range containers {
		if err := r.removeContainer(ctx, container.ID); err != nil {
			return err
		}
	}

	volumesResp, err := r.client.VolumeList(ctx, filters)
	if err != nil {
		return err
	}

	for _, volume := range volumesResp.Volumes {
		if err := r.removeVolume(ctx, volume.Name); err != nil {
			return err
		}
	}

	networks, err := r.client.NetworkList(ctx, types.NetworkListOptions{Filters: filters})
	if err != nil {
		return err
	}

	for _, network := range networks {
		if err := r.removeNetwork(ctx, network.ID); err != nil {
			return err
		}
	}

	return nil
}

func (r *DockerRuntime) Run(ctx context.Context, options RuntimeRunOptions) error {
	m := structs.Map(options)

	if err := r.createNetwork(ctx, r.config.App.Network); err != nil {
		return err
	}

	appContainerName := r.renderName(r.config.App.Name, m)
	if err := r.runContainer(ctx, DockerRuntimeRunContainerOptions{
		Name:    appContainerName,
		Image:   options.Image,
		Cmd:     options.Cmd,
		Network: r.config.App.Network,
		Volumes: r.config.App.Volumes,
		Labels: r.labels(map[string]string{
			"kess-app": options.Name,
		}),
	}); err != nil {
		return err
	}

	if err := r.runContainer(ctx, DockerRuntimeRunContainerOptions{
		Name:    r.renderName(r.config.Sidecar.Name, m),
		Image:   r.config.Sidecar.Image,
		Cmd:     append(r.config.Sidecar.Cmd, "--app-id", options.Name, "--app-port", options.Port),
		Network: r.renderName(r.config.Sidecar.Network, map[string]interface{}{"Name": appContainerName}),
		Volumes: r.config.Sidecar.Volumes,
		Labels: r.labels(map[string]string{
			"kess-app":         options.Name,
			"kess-app-sidecar": options.Name,
		}),
	}); err != nil {
		return err
	}

	return nil
}

func (r *DockerRuntime) Remove(ctx context.Context, options RuntimeRemoveOptions) error {
	m := structs.Map(options)

	if err := r.removeContainer(ctx, r.renderName(r.config.Sidecar.Name, m)); err != nil {
		return err
	}

	if err := r.removeContainer(ctx, r.renderName(r.config.App.Name, m)); err != nil {
		return err
	}

	return nil
}

func (r *DockerRuntime) renderName(tpl string, m map[string]interface{}) string {
	return fasttemplate.New(tpl, "{", "}").ExecuteString(m)
}

type DockerRuntimeRunContainerOptions struct {
	Name    string
	Image   string
	Cmd     []string
	Network string
	Ports   []string
	Volumes []string
	Links   []string
	Labels  map[string]string
}

func (r *DockerRuntime) runContainer(ctx context.Context, options DockerRuntimeRunContainerOptions) error {
	if _, _, err := r.client.ImageInspectWithRaw(ctx, options.Image); err != nil {
		if !client.IsErrNotFound(err) {
			return errors.WithStack(err)
		}
		reader, err := r.client.ImagePull(ctx, options.Image, types.ImagePullOptions{})
		if err != nil {
			return errors.WithStack(err)
		}
		if r.config.Debug {
			io.Copy(os.Stdout, reader)
		}
	}

	exposedports, portbindings, err := nat.ParsePortSpecs(options.Ports)
	if err != nil {
		return err
	}

	resp, err := r.client.ContainerCreate(ctx, &container.Config{
		Image:        options.Image,
		Cmd:          options.Cmd,
		ExposedPorts: exposedports,
		Labels:       options.Labels,
	}, &container.HostConfig{
		NetworkMode:   container.NetworkMode(options.Network),
		PortBindings:  portbindings,
		Binds:         options.Volumes,
		RestartPolicy: container.RestartPolicy{Name: "always"},
	}, nil, nil, options.Name)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := r.client.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *DockerRuntime) removeContainer(ctx context.Context, name string) error {
	if err := r.client.ContainerRemove(ctx, name, types.ContainerRemoveOptions{Force: true}); err != nil {
		if !client.IsErrNotFound(err) {
			return errors.WithStack(err)
		}
	}
	return nil
}

func (r *DockerRuntime) createNetwork(ctx context.Context, name string) error {
	if _, err := r.client.NetworkInspect(ctx, name, types.NetworkInspectOptions{}); err != nil {
		if _, err := r.client.NetworkCreate(ctx, name, types.NetworkCreate{Labels: r.labels(nil)}); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func (r *DockerRuntime) removeVolume(ctx context.Context, name string) error {
	if err := r.client.VolumeRemove(ctx, name, true); err != nil {
		if !client.IsErrNotFound(err) {
			return errors.WithStack(err)
		}
	}
	return nil
}

func (r *DockerRuntime) createVolume(ctx context.Context, name string) error {
	if _, err := r.client.VolumeInspect(ctx, name); err != nil {
		if _, err := r.client.VolumeCreate(ctx, volume.VolumeCreateBody{Name: name, Labels: r.labels(nil)}); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func (r *DockerRuntime) removeNetwork(ctx context.Context, name string) error {
	if err := r.client.NetworkRemove(ctx, name); err != nil {
		if !client.IsErrNotFound(err) {
			return errors.WithStack(err)
		}
	}
	return nil
}

func (r *DockerRuntime) labels(m map[string]string) map[string]string {
	l := map[string]string{"kess": ""}
	for k, v := range m {
		l[k] = v
	}
	return l
}
