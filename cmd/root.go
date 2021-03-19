package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yamajik/kess/runtimes"
)

var (
	debug         bool
	runtimeConfig runtimes.RuntimeConfig
	runtime       runtimes.Runtime

	RootCMD = &cobra.Command{
		Use:   "kess",
		Short: "Kess CLI",
		Long:  "Kess DAR based on Dapr",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
)

func init() {
	RootCMD.PersistentFlags().BoolVarP(&debug, "debug", "D", false, "Debug")
}

func Execute(version string) {
	RootCMD.Version = version
	if err := RootCMD.Execute(); err != nil {
		if debug {
			fmt.Fprintf(os.Stderr, "%+v", err)
		} else {
			fmt.Fprintln(os.Stderr, err)

		}
		os.Exit(1)
	}
}
