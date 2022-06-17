package config

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

// delimiter allows to have the hierarchy between configuration name
const delimiter = "."

// envPrefix is the prefix the environment variables that will be taken into account
const envPrefix = "FAMED_"

// ??
const FamedLabelKey = "famed"

// NewConfig loads from the different available sources the whole configuration
// the function returns the completed structure with all parameters loaded.
// The order of loading is as follows:
// Defaults values, which can be overridden by
// YAML config from XDG path, which can be overridden by
// dotenv file (./.env file), which can be overridden by
// environment variables
func NewConfig(filePath string) (config *Config, err error) {
	koanf := koanf.New(delimiter)

	// Load defaults values.
	if err := koanf.Load(confmap.Provider(defaultConfig, delimiter), nil); err != nil {
		return nil, eris.Wrap(err, "failed to load configuration from default values")
	}

	// if the XDG yaml config file exists then load the configuration from this file.
	if yamlFile, err := os.Stat(filePath); err == nil {
		if err := koanf.Load(file.Provider(yamlFile.Name()), yaml.Parser()); err != nil {
			return nil, eris.Wrap(err, "failed to load configuration from yaml config")
		}
	}

	// if the dotenv (./.env) file exists then load the configuration from this file.
	if envFile, err := os.Stat(".env"); err == nil {
		if err := koanf.Load(file.Provider(envFile.Name()), dotenv.Parser()); err != nil {
			return nil, eris.Wrap(err, "failed to load configuration from .env file")
		}
	}

	// Load from environment variables
	err = koanf.Load(env.Provider(envPrefix, delimiter, func(s string) string {
		return strings.ReplaceAll(strings.ToLower(
			strings.TrimPrefix(s, envPrefix)), "_", delimiter)
	}), nil)
	if err != nil {
		return nil, eris.Wrap(err, "failed to load configuration from environment variables")
	}

	// Try to unmarshal config from all loaders.
	if err = koanf.Unmarshal("", &config); err != nil {
		return nil, eris.Wrap(err, "failed to unmarshal config")
	}

	return config, nil
}
