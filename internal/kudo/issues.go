package kudo

import (
	"context"
	"log"
	"sort"

	"github.com/google/go-github/v41/github"
	"github.com/labstack/echo/v4"
)

// issuesToContributors generates a contributor list based on a list of issues
func (bG *boardGenerator) issuesToContributors(ctx context.Context, issues []*github.Issue) ([]*Contributor, error) {
	var (
		contributorsArray = []*Contributor{}
		filteredIssues    []*github.Issue
	)

	if len(issues) == 0 {
		return contributorsArray, nil
	}

	usdToEthRate, err := bG.currencyClient.GetUSDToETHConversion(ctx)
	if err != nil {
		return nil, echo.ErrBadGateway.SetInternal(err)
	}

	for _, issue := range issues {
		if _, err := IsIssueValid(issue); err != nil {
			log.Printf("[issuesToContributors] issue invalid with ID: %d, error: %v \n", issue.ID, err)
			continue
		}
		filteredIssues = append(filteredIssues, issue)
	}

	events, err := bG.getEvents(ctx, filteredIssues, bG.repo)
	if err != nil {
		return nil, echo.ErrBadGateway.SetInternal(err)
	}

	contributors := GenerateContributors(filteredIssues, events, bG.config.Currency, bG.config.Rewards, usdToEthRate)

	// Transformation of contributors map to contributors array
	for _, contributor := range contributors {
		contributorsArray = append(contributorsArray, contributor)
	}

	// Sort contributors array by total rewards
	sort.SliceStable(contributorsArray, func(i, j int) bool {
		return contributorsArray[i].RewardSum > contributorsArray[j].RewardSum
	})

	return contributorsArray, nil
}
