package famed_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/morphysm/famed-github-backend/internal/famed"
	"github.com/morphysm/famed-github-backend/internal/respositories/github/model"
)

func TestConfig(t *testing.T) {
	t.Parallel()

	cfg := NewTestConfig()

	assert.Equal(t, cfg.Currency, "POINTS")
	assert.Equal(t, cfg.Rewards, map[model.IssueSeverity]float64{
		model.Info:     0,
		model.Low:      1000,
		model.Medium:   2000,
		model.High:     3000,
		model.Critical: 4000,
	})
	assert.Equal(t, cfg.Labels, map[string]model.Label{
		"famed": {
			Name:        "famed",
			Color:       "testColor",
			Description: "testDescription",
		},
	})
	assert.Equal(t, cfg.BotLogin, "b")
}

func NewTestConfig() famed.Config {
	rewards := map[model.IssueSeverity]float64{
		model.Info:     0,
		model.Low:      1000,
		model.Medium:   2000,
		model.High:     3000,
		model.Critical: 4000,
	}
	labels := map[string]model.Label{
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
	)
}
