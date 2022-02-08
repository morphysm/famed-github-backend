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
	if len(issues) == 0 {
		return []*Contributor{}, nil
	}

	// Get usd to eth conversion rate
	usdToEthRate, err := bG.currencyClient.GetUSDToETHConversion(ctx)
	if err != nil {
		return nil, echo.ErrBadGateway.SetInternal(err)
	}

	// Filter issues for missing data
	filteredIssues := filterIssues(issues, bG.config.Label)

	// Get all events for each issue
	events, err := bG.getEvents(ctx, filteredIssues, bG.repo)
	if err != nil {
		return nil, echo.ErrBadGateway.SetInternal(err)
	}

	// Generate the contributors from the issues and events
	githubData := GithubData{
		issues:        filteredIssues,
		eventsByIssue: events,
	}
	boardOptions := BoardOptions{
		currency:     bG.config.Currency,
		rewards:      bG.config.Rewards,
		usdToEthRate: usdToEthRate,
	}
	contributors := GenerateContributors(githubData, boardOptions)
	// Transformation of contributors map to contributors array
	contributorsArray := mapToSlice(contributors)
	// Sort contributors array by total rewards
	sortContributors(contributorsArray)

	return contributorsArray, nil
}

// filterIssues filters for valid issues.
func filterIssues(issues []*github.Issue, kudoLabel string) []*github.Issue {
	filteredIssues := make([]*github.Issue, 0)
	for _, issue := range issues {
		if _, err := IsIssueValid(issue, kudoLabel); err != nil {
			log.Printf("[issuesToContributors] issue invalid with ID: %d, error: %v \n", issue.ID, err)
			continue
		}
		filteredIssues = append(filteredIssues, issue)
	}

	return filteredIssues
}

// mapToSlice transforms the contributors map to a contributors slice.
func mapToSlice(contributors Contributors) []*Contributor {
	contributorsArray := make([]*Contributor, 0)
	for _, contributor := range contributors {
		contributorsArray = append(contributorsArray, contributor)
	}

	return contributorsArray
}

// sortContributors sorts the contributors by descending reward sum.
func sortContributors(contributors []*Contributor) {
	sort.SliceStable(contributors, func(i, j int) bool {
		return contributors[i].RewardSum > contributors[j].RewardSum
	})
}
