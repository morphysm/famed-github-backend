package famed

import (
	"github.com/morphysm/famed-github-backend/internal/client/installation"
	"github.com/morphysm/famed-github-backend/internal/config"
)

type Config struct {
	Currency string
	Rewards  map[config.IssueSeverity]float64
	Labels   map[string]installation.Label
	BotLogin string
}

// NewFamedConfig returns a new instance of the famed config.
func NewFamedConfig(currency string, rewards map[config.IssueSeverity]float64, labels map[string]installation.Label, botLogin string) Config {
	return Config{
		Currency: currency,
		Rewards:  rewards,
		Labels:   labels,
		BotLogin: botLogin,
	}
}
