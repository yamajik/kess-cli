package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yamajik/kess/runtimes"
)

var (
	DockerCMD = &cobra.Command{
		Use: "docker",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
			runtimeConfig.Type = "docker"
			runtimeConfig.Docker.Debug = debug
			runtime, err = runtimes.New(runtimeConfig)
			return err
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
)

func init() {
	RootCMD.AddCommand(DockerCMD)
}
