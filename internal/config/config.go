package config

type Config struct {
	App struct {
		Host string
		Port string
	}

	Github struct {
		Key           string
		WebhookSecret string
		KudoLabel     string
	}
}

const (
	githubKeyEnvName       = "GITHUB_API_KEY"
	githubWHSecretEnvName  = "GITHUB_WEBHOOK_SECRET" //nolint:gosec
	githubKudoLabelEnvName = "GITHUB_KUDO_LABEL"
)

func Load() (*Config, error) {
	config := Config{}

	config.App.Host = "127.0.0.1"
	config.App.Port = "8080"

	// GitHub api key
	err := bindString(&config.Github.Key, githubKeyEnvName)
	if err != nil {
		return nil, err
	}

	// GitHub api key
	err = bindString(&config.Github.WebhookSecret, githubWHSecretEnvName)
	if err != nil {
		return nil, err
	}

	// GitHub Kudo issue label
	err = bindString(&config.Github.KudoLabel, githubKudoLabelEnvName)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
