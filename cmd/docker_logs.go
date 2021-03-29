package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/yamajik/kess/runtimes"
)

var (
	dockerLogsOptions runtimes.RuntimeLogsOptions

	DockerLogsCMD = &cobra.Command{
		Use:  "logs",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dockerLogsOptions.AppID = args[0]
			ctx := context.Background()
			return runtime.Logs(ctx, dockerLogsOptions)
		},
	}
)

func init() {
	DockerLogsCMD.PersistentFlags().BoolVarP(&dockerLogsOptions.Follow, "follow", "f", false, "Follow logs")
	DockerLogsCMD.PersistentFlags().StringVarP(&dockerLogsOptions.Tail, "tail", "", "", "Tail logs")
	DockerCMD.AddCommand(DockerLogsCMD)
}
