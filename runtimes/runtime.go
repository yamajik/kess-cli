package runtimes

import (
	"context"

	"github.com/pkg/errors"
	"github.com/yamajik/kess/dapr"
)

type Runtime interface {
	Install(ctx context.Context, options RuntimeInstallOptions) error
	Uninstall(ctx context.Context, options RuntimeUninstallOptions) error
	Run(ctx context.Context, options RuntimeRunOptions) error
	Remove(ctx context.Context, options RuntimeRemoveOptions) error
	Logs(ctx context.Context, options RuntimeLogsOptions) error
	Dashboard(ctx context.Context, options RuntimeDashboardOptions) error
}

type RuntimeConfig struct {
	Type       string
	Slim       SlimRuntimeConfig
	Docker     DockerRuntimeConfig
	Kubernetes KubernetesRuntimeConfig
}

type RuntimeInstallOptions struct {
	RuntimeVersion   string
	DashboardVersion string
}

type RuntimeUninstallOptions struct {
}

type RuntimeRunOptions struct {
	dapr.StandaloneRunConfig
	AppImage string
}

type RuntimeRemoveOptions struct {
	AppID string
}

type RuntimeLogsOptions struct {
	AppID  string
	Follow bool
	Tail   string
}

type RuntimeDashboardOptions struct {
	dapr.DashboardRunConfig
}

func New(config RuntimeConfig) (Runtime, error) {
	switch config.Type {
	case "slim":
		return Slim(config.Slim)
	case "docker":
		return Docker(config.Docker)
	case "kubernetes":
		return Kubernetes(config.Kubernetes)
	default:
		return nil, errors.Errorf("Unknown runtime type: %s", config.Type)
	}
}

func Slim(config SlimRuntimeConfig) (Runtime, error) {
	return NewSlimRuntime(config)
}

func Docker(config DockerRuntimeConfig) (Runtime, error) {
	return NewDockerRuntime(config)
}

func Kubernetes(config KubernetesRuntimeConfig) (Runtime, error) {
	return NewKubernetesRuntime(config)
}
