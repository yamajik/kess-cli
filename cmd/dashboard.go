package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yamajik/kess/dapr"
)

var (
	dashboardRunConfig dapr.DashboardRunConfig

	DashboardCMD = &cobra.Command{
		Use:     "dashboard",
		Aliases: []string{"web"},
		RunE: func(cmd *cobra.Command, args []string) error {
			dapr.DashboardRun(&dashboardRunConfig)
			return nil
		},
	}
)

func init() {
	DashboardCMD.PersistentFlags().IntVarP(&dashboardRunConfig.Port, "port", "p", dapr.DefaultDashboardPort, "The local port on which to serve Dapr dashboard")
	RootCMD.AddCommand(DashboardCMD)
}
