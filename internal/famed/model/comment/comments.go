package comment

import (
	"strings"

	"github.com/morphysm/famed-github-backend/internal/repositories/github/model"
)

type Comments []model.IssueComment

// FindComment finds the last of with the commentType and posted by the user with a login equal to botLogin
func (cs Comments) FindComment(botLogin string, commentType Type) (model.IssueComment, bool) {
	for _, comment := range cs {
		if VerifyComment(comment, botLogin, commentType) {
			return comment, true
		}
	}

	return model.IssueComment{}, false
}

// VerifyComment return true if the given comment is of the given comment type and was authored by a user with the given login.
func VerifyComment(comment model.IssueComment, login string, commentType Type) bool {
	if comment.User.Login == login &&
		verifyCommentType(comment.Body, commentType) {
		return true
	}
	return false
}

// verifyCommentType checks if a given string is of a given commentType
// TODO: add detection by matching meta json
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
