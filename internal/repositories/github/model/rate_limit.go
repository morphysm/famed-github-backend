package model

import (
	"time"

	"github.com/google/go-github/v41/github"
)

type RateLimit struct {
	Limit     int       `json:"limit"`
	Remaining int       `json:"remaining"`
	Reset     time.Time `json:"reset"`
}

func NewRateLimit(rateLimits *github.RateLimits) (RateLimit, error) {
	if rateLimits == nil ||
		rateLimits.Core == nil {
		return RateLimit{}, ErrRateLimitMissingData
	}

	reset := rateLimits.Core.Reset.Time

	return RateLimit{
		Limit:     rateLimits.Core.Limit,
		Remaining: rateLimits.Core.Remaining,
		Reset:     reset,
	}, nil
}
