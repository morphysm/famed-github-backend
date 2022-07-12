package famed_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/go-github/v41/github"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/morphysm/famed-github-backend/internal/famed"
	"github.com/morphysm/famed-github-backend/internal/repositories/github/model"
	"github.com/morphysm/famed-github-backend/internal/repositories/github/providers/providersfakes"
	"github.com/morphysm/famed-github-backend/pkg/pointer"
)

const (
	eligibleCommentV1 = "<!--{\"type\":\"eligible\",\"version\":\"TODO\"}-->\n" +
		"🤖 Assignees for issue **TestIssue #0** are now eligible to Get Famed.\n\n" +
		"✅ Add assignees to track contribution times of the issue 🦸‍♀️🦹️\n" +
		"✅ Add a single severity (CVSS) label to compute the score 🏷️️\n\n" +
		"Happy hacking! 🦾💙❤️️"
	rewardCommentV1 = "<!--{\"type\":\"reward\",\"version\":\"TODO\"}-->\n@testUser - you Got Famed! 💎 Check out your new score here: https://www.famed.morphysm.com/teams/testOwner/testRepo\n| Contributor | Time | Reward |\n| ----------- | ----------- | ----------- |\n|testUser|24h0m0s|975 POINTS|"
)

func TestGetUpdateComment(t *testing.T) {
	t.Parallel()

	open := time.Date(2022, 4, 4, 0, 0, 0, 0, time.UTC)
	assigned := open
	closed := open.Add(24 * time.Hour)
	famedConfig := NewTestConfig()
	botUser := famedConfig.BotLogin
	owner := "testOwner"
	repoName := "testRepo"

	testCases := []struct {
		Name                               string
		Issues                             map[int]model.EnrichedIssue
		Events                             []model.IssueEvent
		Comments                           []model.IssueComment
		PullRequest                        *string
		ExpectedGetEnrichedIssuesCallCount int
		ExpectedGetCommentsCallCount       int
		ExpectedPostCommentCallCount       int
		ExpectedUpdateCommentCallCount     int
		ExpectedDeleteCommentCallCount     int
		ExpectedComments                   []string
		ExpectedResponse                   string
		ExpectedErr                        *echo.HTTPError
	}{
		{
			Name: "No Update",
			Issues: map[int]model.EnrichedIssue{0: {
				Issue: model.Issue{
					ID:         0,
					Number:     0,
					HTMLURL:    "TestURL",
					Title:      "TestIssue",
					CreatedAt:  open,
					ClosedAt:   &closed,
					Assignees:  []model.User{{Login: "testUser"}},
					Severities: []model.IssueSeverity{model.IssueSeverity("low")},
					Migrated:   false,
				},
				PullRequest: nil,
				Events:      nil,
			}},
			Events: []model.IssueEvent{
				{
					Event:     "assigned",
					CreatedAt: assigned,
					Assignee:  &model.User{Login: "testUser"},
				},
				{
					Event:     "closed",
					CreatedAt: closed,
					Assignee:  &model.User{Login: "testUser"},
				},
			},
			Comments:                           []model.IssueComment{{ID: 1, User: model.User{Login: botUser}, Body: eligibleCommentV1}, {ID: 2, User: model.User{Login: botUser}, Body: rewardCommentV1}},
			PullRequest:                        pointer.String("test"),
			ExpectedGetEnrichedIssuesCallCount: 1,
			ExpectedGetCommentsCallCount:       1,
			ExpectedPostCommentCallCount:       0,
			ExpectedUpdateCommentCallCount:     0,
			ExpectedDeleteCommentCallCount:     0,
			ExpectedComments:                   []string{},
			ExpectedResponse:                   "{\"updates\":{}}\n",
		},
		{
			Name: "Update Eligible",
			Issues: map[int]model.EnrichedIssue{0: {
				Issue: model.Issue{
					ID:         0,
					Number:     0,
					HTMLURL:    "TestURL",
					Title:      "TestIssue",
					CreatedAt:  open,
					ClosedAt:   &closed,
					Assignees:  []model.User{{Login: "testUser"}},
					Severities: []model.IssueSeverity{model.IssueSeverity("low")},
					Migrated:   false,
				},
				PullRequest: nil,
				Events:      nil,
			}},
			Events: []model.IssueEvent{
				{
					Event:     "assigned",
					CreatedAt: assigned,
					Assignee:  &model.User{Login: "testUser"},
				},
				{
					Event:     "closed",
					CreatedAt: closed,
					Assignee:  &model.User{Login: "testUser"},
				},
			},
			Comments:                           []model.IssueComment{{ID: 1, User: model.User{Login: botUser}, Body: eligibleCommentV1}, {ID: 2, User: model.User{Login: botUser}, Body: rewardCommentV1}},
			PullRequest:                        pointer.String("test"),
			ExpectedGetEnrichedIssuesCallCount: 1,
			ExpectedGetCommentsCallCount:       1,
			ExpectedPostCommentCallCount:       0,
			ExpectedUpdateCommentCallCount:     1,
			ExpectedDeleteCommentCallCount:     0,
			ExpectedComments:                   []string{},
			ExpectedResponse:                   "{\"updates\":{\"1\":{\"eligibleComment\":{\"actions\":[\"update\"],\"errors\":[]},\"rewardComment\":{\"actions\":[],\"errors\":[]}}}}\n",
		},
		{
			Name: "Update Reward",
			Issues: map[int]model.EnrichedIssue{0: {
				Issue: model.Issue{
					ID:         0,
					Number:     0,
					HTMLURL:    "TestURL",
					Title:      "TestIssue",
					CreatedAt:  open,
					ClosedAt:   &closed,
					Assignees:  []model.User{{Login: "testUser"}},
					Severities: []model.IssueSeverity{model.IssueSeverity("low")},
					Migrated:   false,
				},
				PullRequest: nil,
				Events:      nil,
			}},
			Events: []model.IssueEvent{
				{
					Event:     "assigned",
					CreatedAt: assigned,
					Assignee:  &model.User{Login: "testUser2"},
				},
				{
					Event:     "closed",
					CreatedAt: closed,
					Assignee:  &model.User{Login: "testUser"},
				},
			},
			Comments:                           []model.IssueComment{{ID: 1, User: model.User{Login: botUser}, Body: eligibleCommentV1}, {ID: 2, User: model.User{Login: botUser}, Body: rewardCommentV1}},
			PullRequest:                        pointer.String("test"),
			ExpectedGetEnrichedIssuesCallCount: 1,
			ExpectedGetCommentsCallCount:       1,
			ExpectedPostCommentCallCount:       0,
			ExpectedUpdateCommentCallCount:     1,
			ExpectedDeleteCommentCallCount:     0,
			ExpectedComments:                   []string{},
			ExpectedResponse:                   "{\"updates\":{\"0\":{\"eligibleComment\":{\"actions\":[],\"errors\":[]},\"rewardComment\":{\"actions\":[\"update\"],\"errors\":[]}}}}\n",
		},
		{
			Name: "Post Reward",
			Issues: map[int]model.EnrichedIssue{0: {
				Issue: model.Issue{
					ID:         0,
					Number:     0,
					HTMLURL:    "TestURL",
					Title:      "TestIssue",
					CreatedAt:  open,
					ClosedAt:   &closed,
					Assignees:  []model.User{{Login: "testUser"}},
					Severities: []model.IssueSeverity{model.IssueSeverity("low")},
					Migrated:   false,
				},
				PullRequest: nil,
				Events:      nil,
			}},
			Events: []model.IssueEvent{
				{
					Event:     "assigned",
					CreatedAt: assigned,
					Assignee:  &model.User{Login: "testUser2"},
				},
				{
					Event:     "closed",
					CreatedAt: closed,
					Assignee:  &model.User{Login: "testUser"},
				},
			},
			Comments:                           []model.IssueComment{{ID: 1, User: model.User{Login: botUser}, Body: eligibleCommentV1}},
			PullRequest:                        pointer.String("test"),
			ExpectedGetEnrichedIssuesCallCount: 1,
			ExpectedGetCommentsCallCount:       1,
			ExpectedPostCommentCallCount:       1,
			ExpectedUpdateCommentCallCount:     0,
			ExpectedDeleteCommentCallCount:     0,
			ExpectedComments:                   []string{},
			ExpectedResponse:                   "{\"updates\":{\"0\":{\"eligibleComment\":{\"actions\":[],\"errors\":[]},\"rewardComment\":{\"actions\":[\"update\"],\"errors\":[]}}}}\n",
		},
		{
			Name: "Post Eligible - Rotate Reward",
			Issues: map[int]model.EnrichedIssue{0: {
				Issue: model.Issue{
					ID:         0,
					Number:     0,
					HTMLURL:    "TestURL",
					Title:      "TestIssue",
					CreatedAt:  open,
					ClosedAt:   &closed,
					Assignees:  []model.User{{Login: "testUser"}},
					Severities: []model.IssueSeverity{model.IssueSeverity("low")},
					Migrated:   false,
				},
				PullRequest: nil,
				Events:      nil,
			}},
			Events: []model.IssueEvent{
				{
					Event:     "assigned",
					CreatedAt: assigned,
					Assignee:  &model.User{Login: "testUser2"},
				},
				{
					Event:     "closed",
					CreatedAt: closed,
					Assignee:  &model.User{Login: "testUser"},
				},
			},
			Comments:                           []model.IssueComment{{ID: 2, User: model.User{Login: botUser}, Body: rewardCommentV1}},
			PullRequest:                        pointer.String("test"),
			ExpectedGetEnrichedIssuesCallCount: 1,
			ExpectedGetCommentsCallCount:       1,
			ExpectedPostCommentCallCount:       1,
			ExpectedUpdateCommentCallCount:     1,
			ExpectedDeleteCommentCallCount:     0,
			ExpectedComments:                   []string{},
			ExpectedResponse:                   "{\"updates\":{\"0\":{\"eligibleComment\":{\"actions\":[\"update\"],\"errors\":[]},\"rewardComment\":{\"actions\":[\"update\"],\"errors\":[]}}}}\n",
		},
		{
			Name: "Rotate",
			Issues: map[int]model.EnrichedIssue{0: {
				Issue: model.Issue{
					ID:         0,
					Number:     0,
					HTMLURL:    "TestURL",
					Title:      "TestIssue",
					CreatedAt:  open,
					ClosedAt:   &closed,
					Assignees:  []model.User{{Login: "testUser"}},
					Severities: []model.IssueSeverity{model.IssueSeverity("low")},
					Migrated:   false,
				},
				PullRequest: nil,
				Events:      nil,
			}},
			Events: []model.IssueEvent{
				{
					Event:     "assigned",
					CreatedAt: assigned,
					Assignee:  &model.User{Login: "testUser"},
				},
				{
					Event:     "closed",
					CreatedAt: closed,
					Assignee:  &model.User{Login: "testUser"},
				},
			},
			Comments:                           []model.IssueComment{{ID: 1, User: model.User{Login: botUser}, Body: rewardCommentV1}, {ID: 2, User: model.User{Login: botUser}, Body: eligibleCommentV1}},
			PullRequest:                        pointer.String("test"),
			ExpectedGetEnrichedIssuesCallCount: 1,
			ExpectedGetCommentsCallCount:       1,
			ExpectedPostCommentCallCount:       0,
			ExpectedUpdateCommentCallCount:     2,
			ExpectedDeleteCommentCallCount:     0,
			ExpectedComments:                   []string{},
			ExpectedResponse:                   "{\"updates\":{\"0\":{\"eligibleComment\":{\"actions\":[\"order\"],\"errors\":[]},\"rewardComment\":{\"actions\":[\"order\"],\"errors\":[]}}}}\n",
		},
		{
			Name: "Delete Eligible",
			Issues: map[int]model.EnrichedIssue{0: {
				Issue: model.Issue{
					ID:         0,
					Number:     0,
					HTMLURL:    "TestURL",
					Title:      "TestIssue",
					CreatedAt:  open,
					ClosedAt:   &closed,
					Assignees:  []model.User{{Login: "testUser"}},
					Severities: []model.IssueSeverity{model.IssueSeverity("low")},
					Migrated:   false,
				}}},
			Events: []model.IssueEvent{
				{
					Event:     "assigned",
					CreatedAt: assigned,
					Assignee:  &model.User{Login: "testUser"},
				},
				{
					Event:     "closed",
					CreatedAt: closed,
					Assignee:  &model.User{Login: "testUser"},
				},
			},
			Comments:                           []model.IssueComment{{ID: 1, User: model.User{Login: botUser}, Body: eligibleCommentV1}, {ID: 2, User: model.User{Login: botUser}, Body: rewardCommentV1}, {ID: 3, User: model.User{Login: botUser}, Body: eligibleCommentV1}},
			PullRequest:                        pointer.String("test"),
			ExpectedGetEnrichedIssuesCallCount: 1,
			ExpectedGetCommentsCallCount:       1,
			ExpectedPostCommentCallCount:       0,
			ExpectedUpdateCommentCallCount:     0,
			ExpectedDeleteCommentCallCount:     2,
			ExpectedComments:                   []string{},
			ExpectedResponse:                   "{\"updates\":{\"0\":{\"eligibleComment\":{\"actions\":[\"delete\",\"delete\",\"order\"],\"errors\":[]},\"rewardComment\":{\"actions\":[\"order\"],\"errors\":[]}}}}\n",
		},
		{
			Name: "Delete Reward",
			Issues: map[int]model.EnrichedIssue{0: {
				Issue: model.Issue{
					ID:         0,
					Number:     0,
					HTMLURL:    "TestURL",
					Title:      "TestIssue",
					CreatedAt:  open,
					ClosedAt:   &closed,
					Assignees:  []model.User{{Login: "testUser"}},
					Severities: []model.IssueSeverity{model.IssueSeverity("low")},
					Migrated:   false,
				},
				PullRequest: nil,
				Events:      nil,
			}},
			Events: []model.IssueEvent{
				{
					Event:     "assigned",
					CreatedAt: assigned,
					Assignee:  &model.User{Login: "testUser"},
				},
				{
					Event:     "closed",
					CreatedAt: closed,
					Assignee:  &model.User{Login: "testUser"},
				},
			},
			Comments:                           []model.IssueComment{{ID: 1, User: model.User{Login: botUser}, Body: eligibleCommentV1}, {ID: 2, User: model.User{Login: botUser}, Body: rewardCommentV1}, {ID: 3, User: model.User{Login: botUser}, Body: rewardCommentV1}},
			PullRequest:                        pointer.String("test"),
			ExpectedGetEnrichedIssuesCallCount: 1,
			ExpectedGetCommentsCallCount:       1,
			ExpectedPostCommentCallCount:       0,
			ExpectedUpdateCommentCallCount:     0,
			ExpectedDeleteCommentCallCount:     1,
			ExpectedComments:                   []string{},
			ExpectedResponse:                   "{\"updates\":{\"0\":{\"eligibleComment\":{\"actions\":[],\"errors\":[]},\"rewardComment\":{\"actions\":[\"delete\"],\"errors\":[]}}}}\n",
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			// GIVEN
			e := echo.New()

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/github/repos/%s/%s/update", owner, repoName), nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			req.Header.Set(github.EventTypeHeader, "issues")
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.SetParamNames([]string{"owner", "repo_name"}...)
			ctx.SetParamValues([]string{owner, repoName}...)

			fakeInstallationClient := &providersfakes.FakeInstallationClient{}
			fakeInstallationClient.CheckInstallationReturns(true)
			fakeInstallationClient.GetEnrichedIssuesReturns(testCase.Issues, nil)
			fakeInstallationClient.EnrichIssuesStub = func(ctx context.Context, owner string, repoName string, issues []model.Issue) map[int]model.EnrichedIssue {
				enrichedIssues := make(map[int]model.EnrichedIssue, len(issues))
				for _, issue := range issues {
					enrichedIssues[issue.Number] = model.NewEnrichIssue(issue, testCase.PullRequest, testCase.Events)
				}

				return enrichedIssues
			}
			fakeInstallationClient.GetCommentsReturns(testCase.Comments, nil)

			githubHandler := famed.NewHandler(nil, fakeInstallationClient, famedConfig, Now)

			// WHEN
			err := githubHandler.GetUpdateComments(ctx)

			// THEN
			//Get issues
			assert.Equal(t, testCase.ExpectedGetEnrichedIssuesCallCount, fakeInstallationClient.GetEnrichedIssuesCallCount())

			//Get comments
			assert.Equal(t, testCase.ExpectedGetCommentsCallCount, fakeInstallationClient.GetCommentsCallCount())

			// Post comments
			assert.Equal(t, testCase.ExpectedPostCommentCallCount, fakeInstallationClient.PostCommentCallCount())
			if len(testCase.ExpectedComments) > 0 {
				for _, expComment := range testCase.ExpectedComments {
					_, _, _, _, comment := fakeInstallationClient.PostCommentArgsForCall(0)
					assert.Equal(t, expComment, comment)
				}
			}

			// Delete comments
			assert.Equal(t, testCase.ExpectedUpdateCommentCallCount, fakeInstallationClient.UpdateCommentCallCount())

			// Update comments
			assert.Equal(t, testCase.ExpectedDeleteCommentCallCount, fakeInstallationClient.DeleteCommentCallCount())

			// Response
			assert.Equal(t, testCase.ExpectedResponse, rec.Body.String())

			// Error
			if testCase.ExpectedErr != nil {
				echoErr, ok := err.(*echo.HTTPError)
				if ok {
					assert.Equal(t, testCase.ExpectedErr.Code, echoErr.Code)
					assert.Equal(t, testCase.ExpectedErr.Message, echoErr.Message)
				}
			}
		})
	}
}
