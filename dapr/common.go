package dapr

import (
	"os"
	"path/filepath"
)

const (
	DefaultDaprDirname               = ".dapr"
	DefaultDaprBinDirname            = "bin"
	DefaultDaprPlacementFilename     = "placement"
	DefaultDaprDaprdFilename         = "daprd"
	DefaultDaprDashboardDirname      = "dashboard"
	DefaultDaprComponentsDirname     = "components"
	DefaultDaprConfigurationFilename = "config.yaml"
	DefaultAppWaitTimeoutInSeconds   = 60
	DefaultRandomPort                = -1
)

func DefaultDaprDirPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, DefaultDaprDirname)
}

func DefaultDaprBinPath() string {
	return filepath.Join(DefaultDaprDirPath(), DefaultDaprBinDirname)
}

func DefaultDaprPlacementPath() string {
	return filepath.Join(DefaultDaprBinPath(), DefaultDaprPlacementFilename)
}

func DefaultDaprDaprdPath() string {
	return filepath.Join(DefaultDaprBinPath(), DefaultDaprDaprdFilename)
}

func DefaultDaprDashboardPath() string {
	return filepath.Join(DefaultDaprBinPath(), DefaultDaprDashboardDirname)
}

func DefaultComponentsDirPath() string {
	return filepath.Join(DefaultDaprDirPath(), DefaultDaprComponentsDirname)
}

func DefaultConfigFilePath() string {
	return filepath.Join(DefaultDaprDirPath(), DefaultDaprConfigurationFilename)
}
