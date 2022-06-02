package comment_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/morphysm/famed-github-backend/internal/famed/model/comment"
	"github.com/morphysm/famed-github-backend/internal/repositories/github/model"
)

const (
	eligibleCommentLegacy = "🤖 Assignees for Issue **Test 3 #5** are now eligible to Get Famed.\n\n✅ Add assignees to track contribution times of the issue \U0001F9B8‍♀️\U0001F9B9️\n✅ Add a single severity (CVSS) label to compute the score 🏷️️\n✅ Link a PR when closing the issue ♻️ \U0001F9B8‍♀️\U0001F9B9\n\nHappy hacking! \U0001F9BE💙❤️️"
	eligibleCommentV1     = "<!--{\"type\":\"eligible\",\"version\":\"TODO\"}-->\n🤖 Assignees for issue **Test 3 #5** are now eligible to Get Famed.\n\n✅ Add assignees to track contribution times of the issue \U0001F9B8‍♀️\U0001F9B9️\n✅ Add a single severity (CVSS) label to compute the score 🏷️️\n\nHappy hacking! \U0001F9BE💙❤️"

	contributorComment = "This is a contributor comment"
)

func TestFindComment(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name          string
		CommentType   comment.Type
		BotLogin      string
		Comments      []model.IssueComment
		ExpectedFind  model.IssueComment
		ExpectedFound bool
	}{
		{
			Name:          "Single Legacy Eligible Comment",
			CommentType:   comment.EligibleCommentType,
			BotLogin:      "test[bot]",
			Comments:      []model.IssueComment{{User: model.User{Login: "test[bot]"}, Body: eligibleCommentLegacy}},
			ExpectedFind:  model.IssueComment{User: model.User{Login: "test[bot]"}, Body: eligibleCommentLegacy},
			ExpectedFound: true,
		},
		{
			Name:          "Single V1 Eligible Comment",
			CommentType:   comment.EligibleCommentType,
			BotLogin:      "test[bot]",
			Comments:      []model.IssueComment{{User: model.User{Login: "test[bot]"}, Body: eligibleCommentV1}},
			ExpectedFind:  model.IssueComment{User: model.User{Login: "test[bot]"}, Body: eligibleCommentV1},
			ExpectedFound: true,
		},
		{
			Name:        "Multiple Legacy Eligible Comments",
			CommentType: comment.EligibleCommentType,
			BotLogin:    "test[bot]",
			Comments: []model.IssueComment{
				{User: model.User{Login: "test[bot]"}, Body: eligibleCommentLegacy},
				{User: model.User{Login: "test[bot]"}, Body: eligibleCommentLegacy},
			},
			ExpectedFind:  model.IssueComment{User: model.User{Login: "test[bot]"}, Body: eligibleCommentLegacy},
			ExpectedFound: true,
		},
		{
			Name:        "Multiple V1 Eligible Comments",
			CommentType: comment.EligibleCommentType,
			BotLogin:    "test[bot]",
			Comments: []model.IssueComment{
				{User: model.User{Login: "test[bot]"}, Body: eligibleCommentV1},
				{User: model.User{Login: "test[bot]"}, Body: eligibleCommentV1},
			},
			ExpectedFind:  model.IssueComment{User: model.User{Login: "test[bot]"}, Body: eligibleCommentV1},
			ExpectedFound: true,
		},
		{
			Name:        "Legacy and V1 Eligible Comments",
			CommentType: comment.EligibleCommentType,
			BotLogin:    "test[bot]",
			Comments: []model.IssueComment{
				{User: model.User{Login: "test[bot]"}, Body: eligibleCommentLegacy},
				{User: model.User{Login: "test[bot]"}, Body: eligibleCommentV1},
			},
			ExpectedFind:  model.IssueComment{User: model.User{Login: "test[bot]"}, Body: eligibleCommentLegacy},
			ExpectedFound: true,
		},
		{
			Name:        "Contributor and V1 Eligible Comments",
			CommentType: comment.EligibleCommentType,
			BotLogin:    "test[bot]",
			Comments: []model.IssueComment{
				{User: model.User{Login: "contributor"}, Body: contributorComment},
				{User: model.User{Login: "contributor"}, Body: contributorComment},
				{User: model.User{Login: "test[bot]"}, Body: eligibleCommentV1},
				{User: model.User{Login: "contributor"}, Body: contributorComment},
			},
			ExpectedFind:  model.IssueComment{User: model.User{Login: "test[bot]"}, Body: eligibleCommentV1},
			ExpectedFound: true,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			// GIVEN

			// WHEN
			foundComment, found := comment.Comments(testCase.Comments).FindComment(testCase.BotLogin, testCase.CommentType)

			// THEN
			if testCase.ExpectedFound {
				assert.True(t, found)
				assert.Equal(t, testCase.ExpectedFind, foundComment)
				return
			}

			assert.False(t, found)
		})
	}
}
