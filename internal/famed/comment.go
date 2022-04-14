package famed

import (
	"errors"
	"fmt"
	"strings"

	"github.com/morphysm/famed-github-backend/internal/client/github"
)

var ErrNoContributors = errors.New("GitHub data incomplete")

// rewardComment generates a reward comment.
func rewardComment(contributors Contributors, currency string, owner string, repoName string) string {
	if len(contributors) == 0 {
		return rewardCommentFromError(ErrNoContributors)
	}

	sortedContributors := contributors.toSortedSlice()

	comment := ""
	for _, contributor := range sortedContributors {
		comment = fmt.Sprintf("%s@%s ", comment, contributor.Login)
	}

	comment = fmt.Sprintf("%s- you Got Famed! 💎 Check out your new score here: https://www.famed.morphysm.com/teams/%s/%s", comment, owner, repoName)
	comment = fmt.Sprintf("%s\n| Contributor | Time | Reward |\n| ----------- | ----------- | ----------- |", comment)

	for _, contributor := range sortedContributors {
		comment = fmt.Sprintf("%s\n|%s|%s|%d %s|", comment, contributor.Login, contributor.TotalWorkTime, int(contributor.RewardSum), currency)
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
		comment = fmt.Sprintf("%sThe data provided by GitHub is not sufficient to generate a reward suggestion.", comment)
		return fmt.Sprintf("%s\nThis might be due to an assignement after the issue has been closed. Please assigne assignees in the open state.", comment)

	default:
		return fmt.Sprintf("%s Unknown.", comment)
	}
}

// issueEligibleComment generate an issue eligible comment.
func issueEligibleComment(issue github.Issue, pullRequest *github.PullRequest) string {
	comment := fmt.Sprintf("🤖 Assignees for Issue **%s #%d** are now eligible to Get Famed.\n", issue.Title, issue.Number)

	// Check that an assignee is assigned
	comment = fmt.Sprintf("%s\n%s️", comment, assigneeComment(issue.Assignees))

	// Check that a valid severity label is assigned
	comment = fmt.Sprintf("%s\n%s️", comment, severityComment(issue.Labels))

	// Check that a PR is assigned
	// TODO create rule
	// comment = fmt.Sprintf("%s\n%s", comment, prComment(pullRequest))

	// Final note
	comment = fmt.Sprintf("%s\n\nHappy hacking! 🦾💙❤️️", comment)

	return comment
}

func assigneeComment(assignees []github.User) string {
	const msg = " Add assignees to track contribution times of the issue \U0001F9B8\u200d♀️\U0001F9B9"
	if len(assignees) > 0 {
		return "✅" + msg
	}

	return "❌" + msg
}

func severityComment(labels []github.Label) string {
	const msg = " Add a single severity (CVSS) label to compute the score 🏷️"
	if _, err := severity(labels); err == nil {
		return "✅" + msg
	}

	return "❌" + msg
}

func prComment(pullRequest *github.PullRequest) string {
	const msg = " Link a PR when closing the issue ♻️ \U0001F9B8‍♀️\U0001F9B9"
	if pullRequest != nil {
		return "✅" + msg
	}

	return "❌" + msg
}

// findComment finds the last of with the commentType and posted by the user with a login equal to botLogin
func findComment(comments []github.IssueComment, botLogin string, commentType commentType) (github.IssueComment, bool) {
	for _, comment := range comments {
		if comment.User.Login == botLogin &&
			verifyCommentType(comment.Body, commentType) {
			return comment, true
		}
	}

	return github.IssueComment{}, false
}

// verifyCommentType checks if a given string is of a given commentType
func verifyCommentType(str string, commentType commentType) bool {
	var substr string
	switch commentType {
	case commentEligible:
		substr = "are now eligible to Get Famed."
	case commentReward:
		substr = "| Contributor | Time | Reward |"
		if strings.Contains(str, substr) {
			return true
		}
		substr = "Famed could not generate a reward suggestion."
	}

	return strings.Contains(str, substr)
}
