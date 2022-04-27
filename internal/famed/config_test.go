package famed_test

import (
	"testing"

	"github.com/morphysm/famed-github-backend/internal/client/github"
	"github.com/morphysm/famed-github-backend/internal/config"
	"github.com/morphysm/famed-github-backend/internal/famed"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	t.Parallel()

	cfg := NewTestConfig()

	assert.Equal(t, cfg.Currency, "POINTS")
	assert.Equal(t, cfg.Rewards, map[config.IssueSeverity]float64{
		config.CVSSInfo:     0,
		config.CVSSLow:      1000,
		config.CVSSMedium:   2000,
		config.CVSSHigh:     3000,
		config.CVSSCritical: 4000,
	})
	assert.Equal(t, cfg.Labels, map[string]github.Label{
		"famed": {
			Name:        "famed",
			Color:       "testColor",
			Description: "testDescription",
		},
	})
	assert.Equal(t, cfg.BotLogin, "b")
}

func NewTestConfig() famed.Config {
	rewards := map[config.IssueSeverity]float64{
		config.CVSSInfo:     0,
		config.CVSSLow:      1000,
		config.CVSSMedium:   2000,
		config.CVSSHigh:     3000,
		config.CVSSCritical: 4000,
	}
	labels := map[string]github.Label{
		"famed": {
			Name:        "famed",
			Color:       "testColor",
			Description: "testDescription",
		},
	}

	return famed.NewFamedConfig("POINTS",
		rewards,
		labels,
		40,
		"b",
		map[string]string{"testUser": "testUser"},
	)
}
