package cmd

import (
	"github.com/dapr/cli/pkg/standalone"
	"github.com/spf13/cobra"
	"github.com/yamajik/kess/dapr"
)

var (
	runConfig dapr.StandaloneRunConfig

	RunCMD = &cobra.Command{
		Use: "run",
		RunE: func(cmd *cobra.Command, args []string) error {
			runConfig.Arguments = args
			dapr.StandaloneRun(&runConfig)
			return nil
		},
	}
)

func init() {
	RunCMD.PersistentFlags().StringVarP(&runConfig.AppID, "app-id", "a", "", "The id for your application, used for service discovery")
	RunCMD.MarkPersistentFlagRequired("app-id")
	RunCMD.PersistentFlags().IntVarP(&runConfig.AppPort, "app-port", "p", dapr.DefaultRandomPort, "The port your application is listening on")
	RunCMD.PersistentFlags().StringVarP(&runConfig.ConfigFile, "config", "c", dapr.DefaultConfigFilePath(), "Dapr configuration file")
	RunCMD.PersistentFlags().IntVarP(&runConfig.HTTPPort, "dapr-http-port", "H", dapr.DefaultRandomPort, "The http port for Dapr to listen on")
	RunCMD.PersistentFlags().IntVarP(&runConfig.GRPCPort, "dapr-grpc-port", "G", dapr.DefaultRandomPort, "The gRPC port for Dapr to listen on")
	RunCMD.PersistentFlags().BoolVar(&runConfig.EnableProfiling, "enable-profiling", false, "Enable pprof profiling via an HTTP endpoint")
	RunCMD.PersistentFlags().IntVarP(&runConfig.ProfilePort, "profile-port", "", dapr.DefaultRandomPort, "The port for the profile server to listen on")
	RunCMD.PersistentFlags().StringVarP(&runConfig.LogLevel, "log-level", "", "info", "The log verbosity. Valid values are: debug, info, warn, error, fatal, or panic")
	RunCMD.PersistentFlags().IntVarP(&runConfig.MaxConcurrency, "app-max-concurrency", "", dapr.DefaultRandomPort, "The concurrency level of the application, otherwise is unlimited")
	RunCMD.PersistentFlags().StringVarP(&runConfig.Protocol, "app-protocol", "P", "http", "The protocol (gRPC or HTTP) Dapr uses to talk to the application")
	RunCMD.PersistentFlags().StringVarP(&runConfig.ComponentsPath, "components-path", "d", standalone.DefaultComponentsDirPath(), "The path for components directory")
	RunCMD.PersistentFlags().StringVarP(&runConfig.PlacementHost, "placement-host-address", "", "localhost", "The host on which the placement service resides")
	RunCMD.PersistentFlags().BoolVar(&runConfig.AppSSL, "app-ssl", false, "Enable https when Dapr invokes the application")
	RunCMD.PersistentFlags().IntVarP(&runConfig.MetricsPort, "metrics-port", "M", dapr.DefaultRandomPort, "The port of metrics on dapr")
	RunCMD.PersistentFlags().StringVarP(&runConfig.AppPwd, "pwd", "", "", "The dir to run cmd in")
	RunCMD.PersistentFlags().IntVarP(&runConfig.AppWaitTimeoutInSeconds, "wait-timeout", "", dapr.DefaultAppWaitTimeoutInSeconds, "The timeout in second to wait for app start")
	RootCMD.AddCommand(RunCMD)
}
