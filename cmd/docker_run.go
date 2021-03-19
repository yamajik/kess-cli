package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/yamajik/kess/runtimes"
)

var (
	dockerRunOptions runtimes.RuntimeRunOptions

	DockerRunCMD = &cobra.Command{
		Use: "run",
		RunE: func(cmd *cobra.Command, args []string) error {
			dockerRunOptions.Cmd = args
			ctx := context.Background()
			return runtime.Run(ctx, dockerRunOptions)
		},
	}
)

func init() {
	DockerRunCMD.PersistentFlags().StringVarP(&dockerRunOptions.Name, "name", "n", "", "The id for your application, used for service discovery")
	DockerRunCMD.MarkPersistentFlagRequired("name")
	DockerRunCMD.PersistentFlags().StringVarP(&dockerRunOptions.Image, "image", "i", "", "The image your application used")
	DockerRunCMD.MarkPersistentFlagRequired("image")
	DockerRunCMD.PersistentFlags().StringVarP(&dockerRunOptions.Port, "port", "p", "", "The port your application is listening on")
	DockerRunCMD.MarkPersistentFlagRequired("port")
	DockerCMD.AddCommand(DockerRunCMD)
}
