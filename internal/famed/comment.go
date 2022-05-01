package famed

import (
	"errors"
	"fmt"
	"strings"

	"github.com/morphysm/famed-github-backend/internal/respositories/github/model"
)

var ErrNoContributors = errors.New("GitHub data incomplete")

const (
	rewardCommentErrorHeader       = "### Famed could not generate a reward suggestion."
	rewardCommentTableHeader       = "| Contributor | Time | Reward |\n| ----------- | ----------- | ----------- |"
	eligibleCommentHeaderBeginning = "🤖 Assignees for issue"
)

// rewardComment generates a reward comment.
func rewardComment(contributors contributors, currency string, owner string, repoName string) string {
	if len(contributors) == 0 {
		return rewardCommentFromError(ErrNoContributors)
	}

	sortedContributors := contributors.toSortedSlice()

	comment := ""
	for _, contributor := range sortedContributors {
		comment = fmt.Sprintf("%s@%s ", comment, contributor.Login)
	}

	comment = fmt.Sprintf("%s- you Got Famed! 💎 Check out your new score here: https://www.famed.morphysm.com/teams/%s/%s", comment, owner, repoName)
	comment = fmt.Sprintf("%s\n%s", comment, rewardCommentTableHeader)

	for _, contributor := range sortedContributors {
		comment = fmt.Sprintf("%s\n|%s|%s|%d %s|", comment, contributor.Login, contributor.TotalWorkTime, int(contributor.RewardSum), currency)
	}

	return comment
}

// rewardComment generates a rewar from an error.
func rewardCommentFromError(err error) string {
	comment := rewardCommentErrorHeader +
		"\nReason: "

	switch err {
	case ErrIssueMissingPullRequest:
		return fmt.Sprintf("%sThe issue is missing a pull request.", comment)

	case ErrIssueMissingAssignee:
		return fmt.Sprintf("%sThe issue is missing an assignee.", comment)

	case model.ErrIssueMissingSeverityLabel:
		return fmt.Sprintf("%sThe issue is missing a severity label.", comment)

	case model.ErrIssueMultipleSeverityLabels:
		return fmt.Sprintf("%sThe issue has more than one severity label.", comment)

	case ErrNoContributors:
		comment = fmt.Sprintf("%sThe data provided by GitHub is not sufficient to generate a reward suggestion.", comment)
		return fmt.Sprintf("%s\nThis might be due to an assignment after the issue has been closed. Please assign assignees in the open state.", comment)

	default:
		return fmt.Sprintf("%s Unknown.", comment)
	}
}

// issueEligibleComment generate an issue eligible comment.
func issueEligibleComment(issue model.Issue, pullRequest *string) string {
	comment := fmt.Sprintf("%s **%s #%d** are now eligible to Get Famed.\n", eligibleCommentHeaderBeginning, issue.Title, issue.Number)

	// Check that an assignee is assigned
	comment = fmt.Sprintf("%s\n%s️", comment, assigneeComment(issue.Assignees))

	// Check that a valid severity label is assigned
	comment = fmt.Sprintf("%s\n%s️", comment, severityComment(issue.Severities))

	// Check that a PR is assigned
	// TODO create rule
	// comment = fmt.Sprintf("%s\n%s", comment, prComment(pullRequest))

	// Final note
	comment = fmt.Sprintf("%s\n\nHappy hacking! 🦾💙❤️️", comment)

	return comment
}

func assigneeComment(assignees []model.User) string {
	const msg = " Add assignees to track contribution times of the issue \U0001F9B8\u200d♀️\U0001F9B9"
	if len(assignees) > 0 {
		return "✅" + msg
	}

	return "❌" + msg
}

func severityComment(severities []model.IssueSeverity) string {
	const msg = " Add a single severity (CVSS) label to compute the score 🏷️"
	if len(severities) == 1 {
		return "✅" + msg
	}

	return "❌" + msg
}

// TODO commented out for DevConnect
//func prComment(pullRequest *github.PullRequest) string {
//	const msg = " Link a PR when closing the issue ♻️ \U0001F9B8‍♀️\U0001F9B9"
//	if pullRequest != nil {
//		return "✅" + msg
//	}
//
//	return "❌" + msg
//}

// findComment finds the last of with the commentType and posted by the user with a login equal to botLogin
func findComment(comments []model.IssueComment, botLogin string, commentType commentType) (model.IssueComment, bool) {
	for _, comment := range comments {
		if comment.User.Login == botLogin &&
			verifyCommentType(comment.Body, commentType) {
			return comment, true
		}
	}

	return model.IssueComment{}, false
}

// verifyCommentType checks if a given string is of a given commentType
func verifyCommentType(str string, commentType commentType) bool {
	var substr string
	switch commentType {
	case commentEligible:
		substr = eligibleCommentHeaderBeginning
	case commentReward:
		if strings.Contains(str, rewardCommentTableHeader) {
			return true
		}
		substr = rewardCommentErrorHeader
	}

	return strings.Contains(str, substr)
}
