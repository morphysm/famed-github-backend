package config

import (
	"encoding/json"
	"os"

	"github.com/morphysm/famed-github-backend/internal/famed"
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
	}

	Famed struct {
		Label    string
		Rewards  map[famed.IssueSeverity]float64
		Currency string
	}
}

const (
	githubKeyEnvName        = "GITHUB_API_KEY"
	githubWHSecretEnvName   = "GITHUB_WEBHOOK_SECRET" //nolint:gosec
	githubFamedLabelEnvName = "GITHUB_FAMED_LABEL"
	githubAppIDEnvName      = "GITHUB_APP_ID"
	githubBotID             = "GITHUB_BOT_ID"
)

func Load() (*Config, error) {
	config := Config{}

	config.App.Host = "127.0.0.1"
	config.App.Port = "8080"

	config.Github.Host = "https://api.github.com"
	config.Currency.Host = "https://cdn.jsdelivr.net/gh/fawazahmed0/currency-api@1"

	// Read json config file
	err := bindConfigFile(&config)
	if err != nil {
		return nil, err
	}

	// GitHub api key
	err = bindString(&config.Github.Key, githubKeyEnvName)
	if err != nil {
		return nil, err
	}

	// GitHub api key
	err = bindString(&config.Github.WebhookSecret, githubWHSecretEnvName)
	if err != nil {
		return nil, err
	}

	// GitHub Famed issue label
	err = bindString(&config.Famed.Label, githubFamedLabelEnvName)
	if err != nil {
		return nil, err
	}

	// GitHub Famed app id
	err = bindInt64(&config.Github.AppID, githubAppIDEnvName)
	if err != nil {
		return nil, err
	}

	// Github bot id
	err = bindInt64(&config.Github.BotID, githubBotID)
	if err != nil {
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
