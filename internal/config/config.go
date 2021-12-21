package config

type Config struct {
	App struct {
		Host string
		Port string
	}


}

const (
	githubTokenEnvName        = "GITHUB_TOKEN"
)

func Load() (*Config, error) {
	var config = Config{}

	config.App.Host = "127.0.0.1"
	config.App.Port = "8080"

	return &config, nil
}