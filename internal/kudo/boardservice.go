package kudo

import (
	"context"

	"github.com/morphysm/kudos-github-backend/internal/client/currency"
	"github.com/morphysm/kudos-github-backend/internal/client/installation"
)

type BoardGenerator interface {
	GetContributors(ctx context.Context) ([]*Contributor, error)
}

type boardGenerator struct {
	config             Config
	installationClient installation.Client
	currencyClient     currency.Client
	repo               string
}

// NewBoardGenerator returns a new instance of the comment generator.
func NewBoardGenerator(config Config, installationClient installation.Client, currencyClient currency.Client, repo string) BoardGenerator {
	return &boardGenerator{
		config:             config,
		installationClient: installationClient,
		currencyClient:     currencyClient,
		repo:               repo,
	}
}

func (bG *boardGenerator) GetContributors(ctx context.Context) ([]*Contributor, error) {
	// Get all issues in repo
	issuesResponse, err := bG.installationClient.GetIssuesByRepo(ctx, bG.repo, []string{bG.config.Label}, installation.Closed)
	if err != nil {
		return nil, err
	}

	// Use issues to generate contributor list
	contributors, err := bG.issuesToContributors(ctx, issuesResponse)
	if err != nil {
		return nil, err
	}

	return contributors, nil
}
