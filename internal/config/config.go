package config

type Config struct {
	App struct {
		Host string
		Port string
	}

	Github struct {
		Key            string
		WebhookSecret  string
		KudoLabel      string
		AppID          int64
		Owner          string
		RepoIDs        []int64
		InstallationID int64
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

	// Github installation id
	err = bindInt64(&config.Github.InstallationID, githubInstallationID)
	if err != nil {
		return nil, err
	}

	// Github repos
	err = bindInt64Slice(&config.Github.RepoIDs, githubRepoIDs)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
