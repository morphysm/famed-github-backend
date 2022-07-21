package model

import (
	"time"

	"github.com/google/go-github/v41/github"
)

type RateLimits struct {
	Core   Rate `json:"core"`
	Search Rate `json:"search"`
}

type Rate struct {
	Limit     int       `json:"limit"`
	Remaining int       `json:"remaining"`
	Reset     time.Time `json:"reset"`
}

func NewRateLimit(rateLimits *github.RateLimits) (RateLimits, error) {
	if rateLimits == nil ||
		rateLimits.Core == nil ||
		rateLimits.Search == nil {
		return RateLimits{}, ErrRateLimitMissingData
	}

	return RateLimits{
		Core: Rate{
			Limit:     rateLimits.Core.Limit,
			Remaining: rateLimits.Core.Remaining,
			Reset:     rateLimits.Core.Reset.Time,
		},
		Search: Rate{
			Limit:     rateLimits.Search.Limit,
			Remaining: rateLimits.Search.Remaining,
			Reset:     rateLimits.Search.Reset.Time,
		},
	}, nil
}
