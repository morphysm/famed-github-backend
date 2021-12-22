package config

type Config struct {
	App struct {
		Host string
		Port string
	}

	Github struct{
		Key string
	}

}

const (
	githubKeyEnvName = "GITHUB_API_KEY"
)

func Load() (*Config, error) {
	var config = Config{}

	config.App.Host = "127.0.0.1"
	config.App.Port = "8080"

	// COMPANY EMAIL
	err := bindString(&config.Github.Key, githubKeyEnvName)
	if err != nil {
		return nil, err
	}

	return &config, nil
}