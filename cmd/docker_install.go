package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

var (
	DockerInstallCMD = &cobra.Command{
		Use:     "install",
		Aliases: []string{"init", "setup"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			return runtime.Install(ctx)
		},
	}
)

func init() {
	DockerCMD.AddCommand(DockerInstallCMD)
}
