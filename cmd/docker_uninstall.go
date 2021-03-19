package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

var (
	DockerUninstallCMD = &cobra.Command{
		Use:     "uninstall",
		Aliases: []string{"unsetup"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			return runtime.Uninstall(ctx)
		},
	}
)

func init() {
	DockerCMD.AddCommand(DockerUninstallCMD)
}
