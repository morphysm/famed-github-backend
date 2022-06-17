package configuration

import (
	"sync"

	"github.com/knadh/koanf"
)

type Configuration struct {
	mutex          *sync.Mutex
	koanf          *koanf.Koanf
	configFilePath string
}

func NewConfig() (config *Configuration, err error) {
	config = &Configuration{
		mutex:          &sync.Mutex{},
		koanf:          koanf.New("."),
		configFilePath: "",
	}

	err = config.ResetAllDefaults()
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Configuration) ResetAllDefaults() error {
	for id, value := range EnumsDefault() {
		err := c.Set(Enum(id), value)
		if err != nil {
			return err
		}
	}

	return nil
}
