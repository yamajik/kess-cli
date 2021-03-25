package dapr

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dapr/cli/pkg/metadata"
	"github.com/dapr/cli/pkg/print"
	"github.com/dapr/cli/pkg/standalone"
	"github.com/dapr/cli/utils"
)

func StandaloneInstall(runtimeVersion string, dashboardVersion string) error {
	if err := standalone.Init(runtimeVersion, dashboardVersion, "", true); err != nil {
		return err
	}
	return nil
}

func StandaloneUninstall() error {
	if err := standalone.Uninstall(false, ""); err != nil {
		return err
	}
	return nil
}

type StandaloneRunConfig struct {
	standalone.RunConfig
	AppWaitTimeoutInSeconds int
	AppPwd                  string
}

func (c *StandaloneRunConfig) Default() error {
	if c.AppWaitTimeoutInSeconds == 0 {
		c.AppWaitTimeoutInSeconds = 60
	}
	return nil
}

func StandaloneRun(config *StandaloneRunConfig) {
	if err := config.Default(); err != nil {
		print.FailureStatusEvent(os.Stdout, err.Error())
		return
	}

	output, err := standalone.Run(&config.RunConfig)
	if err != nil {
		print.FailureStatusEvent(os.Stdout, err.Error())
		return
	}

	sigCh := make(chan os.Signal, 1)
	SetupShutdownNotify(sigCh)

	daprRunning := make(chan bool, 1)
	appRunning := make(chan bool, 1)

	go func() {
		print.InfoStatusEvent(
			os.Stdout,
			fmt.Sprintf(
				"Starting Dapr with id %s. HTTP Port: %v. gRPC Port: %v",
				output.AppID,
				output.DaprHTTPPort,
				output.DaprGRPCPort))

		output.DaprCMD.Stdout = os.Stdout
		output.DaprCMD.Stderr = os.Stderr

		err = output.DaprCMD.Start()
		if err != nil {
			print.FailureStatusEvent(os.Stdout, err.Error())
			os.Exit(1)
		}

		if config.AppPort <= 0 {
			// If app does not listen to port, we can check for Dapr's sidecar health before starting the app.
			// Otherwise, it creates a deadlock.
			sidecarUp := true
			print.InfoStatusEvent(os.Stdout, "Checking if Dapr sidecar is listening on HTTP port %v", output.DaprHTTPPort)
			err = utils.IsDaprListeningOnPort(output.DaprHTTPPort, time.Duration(config.AppWaitTimeoutInSeconds)*time.Second)
			if err != nil {
				sidecarUp = false
				print.WarningStatusEvent(os.Stdout, "Dapr sidecar is not listening on HTTP port: %s", err.Error())
			}

			print.InfoStatusEvent(os.Stdout, "Checking if Dapr sidecar is listening on GRPC port %v", output.DaprGRPCPort)
			err = utils.IsDaprListeningOnPort(output.DaprGRPCPort, time.Duration(config.AppWaitTimeoutInSeconds)*time.Second)
			if err != nil {
				sidecarUp = false
				print.WarningStatusEvent(os.Stdout, "Dapr sidecar is not listening on GRPC port: %s", err.Error())
			}

			if sidecarUp {
				print.InfoStatusEvent(os.Stdout, "Dapr sidecar is up and running.")
			} else {
				print.WarningStatusEvent(os.Stdout, "Dapr sidecar might not be responding.")
			}
		}

		daprRunning <- true
	}()

	<-daprRunning

	go func() {
		if output.AppCMD == nil {
			appRunning <- true
			return
		}

		if config.AppPwd != "" {
			output.AppCMD.Dir = config.AppPwd
		}

		stdErrPipe, pipeErr := output.AppCMD.StderrPipe()
		if pipeErr != nil {
			print.FailureStatusEvent(os.Stdout, fmt.Sprintf("Error creating stderr for App: %s", err.Error()))
			os.Exit(1)
		}

		stdOutPipe, pipeErr := output.AppCMD.StdoutPipe()
		if pipeErr != nil {
			print.FailureStatusEvent(os.Stdout, fmt.Sprintf("Error creating stdout for App: %s", err.Error()))
			os.Exit(1)
		}

		errScanner := bufio.NewScanner(stdErrPipe)
		outScanner := bufio.NewScanner(stdOutPipe)
		go func() {
			for errScanner.Scan() {
				fmt.Println(print.Blue(fmt.Sprintf("== APP == %s\n", errScanner.Text())))
			}
		}()

		go func() {
			for outScanner.Scan() {
				fmt.Println(print.Blue(fmt.Sprintf("== APP == %s\n", outScanner.Text())))
			}
		}()

		err = output.AppCMD.Start()
		if err != nil {
			print.FailureStatusEvent(os.Stdout, err.Error())
			os.Exit(1)
		}

		appRunning <- true
	}()

	<-appRunning

	// Metadata API is only available if app has started listening to port, so wait for app to start before calling metadata API.
	err = metadata.Put(output.DaprHTTPPort, "cliPID", strconv.Itoa(os.Getpid()))
	if err != nil {
		print.WarningStatusEvent(os.Stdout, "Could not update sidecar metadata for cliPID: %s", err.Error())
	}

	if output.AppCMD != nil {
		appCommand := strings.Join(config.Arguments, " ")
		print.InfoStatusEvent(os.Stdout, fmt.Sprintf("Updating metadata for app command: %s", appCommand))
		err = metadata.Put(output.DaprHTTPPort, "appCommand", appCommand)
		if err != nil {
			print.WarningStatusEvent(os.Stdout, "Could not update sidecar metadata for appCommand: %s", err.Error())
		} else {
			print.SuccessStatusEvent(os.Stdout, "You're up and running! Both Dapr and your app logs will appear here.\n")
		}
	} else {
		print.SuccessStatusEvent(os.Stdout, "You're up and running! Dapr logs will appear here.\n")
	}

	<-sigCh
	print.InfoStatusEvent(os.Stdout, "\nterminated signal received: shutting down")

	err = output.DaprCMD.Process.Kill()
	if err != nil {
		print.FailureStatusEvent(os.Stdout, fmt.Sprintf("Error exiting Dapr: %s", err))
	} else {
		print.SuccessStatusEvent(os.Stdout, "Exited Dapr successfully")
	}

	if output.AppCMD != nil {
		err = output.AppCMD.Process.Kill()
		if err != nil {
			print.FailureStatusEvent(os.Stdout, fmt.Sprintf("Error exiting App: %s", err))
		} else {
			print.SuccessStatusEvent(os.Stdout, "Exited App successfully")
		}
	}
}
