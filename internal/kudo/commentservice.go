package kudo

import (
	"context"
	"log"

	"github.com/google/go-github/v41/github"

	"github.com/morphysm/kudos-github-backend/internal/client/currency"
	"github.com/morphysm/kudos-github-backend/internal/client/installation"
)

type CommentGenerator interface {
	GetComment(ctx context.Context) (string, error)
}

type commentGenerator struct {
	config             Config
	installationClient installation.Client
	currencyClient     currency.Client
	event              *github.IssuesEvent
}

// NewCommentGenerator returns a new instance of the comment generator.
func NewCommentGenerator(config Config, installationClient installation.Client, currencyClient currency.Client, event *github.IssuesEvent) CommentGenerator {
	return &commentGenerator{
		config:             config,
		installationClient: installationClient,
		currencyClient:     currencyClient,
		event:              event,
	}
}

func (cG *commentGenerator) GetComment(ctx context.Context) (string, error) {
	if _, err := IsValidCloseEvent(cG.event, cG.config.Label); err != nil {
		switch err {
		case ErrIssueMissingAssignee:
			return GenerateCommentFromError(err), nil
		default:
			return "", err
		}
	}
	// Get issue events
	events, err := cG.installationClient.GetIssueEvents(ctx, *cG.event.Repo.Name, *cG.event.Issue.Number)
	if err != nil {
		log.Printf("[handleIssueEvent] error getting issue events: %v", err)
		return "", err
	}

	// Get USD to ETH conversion rate
	usdToEthRate, err := cG.currencyClient.GetUSDToETHConversion(ctx)
	if err != nil {
		log.Printf("[handleIssueEvent] error getting usd eth conversion rate: %v", err)
		return "", err
	}

	// Generate comments from issue, events, currency, rewards and conversion rate
	comment := GenerateComment(cG.event.Issue, events, cG.config.Currency, cG.config.Rewards, usdToEthRate)

	return comment, nil
}
