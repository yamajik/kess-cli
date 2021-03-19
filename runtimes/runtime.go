package runtimes

import (
	"context"

	"github.com/pkg/errors"
)

type Runtime interface {
	Install(ctx context.Context) error
	Uninstall(ctx context.Context) error
	Run(ctx context.Context, options RuntimeRunOptions) error
	Remove(ctx context.Context, options RuntimeRemoveOptions) error
}

type RuntimeConfig struct {
	Type       string
	Slim       SlimRuntimeConfig
	Docker     DockerRuntimeConfig
	Kubernetes KubernetesRuntimeConfig
}

type RuntimeRunOptions struct {
	Name  string
	Image string
	Cmd   []string
	Port  string
}

type RuntimeRemoveOptions struct {
	Name string
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
