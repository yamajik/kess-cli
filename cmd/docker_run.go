package cmd

import (
	"context"

	"github.com/dapr/cli/pkg/standalone"
	"github.com/spf13/cobra"
	"github.com/yamajik/kess/dapr"
	"github.com/yamajik/kess/runtimes"
)

var (
	dockerRunOptions runtimes.RuntimeRunOptions

	DockerRunCMD = &cobra.Command{
		Use: "run",
		RunE: func(cmd *cobra.Command, args []string) error {
			dockerRunOptions.Arguments = args
			ctx := context.Background()
			return runtime.Run(ctx, dockerRunOptions)
		},
	}
)

func init() {
	DockerRunCMD.PersistentFlags().StringVarP(&dockerRunOptions.AppID, "app-id", "a", "", "The id for your application, used for service discovery")
	DockerRunCMD.MarkPersistentFlagRequired("app-id")
	DockerRunCMD.PersistentFlags().StringVarP(&dockerRunOptions.AppImage, "app-image", "i", "", "The image your application used")
	DockerRunCMD.PersistentFlags().IntVarP(&dockerRunOptions.AppPort, "app-port", "p", -1, "The port your application is listening on")
	DockerRunCMD.PersistentFlags().StringVarP(&dockerRunOptions.ConfigFile, "config", "c", dapr.DefaultConfigFilePath(), "Dapr configuration file")
	DockerRunCMD.PersistentFlags().IntVarP(&dockerRunOptions.HTTPPort, "dapr-http-port", "H", -1, "The http port for Dapr to listen on")
	DockerRunCMD.PersistentFlags().IntVarP(&dockerRunOptions.GRPCPort, "dapr-grpc-port", "G", -1, "The gRPC port for Dapr to listen on")
	DockerRunCMD.PersistentFlags().BoolVar(&dockerRunOptions.EnableProfiling, "enable-profiling", false, "Enable pprof profiling via an HTTP endpoint")
	DockerRunCMD.PersistentFlags().IntVarP(&dockerRunOptions.ProfilePort, "profile-port", "", -1, "The port for the profile server to listen on")
	DockerRunCMD.PersistentFlags().StringVarP(&dockerRunOptions.LogLevel, "log-level", "", "info", "The log verbosity. Valid values are: debug, info, warn, error, fatal, or panic")
	DockerRunCMD.PersistentFlags().IntVarP(&dockerRunOptions.MaxConcurrency, "app-max-concurrency", "", -1, "The concurrency level of the application, otherwise is unlimited")
	DockerRunCMD.PersistentFlags().StringVarP(&dockerRunOptions.Protocol, "app-protocol", "P", "http", "The protocol (gRPC or HTTP) Dapr uses to talk to the application")
	DockerRunCMD.PersistentFlags().StringVarP(&dockerRunOptions.ComponentsPath, "components-path", "d", standalone.DefaultComponentsDirPath(), "The path for components directory")
	DockerRunCMD.PersistentFlags().StringVarP(&dockerRunOptions.PlacementHost, "placement-host-address", "", "localhost", "The host on which the placement service resides")
	DockerRunCMD.PersistentFlags().BoolVar(&dockerRunOptions.AppSSL, "app-ssl", false, "Enable https when Dapr invokes the application")
	DockerRunCMD.PersistentFlags().IntVarP(&dockerRunOptions.MetricsPort, "metrics-port", "M", -1, "The port of metrics on dapr")
	DockerRunCMD.PersistentFlags().StringVarP(&dockerRunOptions.AppPwd, "pwd", "", "", "The dir to run cmd in")
	DockerRunCMD.PersistentFlags().IntVarP(&dockerRunOptions.AppWaitTimeoutInSeconds, "wait-timeout", "", 60, "The timeout in second to wait for app start")
	DockerCMD.AddCommand(DockerRunCMD)
}
