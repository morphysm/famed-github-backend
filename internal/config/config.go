package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/morphysm/famed-github-backend/internal/client/github"
)

type Config struct {
	App struct {
		Host string
		Port string
	}

	NewRelic struct {
		Name    string
		Key     string
		Enabled bool
	}

	Github struct {
		Host          string
		Key           string
		WebhookSecret string
		AppID         int64
		BotLogin      string
	}

	Famed struct {
		Labels          map[string]github.Label
		Rewards         map[github.IssueSeverity]float64
		Currency        string
		DaysToFix       int
		UpdateFrequency int
	}

	// TODO this should probably not be in memory
	RedTeamLogins map[string]string

	Admin struct {
		Username string
		Password string
	}
}

const (
	newRelicNameEnvName    = "NEWRELIC_NAME"
	newRelicKeyEnvName     = "NEWRELIC_KEY"
	newRelicEnabledEnvName = "NEWRELIC_ENABLED"

	githubKeyEnvName      = "GITHUB_API_KEY"
	githubWHSecretEnvName = "GITHUB_WEBHOOK_SECRET" //nolint:gosec
	githubAppIDEnvName    = "GITHUB_APP_ID"
	githubBotLogin        = "GITHUB_BOT_LOGIN"

	adminUsername = "ADMIN_USERNAME"
	adminPassword = "ADMIN_PASSWORD"

	FamedLabelKey = "famed"
)

func Load() (*Config, error) {
	cfg := Config{}

	// Read json cfg file
	if err := bindConfigFile(&cfg); err != nil {
		return nil, err
	}

	// Verify cfg read from file
	if err := verifyConfig(cfg); err != nil {
		return nil, err
	}

	// NewRelic
	if err := loadNewRelic(&cfg); err != nil {
		return nil, err
	}

	// GitHub api key
	if err := bindString(&cfg.Github.Key, githubKeyEnvName); err != nil {
		return nil, err
	}

	// GitHub api key
	if err := bindString(&cfg.Github.WebhookSecret, githubWHSecretEnvName); err != nil {
		return nil, err
	}

	// GitHub Famed app id
	if err := bindInt64(&cfg.Github.AppID, githubAppIDEnvName); err != nil {
		return nil, err
	}

	// GitHub bot name
	if err := bindString(&cfg.Github.BotLogin, githubBotLogin); err != nil {
		return nil, err
	}

	// Admin username
	if err := bindString(&cfg.Admin.Username, adminUsername); err != nil {
		return nil, err
	}

	// Admin password
	if err := bindString(&cfg.Admin.Password, adminPassword); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// loadNewRelic loads the config for NewRelic
func loadNewRelic(cfg *Config) error {
	// NewRelic enabled
	if err := bindBool(&cfg.NewRelic.Enabled, newRelicEnabledEnvName); err != nil {
		log.Printf("%s not found", newRelicEnabledEnvName)
		cfg.NewRelic.Enabled = false
	}
	if cfg.NewRelic.Enabled {
		// NewRelic api key
		if err := bindString(&cfg.NewRelic.Key, newRelicKeyEnvName); err != nil {
			return err
		}
		// NewRelic app name
		if err := bindString(&cfg.NewRelic.Name, newRelicNameEnvName); err != nil {
			return err
		}
	}

	return nil
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
	if cfg.Github.Host == "" {
		return errors.New("config.json github.host must be set")
	}
	if cfg.Famed.DaysToFix == 0 {
		return errors.New("config.json famed.daysToFix must be set")
	}
	if cfg.Famed.UpdateFrequency == 0 {
		return errors.New("config.json famed.updateFrequency must be set")
	}

	if err := verifyLabel(cfg, FamedLabelKey); err != nil {
		return err
	}
	if err := verifyLabel(cfg, string(github.Info)); err != nil {
		return err
	}
	if err := verifyLabel(cfg, string(github.Low)); err != nil {
		return err
	}
	if err := verifyLabel(cfg, string(github.Medium)); err != nil {
		return err
	}
	if err := verifyLabel(cfg, string(github.High)); err != nil {
		return err
	}
	if err := verifyLabel(cfg, string(github.Critical)); err != nil {
		return err
	}

	if err := verifyReward(cfg, github.Info); err != nil {
		return err
	}
	if err := verifyReward(cfg, github.Low); err != nil {
		return err
	}
	if err := verifyReward(cfg, github.Medium); err != nil {
		return err
	}
	if err := verifyReward(cfg, github.High); err != nil {
		return err
	}
	if err := verifyReward(cfg, github.Critical); err != nil {
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

func verifyReward(cfg Config, cvss github.IssueSeverity) error {
	if _, ok := cfg.Famed.Rewards[cvss]; !ok {
		return fmt.Errorf("config.json app.famed.rewards.%s must be set", cvss)
	}

	return nil
}
