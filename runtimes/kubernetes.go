package runtimes

import (
	"context"

	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type KubernetesRuntime struct {
	client *kubernetes.Clientset
	config *KubernetesRuntimeConfig
}

type KubernetesRuntimeConfig struct {
	Debug          bool
	KubeconfigPath string
	MasterUrl      string
}

func NewKubernetesRuntime(config KubernetesRuntimeConfig) (*KubernetesRuntime, error) {
	var (
		restconfig *rest.Config
		err        error
	)
	switch {
	case config.MasterUrl != "" || config.KubeconfigPath != "":
		restconfig, err = clientcmd.BuildConfigFromFlags(config.MasterUrl, config.KubeconfigPath)
	default:
		restconfig, err = rest.InClusterConfig()
	}
	if err != nil {
		return nil, errors.WithStack(err)
	}

	c, err := kubernetes.NewForConfig(restconfig)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	r := KubernetesRuntime{
		client: c,
		config: &config,
	}
	return &r, nil
}

func (r *KubernetesRuntime) Install(ctx context.Context, options RuntimeInstallOptions) error {
	return nil
}

func (r *KubernetesRuntime) Uninstall(ctx context.Context, options RuntimeUninstallOptions) error {
	return nil
}

func (r *KubernetesRuntime) Run(ctx context.Context, options RuntimeRunOptions) error {
	return nil
}

func (r *KubernetesRuntime) Remove(ctx context.Context, options RuntimeRemoveOptions) error {
	return nil
}
