package providers

import (
	"context"

	"github.com/morphysm/famed-github-backend/internal/repositories/github/model"
)

// GetRateLimits returns the rate limits of a user or organization that installed the Famed app.
func (c *githubInstallationClient) GetRateLimits(ctx context.Context, owner string) (model.RateLimits, error) {
	client, _ := c.clients.get(owner)

	rateLimit, _, err := client.RateLimits(ctx)
	if err != nil {
		return model.RateLimits{}, err
	}

	compressedRateLimit, err := model.NewRateLimit(rateLimit)
	if err != nil {
		return model.RateLimits{}, err
	}

	return compressedRateLimit, nil
}
