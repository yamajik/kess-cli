package dapr

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type Configs struct {
	configuration Configuration
	components    []Component

	Dir                   string
	BinDirname            string
	ComponentsDirname     string
	ConfigurationFilename string
}

func (c *Configs) BinDir() string {
	return filepath.Join(c.Dir, c.BinDirname)
}

func (c *Configs) ComponentsDir() string {
	return filepath.Join(c.Dir, c.ComponentsDirname)
}

func (c *Configs) ConfigurationFile() string {
	return filepath.Join(c.Dir, c.ConfigurationFilename)
}

func (c *Configs) SetConfiguration(configuration Configuration) *Configs {
	c.configuration = configuration
	return c
}

func (c *Configs) SetComponents(components []Component) *Configs {
	c.components = components
	return c
}

func (c *Configs) Save() error {
	configurationBytes, err := yaml.Marshal(c.configuration)
	if err != nil {
		return errors.WithStack(err)
	}
	if err := c.writeFile(c.ConfigurationFile(), configurationBytes); err != nil {
		return errors.WithStack(err)
	}

	for _, component := range c.components {
		componentBytes, err := yaml.Marshal(component)
		if err != nil {
			return errors.WithStack(err)
		}
		if err := c.writeFile(filepath.Join(c.ComponentsDir(), fmt.Sprintf("%s.yaml", component.Metadata.Name)), componentBytes); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func (c *Configs) Buffer() (*bytes.Buffer, error) {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	defer tw.Close()

	configurationBytes, err := yaml.Marshal(c.configuration)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if err := c.writeTarFile(tw, filepath.Join(c.ConfigurationFilename), configurationBytes); err != nil {
		return nil, errors.WithStack(err)
	}

	for _, component := range c.components {
		componentBytes, err := yaml.Marshal(component)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if err := c.writeTarFile(tw, filepath.Join(c.ComponentsDirname, fmt.Sprintf("%s.yaml", component.Metadata.Name)), componentBytes); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return &buf, nil
}

func (c *Configs) writeTarFile(tw *tar.Writer, filename string, bytes []byte) error {
	if err := tw.WriteHeader(&tar.Header{
		Name: filename,
		Mode: 0644,
		Size: int64(len(bytes)),
	}); err != nil {
		return errors.WithStack(err)
	}
	if _, err := tw.Write(bytes); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (c *Configs) writeFile(filename string, bytes []byte) error {
	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return errors.WithStack(err)
	}
	if err := ioutil.WriteFile(filename, bytes, 0644); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func DefaultConfigs() *Configs {
	return &Configs{
		Dir:                   DefaultDaprDirPath(),
		BinDirname:            DefaultDaprBinDirname,
		ComponentsDirname:     DefaultDaprComponentsDirname,
		ConfigurationFilename: DefaultDaprConfigurationFilename,
	}
}
