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
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			return runtime.Remove(ctx, dockerRemoveOptions)
		},
	}
)

func init() {
	DockerRemoveCMD.PersistentFlags().StringVarP(&dockerRemoveOptions.AppID, "app-id", "a", "", "The id for your application, used for service discovery")
	DockerRemoveCMD.MarkPersistentFlagRequired("name")
	DockerCMD.AddCommand(DockerRemoveCMD)
}
