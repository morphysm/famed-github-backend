package config

import (
	"github.com/awnumar/memguard"
	"github.com/morphysm/famed-github-backend/internal/repositories/github/model"
	"github.com/phuslu/log"
)

// Config is the complete representation of the configuration, it is authoritative on configuration names, hierarchy, structure and type.
type Config struct {
	App struct {
		Host     string    `koanf:"host"`
		Port     string    `koanf:"port"`
		LogLevel log.Level `koanf:"loglevel"`
	} `koanf:"app"`

	NewRelic struct {
		Name    string `koanf:"name"`
		Key     string `koanf:"key"`
		Enabled bool   `koanf:"enabled"`
	} `koanf:"newrelic"`

	Github struct {
		Host          string            `koanf:"host"`
		KeyEnclave    *memguard.Enclave `koanf:"keyenclave"`
		WebhookSecret string            `koanf:"webhooksecret"`
		AppID         int64             `koanf:"appid"`
		BotLogin      string            `koanf:"botlogin"`
	} `koanf:"github"`

	Famed struct {
		Labels          map[string]model.Label          `koanf:"labels"`
		Rewards         map[model.IssueSeverity]float64 `koanf:"rewards"`
		Currency        string                          `koanf:"currency"`
		DaysToFix       int                             `koanf:"daystofix"`
		UpdateFrequency int                             `koanf:"updatefrequency"`
	} `koanf:"famed"`

	// TODO this should probably not be in memory
	RedTeamLogins map[string]string `koanf:"redteamlogins"`

	Admin struct {
		Username string `koanf:"username"`
		Password string `koanf:"password"`
	} `koanf:"admin"`
}
