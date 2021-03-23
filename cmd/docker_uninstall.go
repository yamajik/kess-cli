package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/yamajik/kess/runtimes"
)

var (
	dockerUninstallOptions runtimes.RuntimeUninstallOptions

	DockerUninstallCMD = &cobra.Command{
		Use:     "uninstall",
		Aliases: []string{"unsetup"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			return runtime.Uninstall(ctx, dockerUninstallOptions)
		},
	}
)

func init() {
	DockerCMD.AddCommand(DockerUninstallCMD)
}
