package famed

import (
	"errors"
	"fmt"
)

var ErrNoContributors = errors.New("GitHub data incomplete")

func (r *repo) comment(issueNumber int) string {
	if err := r.issues[issueNumber].Error; err != nil {
		return commentFromError(err)
	}

	if len(r.contributors) == 0 {
		return commentFromError(ErrNoContributors)
	}

	contributors := r.contributors.toSortedSlice()

	comment := "### Famed suggests:\n" +
		"| Contributor | Time | Reward |\n" +
		"| ----------- | ----------- | ----------- |"
	for _, contributor := range contributors {
		comment = fmt.Sprintf("%s\n|%s|%s|%f %s|", comment, contributor.Login, contributor.TotalWorkTime, contributor.RewardSum, r.config.Currency)
	}

	return comment
}

func commentFromError(err error) string {
	comment := "### Famed could not generate a reward suggestion. \n" +
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
