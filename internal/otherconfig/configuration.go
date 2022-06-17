package otherconfig

import (
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/rotisserie/eris"
	"os"
	"strings"
)

// delimiter ..
const delimiter = "."

// envPrefix ..
const envPrefix = "FAMED_"

func NewConfig(filePath string) (*Config, error) {
	k := koanf.New(delimiter)

	// Load defaults values.
	err := k.Load(confmap.Provider(defaultConfig, "."), nil)
	if err != nil {
		return nil, eris.Wrap(err, "failed to load default values")
	}

	if yamlFile, err := os.Stat(filePath); err == nil {
		// Load YAML config from filePath
		if err := k.Load(file.Provider(yamlFile.Name()), yaml.Parser()); err != nil {
			return nil, eris.Wrap(err, "failed to load yaml config")
		}
	}

	if envFile, err := os.Stat(".env"); err == nil {
		// Load from ./.env file
		if err := k.Load(file.Provider(envFile.Name()), dotenv.Parser()); err != nil {
			return nil, eris.Wrap(err, "failed to load dotenv config")
		}
	}

	// Load from environment variables
	err = k.Load(env.Provider(envPrefix, delimiter, func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, envPrefix)), "_", ".", -1)
	}), nil)
	if err != nil {
		return nil, eris.Wrap(err, "failed to load environment variables")
	}

	var config *Config

	err = k.Unmarshal("", &config)
	if err != nil {
		return nil, eris.Wrap(err, "failed to unmarshal config")
	}

	return config, nil
}
