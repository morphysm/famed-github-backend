package config

import (
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

	// GitHub api key
	if cfg.Github.KeyEnclave == "" {
		return eris.New("Missing GitHub Key")
	}

	// GitHub api key
	if cfg.Github.WebhookSecret == "" {
		return eris.New("Missing GitHub Webhook secret")
	}

	// GitHub Famed app id
	if cfg.Github.AppID == 0 {
		return eris.New("Missing GitHub AppID")
	}

	// GitHub bot name
	if cfg.Github.BotLogin == "" {
		return eris.New("Missing GitHub BotLogin")
	}

	// Admin username
	if cfg.Admin.Username == "" {
		return eris.New("Missing Admin username")
	}

	// Admin password
	if cfg.Admin.Password == "" {
		return eris.New("Missing Admin password")
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
