package comment

import (
	"strings"

	"github.com/morphysm/famed-github-backend/internal/respositories/github/model"
)

type Comments []model.IssueComment

// FindComment finds the last of with the commentType and posted by the user with a login equal to botLogin
func (cs Comments) FindComment(botLogin string, commentType Type) (model.IssueComment, bool) {
	for _, comment := range cs {
		if comment.User.Login == botLogin &&
			verifyCommentType(comment.Body, commentType) {
			return comment, true
		}
	}

	return model.IssueComment{}, false
}

// verifyCommentType checks if a given string is of a given commentType
func verifyCommentType(body string, commentType Type) bool {
	var substrs []string
	switch commentType {
	case EligibleCommentType:
		substrs = append(substrs, EligibleCommentHeaderBeginning)
		substrs = append(substrs, EligibleCommentHeaderLegacy)
	case RewardCommentType:
		substrs = append(substrs, RewardCommentTableHeader)
		substrs = append(substrs, ErrorRewardCommentHeader)
	}

	for _, substr := range substrs {
		if strings.Contains(body, substr) {
			return true
		}
	}

	return false
}
