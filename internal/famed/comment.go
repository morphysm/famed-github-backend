package famed

import (
	"errors"
	"fmt"

	"github.com/morphysm/famed-github-backend/internal/client/installation"
)

var ErrNoContributors = errors.New("GitHub data incomplete")

// rewardComment generates a reward comment.
func rewardComment(contributors Contributors, currency string) string {
	if len(contributors) == 0 {
		return rewardCommentFromError(ErrNoContributors)
	}

	sortedContributors := contributors.toSortedSlice()

	comment := "### Famed suggests:\n" +
		"| Contributor | Time | Reward |\n" +
		"| ----------- | ----------- | ----------- |"
	for _, contributor := range sortedContributors {
		comment = fmt.Sprintf("%s\n|%s|%s|%f %s|", comment, contributor.Login, contributor.TotalWorkTime, contributor.RewardSum, currency)
	}

	return comment
}

// rewardComment generates a reward from an error.
func rewardCommentFromError(err error) string {
	comment := "### Famed could not generate a reward suggestion. \n" +
		"Reason: "

	switch err {
	case ErrIssueMissingPullRequest:
		return fmt.Sprintf("%sThe issue is missing a pull request.", comment)

	case ErrIssueMissingAssignee:
		return fmt.Sprintf("%sThe issue is missing an assignee.", comment)

	case ErrIssueMissingSeverityLabel:
		return fmt.Sprintf("%sThe issue is missing a severity label.", comment)

	case ErrIssueMultipleSeverityLabels:
		return fmt.Sprintf("%sThe issue has more than one severity label.", comment)

	case ErrNoContributors:
		return fmt.Sprintf("%sThe data provided by GitHub is not sufficient to generate a reward suggestion.", comment)

	default:
		return fmt.Sprintf("%s Unknown.", comment)
	}
}

// issueEligibleComment generate an issue eligible comment.
func issueEligibleComment(issue installation.Issue, pullRequest *installation.PullRequest) (string, error) {
	comment := fmt.Sprintf("🤖 Assignees for WrappedIssue **%s #%d** are now eligible to Get Famed.", issue.Title, issue.Number)

	// Check that an assignee is assigned
	comment = fmt.Sprintf("%s\n%s️", comment, assigneeComment(issue))

	// Check that a valid severity label is assigned
	comment = fmt.Sprintf("%s\n%s️", comment, severityComment(issue))

	// Check that a PR is assigned
	comment = fmt.Sprintf("%s\n%s", comment, prComment(pullRequest))

	// Final note
	comment = fmt.Sprintf("%s\n\nHappy hacking! 🦾💙❤️️", comment)

	return comment, nil
}

func assigneeComment(issue installation.Issue) string {
	if issue.Assignee != nil {
		return "- ✅ Add assignees to track contribution times of the issue \U0001F9B8‍♀️\U0001F9B9"
	}

	return "- ❌ Add assignees to track contribution times of the issue \U0001F9B8‍♀️\U0001F9B9"
}

func severityComment(issue installation.Issue) string {
	if _, err := severity(issue); err == nil {
		return "- ✅ Add a severity (CVSS) label to compute the score 🏷️"
	}

	return "- ❌ Add a severity (CVSS) label to compute the score 🏷️"
}

func prComment(pullRequest *installation.PullRequest) string {
	if pullRequest != nil {
		return "- ✅ Link a PR when closing the issue ♻️ \U0001F9B8‍♀️\U0001F9B9"
	}

	return "- ❌ Link a PR when closing the issue ♻️ \U0001F9B8‍♀️\U0001F9B9"
}
