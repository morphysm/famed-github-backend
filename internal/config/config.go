package config

import (
	"encoding/json"
	"os"

	"github.com/morphysm/kudos-github-backend/internal/kudo"
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
		Host           string
		Key            string
		WebhookSecret  string
		AppID          int64
		Owner          string
		RepoIDs        []int64
		InstallationID int64
	}

	Kudo struct {
		Label    string
		Rewards  map[kudo.IssueSeverity]float64
		Currency string
	}
}

const (
	githubKeyEnvName       = "GITHUB_API_KEY"
	githubWHSecretEnvName  = "GITHUB_WEBHOOK_SECRET" //nolint:gosec
	githubKudoLabelEnvName = "GITHUB_KUDO_LABEL"
	githubAppIDEnvName     = "GITHUB_APP_ID"
	githubOwner            = "GITHUB_OWNER"
	githubInstallationID   = "GITHUB_INSTALLATION_ID"
	githubRepoIDs          = "GITHUB_REPO_IDS"
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

	// GitHub Kudo issue label
	err = bindString(&config.Kudo.Label, githubKudoLabelEnvName)
	if err != nil {
		return nil, err
	}

	// GitHub Kudo app id
	err = bindInt64(&config.Github.AppID, githubAppIDEnvName)
	if err != nil {
		return nil, err
	}

	// GitHub Kudo owner
	err = bindString(&config.Github.Owner, githubOwner)
	if err != nil {
		return nil, err
	}

	// GitHub installation id
	err = bindInt64(&config.Github.InstallationID, githubInstallationID)
	if err != nil {
		return nil, err
	}

	// GitHub repos
	err = bindInt64Slice(&config.Github.RepoIDs, githubRepoIDs)
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
