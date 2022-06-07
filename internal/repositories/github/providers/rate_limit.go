package providers

import (
	"context"

	"github.com/morphysm/famed-github-backend/internal/repositories/github/model"
)

func (c *githubInstallationClient) GetRateLimit(ctx context.Context, owner string) (model.RateLimit, error) {
	client, _ := c.clients.get(owner)

	rateLimit, _, err := client.RateLimits(ctx)
	if err != nil {
		return model.RateLimit{}, err
	}

	compressedRateLimit, err := model.NewRateLimit(rateLimit)
	if err != nil {
		return model.RateLimit{}, err
	}

	return compressedRateLimit, nil
}
