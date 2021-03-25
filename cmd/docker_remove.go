package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/yamajik/kess/runtimes"
)

var (
	dockerRemoveOptions runtimes.RuntimeRemoveOptions

	DockerRemoveCMD = &cobra.Command{
		Use:     "remove",
		Aliases: []string{"rm"},
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			for _, appid := range args {
				dockerRemoveOptions.AppID = appid
				if err := runtime.Remove(ctx, dockerRemoveOptions); err != nil {
					return err
				}
			}
			return nil
		},
	}
)

func init() {
	DockerCMD.AddCommand(DockerRemoveCMD)
}
