package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	model2 "github.com/morphysm/famed-github-backend/internal/famed/model"
	"github.com/morphysm/famed-github-backend/internal/repositories/github/model"
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

func NewTestConfig() model2.Config {
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

	return model2.NewFamedConfig("POINTS",
		rewards,
		labels,
		40,
		"b",
	)
}
