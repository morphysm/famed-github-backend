package famed_test

import (
	"testing"

	"github.com/morphysm/famed-github-backend/internal/client/installation"
	"github.com/morphysm/famed-github-backend/internal/config"
	"github.com/morphysm/famed-github-backend/internal/famed"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	t.Parallel()

	cfg := NewTestConfig()

	assert.Equal(t, cfg.Currency, "c")
	assert.Equal(t, cfg.Rewards, map[config.IssueSeverity]float64{"s": 0})
	assert.Equal(t, cfg.Labels, map[string]installation.Label{
		"f": {
			Name:        "n",
			Color:       "c",
			Description: "d",
		},
	})
	assert.Equal(t, cfg.BotLogin, "b")
}

func NewTestConfig() famed.Config {
	return famed.NewFamedConfig("c",
		map[config.IssueSeverity]float64{"s": 0},
		map[string]installation.Label{
			"f": {
				Name:        "n",
				Color:       "c",
				Description: "d",
			},
		},
		"b")
}
