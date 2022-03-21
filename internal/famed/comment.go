package famed

import (
	"errors"
	"fmt"

	"github.com/google/go-github/v41/github"
	"github.com/morphysm/famed-github-backend/internal/client/installation"
)

var ErrNoContributors = errors.New("GitHub data incomplete")

func RewardComment(contributors Contributors, currency string) string {
	if len(contributors) == 0 {
		return RewardCommentFromError(ErrNoContributors)
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

func RewardCommentFromError(err error) string {
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

// IssueEligibleComment generate an issue eligible RewardComment.
func IssueEligibleComment(issue *github.Issue, pullRequest *installation.PullRequest) (string, error) {
	comment := fmt.Sprintf("🤖 Assignees for WrappedIssue **%s #%d** are now eligible to Get Famed.", *issue.Title, *issue.Number)

	// Check that an assignee is assigned
	comment = fmt.Sprintf("%s\n%s️", comment, assigneeComment(issue))

	// Check that a valid severity label is assigned
	comment = fmt.Sprintf("%s\n%s️", comment, severityComment(WrappedIssue{Issue: issue}))

	// Check that a PR is assigned
	comment = fmt.Sprintf("%s\n%s", comment, prComment(pullRequest))

	// Final note
	comment = fmt.Sprintf("%s\n\nHappy hacking! 🦾💙❤️️", comment)

	return comment, nil
}

func assigneeComment(issue *github.Issue) string {
	if issue.Assignee != nil {
		return "- ✅ Add assignees to track contribution times of the issue \U0001F9B8‍♀️\U0001F9B9"
	}

	return "- ❌ Add assignees to track contribution times of the issue \U0001F9B8‍♀️\U0001F9B9"
}

func severityComment(issue WrappedIssue) string {
	if _, err := issue.severity(); err == nil {
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
