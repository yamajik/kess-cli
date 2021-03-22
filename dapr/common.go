package dapr

import (
	"os"
	"path/filepath"
)

const (
	DefaultDaprDirname               = ".kess"
	DefaultDaprBinDirname            = "bin"
	DefaultDaprComponentsDirname     = "components"
	DefaultDaprConfigurationFilename = "config.yaml"
)

func DefaultDaprDirPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, DefaultDaprDirname)
}

func DefaultDaprBinPath() string {
	return filepath.Join(DefaultDaprDirPath(), DefaultDaprBinDirname)
}

func DefaultComponentsDirPath() string {
	return filepath.Join(DefaultDaprDirPath(), DefaultDaprComponentsDirname)
}

func DefaultConfigFilePath() string {
	return filepath.Join(DefaultDaprDirPath(), DefaultDaprConfigurationFilename)
}
