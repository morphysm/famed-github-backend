package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/morphysm/famed-github-backend/internal/client/installation"
)

type Config struct {
	App struct {
		Host string
		Port string
	}

	Currency struct {
		Host string
	}

	Github struct {
		Host          string
		Key           string
		WebhookSecret string
		AppID         int64
		BotID         int64
		BotLogin      string
	}

	Famed struct {
		Labels   map[string]installation.Label
		Rewards  map[IssueSeverity]float64
		Currency string
	}

	Admin struct {
		Username string
		Password string
	}
}

type IssueSeverity string

const (
	githubKeyEnvName      = "GITHUB_API_KEY"
	githubWHSecretEnvName = "GITHUB_WEBHOOK_SECRET" //nolint:gosec
	githubAppIDEnvName    = "GITHUB_APP_ID"
	githubBotID           = "GITHUB_BOT_ID"
	githubBotLogin        = "GITHUB_BOT_LOGIN"
	adminUsername         = "ADMIN_USERNAME"
	adminPassword         = "ADMIN_PASSWORD"

	FamedLabel = "famed"
	// CVSSNone represents a CVSS of 0
	CVSSNone IssueSeverity = "none"
	// CVSSLow represents a CVSS of 0.1-3.9
	CVSSLow IssueSeverity = "low"
	// CVSSMedium represents a CVSS of 4.0-6.9
	CVSSMedium IssueSeverity = "medium"
	// CVSSHigh represents a CVSS of 7.0-8.9
	CVSSHigh IssueSeverity = "high"
	// CVSSCritical represents a CVSS of 9.0-10.0
	CVSSCritical IssueSeverity = "critical"
)

func Load() (*Config, error) {
	config := Config{}

	// Read json config file
	err := bindConfigFile(&config)
	if err != nil {
		return nil, err
	}

	if err := verifyConfig(config); err != nil {
		return nil, err
	}

	// GitHub api key
	if err := bindString(&config.Github.Key, githubKeyEnvName); err != nil {
		return nil, err
	}

	// GitHub api key
	if err := bindString(&config.Github.WebhookSecret, githubWHSecretEnvName); err != nil {
		return nil, err
	}

	// GitHub Famed app id
	if err := bindInt64(&config.Github.AppID, githubAppIDEnvName); err != nil {
		return nil, err
	}

	// GitHub bot id
	if err := bindInt64(&config.Github.BotID, githubBotID); err != nil {
		return nil, err
	}

	// GitHub bot name
	if err := bindString(&config.Github.BotLogin, githubBotLogin); err != nil {
		return nil, err
	}

	// Admin username
	if err := bindString(&config.Admin.Username, adminUsername); err != nil {
		return nil, err
	}

	// Admin password
	if err := bindString(&config.Admin.Password, adminPassword); err != nil {
		return nil, err
	}

	return &config, nil
}

func bindConfigFile(cfg *Config) error {
	configFile, err := os.Open("config.json")
	if err != nil {
		return err
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(cfg)
	if err != nil {
		return err
	}

	return nil
}

func verifyConfig(cfg Config) error {
	if cfg.App.Host == "" {
		return errors.New("config.json app.host must be set")
	}
	if cfg.App.Port == "" {
		return errors.New("config.json app.host must be set")
	}
	if cfg.Currency.Host == "" {
		return errors.New("config.json currency.host must be set")
	}
	if cfg.Github.Host == "" {
		return errors.New("config.json github.host must be set")
	}

	if err := verifyLabel(cfg, FamedLabel); err != nil {
		return err
	}
	if err := verifyLabel(cfg, string(CVSSNone)); err != nil {
		return err
	}
	if err := verifyLabel(cfg, string(CVSSLow)); err != nil {
		return err
	}
	if err := verifyLabel(cfg, string(CVSSMedium)); err != nil {
		return err
	}
	if err := verifyLabel(cfg, string(CVSSHigh)); err != nil {
		return err
	}
	if err := verifyLabel(cfg, string(CVSSCritical)); err != nil {
		return err
	}

	if err := verifyReward(cfg, CVSSNone); err != nil {
		return err
	}
	if err := verifyReward(cfg, CVSSLow); err != nil {
		return err
	}
	if err := verifyReward(cfg, CVSSMedium); err != nil {
		return err
	}
	if err := verifyReward(cfg, CVSSHigh); err != nil {
		return err
	}
	if err := verifyReward(cfg, CVSSCritical); err != nil {
		return err
	}

	return nil
}

func verifyLabel(cfg Config, label string) error {
	if label, ok := cfg.Famed.Labels[label]; !ok || label.Name == "" || label.Color == "" || label.Description == "" {
		return fmt.Errorf("config.json app.famed.labels.%s must be set", label)
	}

	return nil
}

func verifyReward(cfg Config, cvss IssueSeverity) error {
	if _, ok := cfg.Famed.Rewards[cvss]; !ok {
		return fmt.Errorf("config.json app.famed.rewards.%s must be set", cvss)
	}

	return nil
}
