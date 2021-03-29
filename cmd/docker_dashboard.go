package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/yamajik/kess/dapr"
	"github.com/yamajik/kess/runtimes"
)

var (
	dockerDashboardOptions runtimes.RuntimeDashboardOptions

	DockerDashboardCMD = &cobra.Command{
		Use:     "dashboard",
		Aliases: []string{"web"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			return runtime.Dashboard(ctx, dockerDashboardOptions)
		},
	}
)

func init() {
	DockerDashboardCMD.PersistentFlags().IntVarP(&dockerDashboardOptions.Port, "port", "p", dapr.DefaultDashboardPort, "The local port on which to serve Dapr dashboard")
	DockerCMD.AddCommand(DockerDashboardCMD)
}
