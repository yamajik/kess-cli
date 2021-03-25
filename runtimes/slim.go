package runtimes

import "context"

type SlimRuntime struct {
	config *SlimRuntimeConfig
}

type SlimRuntimeConfig struct {
	Debug bool
}

func NewSlimRuntime(config SlimRuntimeConfig) (*SlimRuntime, error) {
	r := SlimRuntime{
		config: &config,
	}
	return &r, nil
}

func (r *SlimRuntime) Install(ctx context.Context, options RuntimeInstallOptions) error {
	return nil
}

func (r *SlimRuntime) Uninstall(ctx context.Context, options RuntimeUninstallOptions) error {
	return nil
}

func (r *SlimRuntime) Run(ctx context.Context, options RuntimeRunOptions) error {
	return nil
}

func (r *SlimRuntime) Remove(ctx context.Context, options RuntimeRemoveOptions) error {
	return nil
}

func (r *SlimRuntime) Logs(ctx context.Context, options RuntimeLogsOptions) error {
	return nil
}
