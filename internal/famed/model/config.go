package model

import (
	"github.com/morphysm/famed-github-backend/internal/respositories/github/model"
)

type Config struct {
	Currency  string
	Rewards   map[model.IssueSeverity]float64
	Labels    map[string]model.Label
	DaysToFix int
	BotLogin  string
}

// NewFamedConfig returns a new instance of the famed config.
func NewFamedConfig(currency string, rewards map[model.IssueSeverity]float64, labels map[string]model.Label, daysToFix int, botLogin string) Config {
	return Config{
		Currency:  currency,
		Rewards:   rewards,
		Labels:    labels,
		DaysToFix: daysToFix,
		BotLogin:  botLogin,
	}
}
