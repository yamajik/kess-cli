package runtimes

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/fatih/structs"
	"github.com/pkg/errors"
	"github.com/valyala/fasttemplate"
	"github.com/yamajik/kess/dapr"
)

var (
	DefaultDockerRuntimeNetwork = "kess"
	DefaultDockerRuntimeVolumes = []string{"kess-configs"}

	DefaultDockerRuntimeToolsName  = "kess-tools-{Time}"
	DefaultDockerRuntimeToolsImage = "alpine:latest"
	DefaultDockerRuntimeToolsCmd   = []string{"sleep", "infinity"}

	DefaultDockerRuntimeRedisName         = "kess-system-redis"
	DefaultDockerRuntimeRedisImage        = "redis:alpine"
	DefaultDockerRuntimeRedisCmd          = []string{"redis-server"}
	DefaultDockerRuntimeRedisNetwork      = DefaultDockerRuntimeNetwork
	DefaultDockerRuntimeRedisPorts        = []string{"50003:6379"}
	DefaultDockerRuntimeRedisExternalHost = "localhost:50003"
	DefaultDockerRuntimeRedisInternalHost = "kess-system-redis:6379"
	DefaultDockerRuntimeRedisPassword     = ""

	DefaultDockerRuntimeZipkinName         = "kess-system-zipkin"
	DefaultDockerRuntimeZipkinImage        = "openzipkin/zipkin:latest"
	DefaultDockerRuntimeZipkinCmd          = []string{""}
	DefaultDockerRuntimeZipkinNetwork      = DefaultDockerRuntimeNetwork
	DefaultDockerRuntimeZipkinPorts        = []string{"50004:9411"}
	DefaultDockerRuntimeZipkinExternalHost = "localhost:50004"
	DefaultDockerRuntimeZipkinInternalHost = "kess-system-zipkin:9411"

	DefaultDockerRuntimePlacementName         = "kess-system-placement"
	DefaultDockerRuntimePlacementImage        = "daprio/dapr"
	DefaultDockerRuntimePlacementCmd          = []string{"./placement"}
	DefaultDockerRuntimePlacementNetwork      = DefaultDockerRuntimeNetwork
	DefaultDockerRuntimePlacementPorts        = []string{"50005:50005"}
	DefaultDockerRuntimePlacementExternalHost = "localhost:50005"
	DefaultDockerRuntimePlacementInternalHost = "kess-system-placement:50005"

	DefaultDockerRuntimeIngressName  = "kess-system-ingress"
	DefaultDockerRuntimeIngressImage = "daprio/daprd:edge"
	DefaultDockerRuntimeIngressCmd   = []string{
		"./daprd",
		"--placement-host-address", "kess-system-placement:50005",
		"--components-path", "/kess-configs/components",
		"--config", "/kess-configs/config.yaml",
		"--app-id", "ingress",
		"--dapr-grpc-port", "50001",
		"--dapr-http-port", "50002",
	}
	DefaultDockerRuntimeIngressNetwork = DefaultDockerRuntimeNetwork
	DefaultDockerRuntimeIngressPorts   = []string{"50001:50001", "50002:50002"}
	DefaultDockerRuntimeIngressVolumes = []string{"kess-configs:/kess-configs"}

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

type DockerRuntimeToolsConfig struct {
	Name  string
	Image string
	Cmd   []string
}

func (c *DockerRuntimeToolsConfig) Default() error {
	if c.Name == "" {
		c.Name = DefaultDockerRuntimeToolsName
	}
	if c.Image == "" {
		c.Image = DefaultDockerRuntimeToolsImage
	}
	if len(c.Cmd) == 0 {
		c.Cmd = DefaultDockerRuntimeToolsCmd
	}
	return nil
}

type DockerRuntimeRedisConfig struct {
	Name         string
	Image        string
	Cmd          []string
	Network      string
	Ports        []string
	ExternalHost string
	InternalHost string
	Password     string
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
	if c.ExternalHost == "" {
		c.ExternalHost = DefaultDockerRuntimeRedisExternalHost
	}
	if c.InternalHost == "" {
		c.InternalHost = DefaultDockerRuntimeRedisInternalHost
	}
	if c.Password == "" {
		c.Password = DefaultDockerRuntimeRedisPassword
	}
	return nil
}

type DockerRuntimeZipkinConfig struct {
	Name         string
	Image        string
	Cmd          []string
	Network      string
	Ports        []string
	ExternalHost string
	InternalHost string
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
	if c.ExternalHost == "" {
		c.ExternalHost = DefaultDockerRuntimeZipkinExternalHost
	}
	if c.InternalHost == "" {
		c.InternalHost = DefaultDockerRuntimeZipkinInternalHost
	}
	return nil
}

type DockerRuntimePlacementConfig struct {
	Name         string
	Image        string
	Cmd          []string
	Network      string
	Ports        []string
	ExternalHost string
	InternalHost string
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
	if c.ExternalHost == "" {
		c.ExternalHost = DefaultDockerRuntimePlacementExternalHost
	}
	if c.InternalHost == "" {
		c.InternalHost = DefaultDockerRuntimePlacementInternalHost
	}
	return nil
}

type DockerRuntimeIngressConfig struct {
	Name    string
	Image   string
	Cmd     []string
	Network string
	Ports   []string
	Volumes []string
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
	if c.Network == "" {
		c.Network = DefaultDockerRuntimeIngressNetwork
	}
	if len(c.Ports) == 0 {
		c.Ports = DefaultDockerRuntimeIngressPorts
	}
	if len(c.Volumes) == 0 {
		c.Volumes = DefaultDockerRuntimeIngressVolumes
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
	Tools     DockerRuntimeToolsConfig
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
	if err := c.Tools.Default(); err != nil {
		return err
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
	externalDaprConfigs := r.getDaprConfigs(r.config.Zipkin.ExternalHost, r.config.Redis.ExternalHost, r.config.Redis.Password)
	if err := externalDaprConfigs.Save(); err != nil {
		return err
	}

	for _, volume := range r.config.Volumes {
		if err := r.createVolume(ctx, volume); err != nil {
			return err
		}

	}

	configsVolume := r.findConfigsVolume(r.config.Ingress.Volumes)
	if configsVolume != "" {
		internalDaprConfigs := r.getDaprConfigs(r.config.Zipkin.InternalHost, r.config.Redis.InternalHost, r.config.Redis.Password)
		buf, err := internalDaprConfigs.Buffer()
		if err != nil {
			return err
		}
		if err := r.copyToVolume(ctx, configsVolume, buf); err != nil {
			return err
		}
	}

	if err := r.createNetwork(ctx, r.config.Network); err != nil {
		return err
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
		Volumes: r.config.Ingress.Volumes,
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
	if err := r.client.ContainerRemove(ctx, name, types.ContainerRemoveOptions{Force: true, RemoveVolumes: true}); err != nil {
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

func (r *DockerRuntime) copyToVolume(ctx context.Context, volume string, reader io.Reader) error {
	dist := strings.SplitN(volume, ":", 2)[1]

	containerName := r.renderName(r.config.Tools.Name, map[string]interface{}{"Time": strconv.FormatInt(time.Now().Unix(), 10)})
	if err := r.runContainer(ctx, DockerRuntimeRunContainerOptions{
		Name:    containerName,
		Image:   r.config.Tools.Image,
		Cmd:     r.config.Tools.Cmd,
		Volumes: []string{volume},
		Labels: r.labels(map[string]string{
			"kess-tools": "",
		}),
	}); err != nil {
		return err
	}

	if err := r.client.CopyToContainer(ctx, containerName, dist, reader, types.CopyToContainerOptions{}); err != nil {
		return errors.WithStack(err)
	}

	if err := r.removeContainer(ctx, containerName); err != nil {
		return err
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

func (r *DockerRuntime) findConfigsVolume(volumes []string) string {
	for _, volume := range volumes {
		if strings.HasPrefix(volume, "kess-configs") {
			return volume
		}
	}
	return ""
}

func (r *DockerRuntime) getDaprConfigs(zipkinHost string, redisHost string, redisPassword string) *dapr.Configs {
	daprConfigs := dapr.DefaultConfigs()
	daprConfigs.SetConfiguration(dapr.CreateConfiguration("kess", dapr.ConfigurationSpec{
		Tracing: dapr.ConfigurationSpecTracing{
			SamplingRate: "1",
			Zipkin: dapr.ConfigurationSpecTracingZipkin{
				EndpointAddress: fmt.Sprintf("http://%s/api/v2/spans", zipkinHost),
			},
		},
	}))
	daprConfigs.SetComponents([]dapr.Component{
		dapr.CreateRedisStateStoreComponent("statestore", dapr.RedisStateStoreComponentOptions{
			Host:            redisHost,
			Password:        redisPassword,
			ActorStateStore: true,
		}),
		dapr.CreateRedisPubsubComponent("pubsub", dapr.RedisPubsubComponentOptions{
			Host:     redisHost,
			Password: redisPassword,
		}),
	})
	return daprConfigs
}
