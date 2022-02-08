package famed

import (
	"errors"
	"fmt"

	"github.com/google/go-github/v41/github"
)

var ErrNoContributors = errors.New("GitHub data incomplete")

func generateComment(issue *github.Issue, events []*github.IssueEvent, boardOptions BoardOptions) string {
	contributors := Contributors{}

	err := contributors.MapIssue(issue, events, boardOptions)
	if err != nil {
		return generateCommentFromError(err)
	}

	return contributors.generateCommentFromContributors(boardOptions.currency)
}

func (contributors Contributors) generateCommentFromContributors(currency string) string {
	if len(contributors) == 0 {
		return generateCommentFromError(ErrNoContributors)
	}

	comment := "### Kudo suggests:\n" +
		"| Contributor | Time | Reward |\n" +
		"| ----------- | ----------- | ----------- |"
	for _, contributor := range contributors {
		comment = fmt.Sprintf("%s\n|%s|%s|%f %s|", comment, contributor.Login, contributor.TotalWorkTime, contributor.RewardSum, currency)
	}

	return comment
}

func generateCommentFromError(err error) string {
	comment := "### Kudo could not generate a reward suggestion. \n" +
		"Reason: "

	if errors.Is(err, ErrIssueMissingAssignee) {
		return fmt.Sprintf("%sThe issue is missing an assignee.", comment)
	}

	if errors.Is(err, ErrIssueMissingSeverityLabel) {
		return fmt.Sprintf("%sThe issue is missing a severity label.", comment)
	}

	if errors.Is(err, ErrIssueMultipleSeverityLabels) {
		return fmt.Sprintf("%sThe issue has more than one severity label.", comment)
	}

	if errors.Is(err, ErrNoContributors) {
		return fmt.Sprintf("%sThe data provided by GitHub is not sufficient to generate a reward suggestion.", comment)
	}

	return fmt.Sprintf("%s Unknown.", comment)
}
