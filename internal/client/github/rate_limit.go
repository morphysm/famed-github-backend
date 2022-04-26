package github

import (
	"context"
	"time"

	"github.com/google/go-github/v41/github"
)

type RateLimit struct {
	Limit     int       `json:"limit"`
	Remaining int       `json:"remaining"`
	Reset     time.Time `json:"reset"`
}

func (c *githubInstallationClient) GetRateLimit(ctx context.Context, owner string) (RateLimit, error) {
	client, _ := c.clients.get(owner)

	rateLimit, _, err := client.RateLimits(ctx)
	if err != nil {
		return RateLimit{}, err
	}

	compressedRateLimit, err := validateRateLimit(rateLimit)
	if err != nil {
		return RateLimit{}, err
	}

	return compressedRateLimit, nil
}

func validateRateLimit(rateLimits *github.RateLimits) (RateLimit, error) {
	if rateLimits == nil ||
		rateLimits.Core == nil {
		return RateLimit{}, ErrRateLimitMissingData
	}

	time := rateLimits.Core.Reset.Time

	return RateLimit{
		Limit:     rateLimits.Core.Limit,
		Remaining: rateLimits.Core.Remaining,
		Reset:     time,
	}, nil
}
