package configuration

import (
	"os"
	"path/filepath"

	"github.com/rotisserie/eris"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
)

const configFilePerm = os.FileMode(0o600)

func (c *Configuration) checkExists() error {
	if _, err := os.Stat(c.configFilePath); err != nil {
		return c.SaveFile()
	}

	return nil
}

func (c *Configuration) LoadFile(configFilePath string) error {
	c.configFilePath = configFilePath

	err := c.checkExists()
	if err != nil {
		return eris.Wrap(err, "unable to check if file already exists")
	}

	if configFileExtension := filepath.Ext(c.configFilePath); configFileExtension != ".yaml" && configFileExtension != ".yml" {
		return eris.New("file is not valid format (.yaml please)")
	}

	if err := c.koanf.Load(file.Provider(c.configFilePath), yaml.Parser()); err != nil {
		return eris.Wrap(err, "unable to load file")
	}
	return nil
}

func (c *Configuration) SaveFile() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	b, err := c.koanf.Marshal(yaml.Parser())
	if err != nil {
		return eris.Wrap(err, "unable to marshal configuration")
	}

	if err := os.WriteFile(c.configFilePath, b, configFilePerm); err != nil {
		return eris.Wrap(err, "unable to write file")
	}

	return nil
}
