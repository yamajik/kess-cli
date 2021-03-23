package dapr

import (
	"github.com/dapr/cli/pkg/standalone"
)

func StandaloneInstall(runtimeVersion string, dashboardVersion string) error {
	if err := standalone.Init(runtimeVersion, dashboardVersion, "", true); err != nil {
		return err
	}
	return nil
}

func StandaloneUninstall() error {
	if err := standalone.Uninstall(false, ""); err != nil {
		return err
	}
	return nil
}

// func StandaloneRun() error {
// 	output, err := standalone.Run(&standalone.RunConfig{})
// 	if err != nil {
// 		return err
// 	}
// }
