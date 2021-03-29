package dapr

import (
	"os"

	"github.com/dapr/cli/pkg/print"
	"github.com/dapr/cli/pkg/standalone"
)

type DashboardRunConfig struct {
	Port int
}

func (c *DashboardRunConfig) Default() error {
	if c.Port == 0 {
		c.Port = DefaultDashboardPort
	}
	return nil
}

func DashboardRun(config *DashboardRunConfig) {
	if err := config.Default(); err != nil {
		print.FailureStatusEvent(os.Stderr, err.Error())
	}

	err := standalone.NewDashboardCmd(config.Port).Run()
	if err != nil {
		print.FailureStatusEvent(os.Stdout, "Dapr dashboard not found. Is Dapr installed?")
	}
}
