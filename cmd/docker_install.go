package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/yamajik/kess/runtimes"
)

var (
	dockerInstallOptions runtimes.RuntimeInstallOptions

	DockerInstallCMD = &cobra.Command{
		Use:     "install",
		Aliases: []string{"init", "setup"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			return runtime.Install(ctx, dockerInstallOptions)
		},
	}
)

func init() {
	DockerInstallCMD.PersistentFlags().StringVarP(&dockerInstallOptions.RuntimeVersion, "runtime-version", "", "latest", "The version of the Dapr runtime to install, for example: 1.0.0")
	DockerInstallCMD.PersistentFlags().StringVarP(&dockerInstallOptions.DashboardVersion, "dashboard-version", "", "latest", "The version of the Dapr dashboard to install, for example: 1.0.0")
	DockerCMD.AddCommand(DockerInstallCMD)
}
