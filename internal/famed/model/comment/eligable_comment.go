package comment

import (
	"fmt"
	"strings"

	"github.com/morphysm/famed-github-backend/internal/repositories/github/model"
)

const (
	EligibleCommentHeaderLegacy    = "🤖 Assignees for Issue"
	EligibleCommentHeaderBeginning = "🤖 Assignees for issue"
)

type EligibleComment struct {
	identifier    Identifier
	headline      string
	assigneeCheck string
	severityCheck string
	footer        string
}

// NewEligibleComment generate an issue eligible comment.
func NewEligibleComment(issue model.Issue, pullRequest *string) EligibleComment {
	eligibleComment := EligibleComment{}
	// TODO add version information.
	eligibleComment.identifier = NewIdentifier(EligibleCommentType, "TODO")

	eligibleComment.headline = fmt.Sprintf("%s **%s #%d** are now eligible to Get Famed.\n", EligibleCommentHeaderBeginning, issue.Title, issue.Number)

	// Check that an assignee is assigned
	eligibleComment.assigneeCheck = fmt.Sprintf("%s️", assigneeCheck(issue.Assignees))

	// Check that a valid severity label is assigned
	eligibleComment.severityCheck = fmt.Sprintf("%s️", severityCheck(issue.Severities))

	// Check that a PR is assigned
	// TODO create rule
	// comment = fmt.Sprintf("%s\n%s", comment, prComment(pullRequest))

	// Final note
	eligibleComment.footer = fmt.Sprintf("Happy hacking! 🦾💙❤️️")

	return eligibleComment
}

func (c EligibleComment) String() (string, error) {
	var sb strings.Builder

	identifier, err := c.identifier.String()
	if err != nil {
		return "", err
	}

	sb.WriteString(identifier)
	sb.WriteString("\n")
	sb.WriteString(c.headline)
	sb.WriteString("\n")
	sb.WriteString(c.assigneeCheck)
	sb.WriteString("\n")
	sb.WriteString(c.severityCheck)
	sb.WriteString("\n\n")
	sb.WriteString(c.footer)

	return sb.String(), nil
}

func (c EligibleComment) Type() Type {
	return c.identifier.Type
}

func assigneeCheck(assignees []model.User) string {
	const msg = " Add assignees to track contribution times of the issue \U0001F9B8\u200d♀️\U0001F9B9"
	if len(assignees) > 0 {
		return "✅" + msg
	}

	return "❌" + msg
}

func severityCheck(severities []model.IssueSeverity) string {
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
