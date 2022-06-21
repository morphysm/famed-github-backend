package config

import (
	"github.com/awnumar/memguard"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/morphysm/famed-github-backend/internal/repositories/github/model"
	"github.com/rotisserie/eris"
	"os"
	"strings"
)

//delimiter delimits hierarchy levels in configuration names
const delimiter = "."

//envPrefix prefixes all environment variables
const envPrefix = "FAMED_"

// FamedLabelKey is the label used in GitHub to tell our backend that this issue should be tracked by famed. // Todo: make it configurable
const FamedLabelKey = "famed"

// NewConfig returns a fully initialized(? maybe not the best word) configuration.
// The configuration can be set and loaded from different sources. The following load order is used:
// Defaults values, which can be overridden by
// YAML config from XDG path, which can be overridden by
// dotenv file (./.env file), which can be overridden by
// environment variables
func NewConfig(filePath string) (config *Config, err error) {
	koanf := koanf.New(delimiter)

	// Load defaults values
	if err := koanf.Load(confmap.Provider(defaultConfig, delimiter), nil); err != nil {
		return nil, eris.Wrap(err, "failed to load configuration from default values")
	}

	// if the XDG json config file exists then load the configuration from this file.
	if jsonFile, err := os.Stat(filePath); err == nil {
		if err := koanf.Load(file.Provider(jsonFile.Name()), json.Parser()); err != nil {
			return nil, eris.Wrap(err, "failed to load configuration from json config")
		}
	}

	// Todo: bad loading
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

	err = verifyConfig(*config)
	if err != nil {
		return nil, err
	}

	config.Github.KeyEnclave = memguard.NewEnclave([]byte(config.Github.Key))
	config.Github.Key = ""

	return config, nil
}

func verifyConfig(cfg Config) error {
	if cfg.App.Host == "" {
		return eris.New("config.json app.host must be set")
	}

	if cfg.App.Port == "" {
		return eris.New("config.json app.host must be set")
	}

	if cfg.Github.Host == "" {
		return eris.New("config.json github.host must be set")
	}

	if cfg.Famed.DaysToFix == 0 {
		return eris.New("config.json famed.daysToFix must be set")
	}

	if cfg.Famed.UpdateFrequency == 0 {
		return eris.New("config.json famed.updateFrequency must be set")
	}

	if err := verifyLabel(cfg, FamedLabelKey); err != nil {
		return err
	}

	if err := verifyLabel(cfg, string(model.Info)); err != nil {
		return err
	}

	if err := verifyLabel(cfg, string(model.Low)); err != nil {
		return err
	}

	if err := verifyLabel(cfg, string(model.Medium)); err != nil {
		return err
	}

	if err := verifyLabel(cfg, string(model.High)); err != nil {
		return err
	}

	if err := verifyLabel(cfg, string(model.Critical)); err != nil {
		return err
	}

	if err := verifyReward(cfg, model.Info); err != nil {
		return err
	}

	if err := verifyReward(cfg, model.Low); err != nil {
		return err
	}

	if err := verifyReward(cfg, model.Medium); err != nil {
		return err
	}

	if err := verifyReward(cfg, model.High); err != nil {
		return err
	}

	if err := verifyReward(cfg, model.Critical); err != nil {
		return err
	}

	if cfg.Github.Key == "" {
		return eris.New("missing github key")
	}

	if cfg.Github.WebhookSecret == "" {
		return eris.New("missing github webhook secret")
	}

	if cfg.Github.AppID == 0 {
		return eris.New("missing github appid")
	}

	if cfg.Github.BotLogin == "" {
		return eris.New("missing github botlogin")
	}

	if cfg.Admin.Username == "" {
		return eris.New("missing admin username")
	}

	if cfg.Admin.Password == "" {
		return eris.New("missing admin password")
	}

	return nil
}

func verifyLabel(cfg Config, label string) error {
	if label, ok := cfg.Famed.Labels[label]; !ok || label.Name == "" || label.Color == "" || label.Description == "" {
		return eris.Errorf("config.json app.famed.labels.%s must be set", label)
	}

	return nil
}

func verifyReward(cfg Config, cvss model.IssueSeverity) error {
	if _, ok := cfg.Famed.Rewards[cvss]; !ok {
		return eris.Errorf("config.json app.famed.rewards.%s must be set", cvss)
	}

	return nil
}
