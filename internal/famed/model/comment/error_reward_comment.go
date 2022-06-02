package comment

import (
	"fmt"
	"strings"

	model2 "github.com/morphysm/famed-github-backend/internal/famed/model"
	"github.com/morphysm/famed-github-backend/internal/repositories/github/model"
)

const ErrorRewardCommentHeader = "### Famed could not generate a reward suggestion."

type ErrorRewardComment struct {
	identifier Identifier
	headline   string
	error      string
}

// NewErrorRewardComment return a ErrorRewardComment.
func NewErrorRewardComment(err error) ErrorRewardComment {
	rewardCommentError := ErrorRewardComment{}
	rewardCommentError.identifier = NewIdentifier(RewardCommentType, "TODO")

	rewardCommentError.headline = ErrorRewardCommentHeader
	rewardCommentError.error = "Reason: "

	switch err {
	case model2.ErrIssueMissingPullRequest:
		rewardCommentError.error = fmt.Sprintf("%sThe issue is missing a pull request.", rewardCommentError.error)

	case model2.ErrIssueMissingAssignee:
		rewardCommentError.error = fmt.Sprintf("%sThe issue is missing an assignee.", rewardCommentError.error)

	case model.ErrIssueMissingSeverityLabel:
		rewardCommentError.error = fmt.Sprintf("%sThe issue is missing a severity label.", rewardCommentError.error)

	case model.ErrIssueMultipleSeverityLabels:
		rewardCommentError.error = fmt.Sprintf("%sThe issue has more than one severity label.", rewardCommentError.error)

	case ErrNoContributors:
		rewardCommentError.error = fmt.Sprintf("%sThe data provided by GitHub is not sufficient to generate a reward suggestion.", rewardCommentError.error)
		rewardCommentError.error = fmt.Sprintf("%s\nThis might be due to an assignment after the issue has been closed. Please assign assignees in the open state.", rewardCommentError.error)

	default:
		rewardCommentError.error = fmt.Sprintf("%s Unknown.", rewardCommentError.error)
	}

	return rewardCommentError
}

func (c ErrorRewardComment) String() (string, error) {
	var sb strings.Builder

	identifier, err := c.identifier.String()
	if err != nil {
		return "", err
	}

	sb.WriteString(identifier)
	sb.WriteString("\n")
	sb.WriteString(c.headline)
	sb.WriteString("\n")
	sb.WriteString(c.error)

	return sb.String(), nil
}

func (c ErrorRewardComment) Type() Type {
	return c.identifier.Type
}
