package famed_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/morphysm/famed-github-backend/internal/client/github"
	"github.com/morphysm/famed-github-backend/internal/famed"
)

func TestConfig(t *testing.T) {
	t.Parallel()

	cfg := NewTestConfig()

	assert.Equal(t, cfg.Currency, "POINTS")
	assert.Equal(t, cfg.Rewards, map[github.IssueSeverity]float64{
		github.Info:     0,
		github.Low:      1000,
		github.Medium:   2000,
		github.High:     3000,
		github.Critical: 4000,
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
	rewards := map[github.IssueSeverity]float64{
		github.Info:     0,
		github.Low:      1000,
		github.Medium:   2000,
		github.High:     3000,
		github.Critical: 4000,
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
	)
}
