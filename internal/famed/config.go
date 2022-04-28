package famed

import (
	"github.com/morphysm/famed-github-backend/internal/client/github"
)

type Config struct {
	Currency  string
	Rewards   map[github.IssueSeverity]float64
	Labels    map[string]github.Label
	DaysToFix int
	BotLogin  string
}

// NewFamedConfig returns a new instance of the famed config.
func NewFamedConfig(currency string, rewards map[github.IssueSeverity]float64, labels map[string]github.Label, daysToFix int, botLogin string) Config {
	return Config{
		Currency:  currency,
		Rewards:   rewards,
		Labels:    labels,
		DaysToFix: daysToFix,
		BotLogin:  botLogin,
	}
}
