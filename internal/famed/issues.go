package famed

import (
	"log"
	"sort"

	"github.com/google/go-github/v41/github"
)

type Issues map[int64]Issues

type Issue struct {
	Issue  *github.Issue
	Events []*github.IssueEvent
	// For issue comment generation
	Error error
}

// issuesToContributors generates a contributor list based on a list of issues
func (r *repo) contributorsArray() []*Contributor {
	// Generate the contributors from the issues and events
	contributors := r.ContributorsForIssues()
	// Transformation of contributors map to contributors array
	contributorsArray := contributors.toSortedSlice()
	// Sort contributors array by total rewards
	sortContributors(contributorsArray)

	return contributorsArray
}

// filterIssues filters for valid issues.
func filterIssues(issues []*github.Issue) []*github.Issue {
	filteredIssues := make([]*github.Issue, 0)
	for _, issue := range issues {
		if _, err := IsIssueValid(issue); err != nil {
			log.Printf("[issuesToContributors] issue invalid with ID: %d, error: %v \n", issue.ID, err)
			continue
		}
		filteredIssues = append(filteredIssues, issue)
	}

	return filteredIssues
}

func (contributors Contributors) toSortedSlice() []*Contributor {
	contributorsSlice := contributors.toSlice()
	sortContributors(contributorsSlice)
	return contributorsSlice
}

// mapToSlice transforms the contributors map to a contributors slice.
func (contributors Contributors) toSlice() []*Contributor {
	contributorsSlice := make([]*Contributor, 0)
	for _, contributor := range contributors {
		contributorsSlice = append(contributorsSlice, contributor)
	}

	return contributorsSlice
}

// sortContributors sorts the contributors by descending reward sum.
func sortContributors(contributors []*Contributor) {
	sort.SliceStable(contributors, func(i, j int) bool {
		return contributors[i].RewardSum > contributors[j].RewardSum
	})
}
