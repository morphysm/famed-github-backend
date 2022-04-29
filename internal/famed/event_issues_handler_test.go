package famed_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/go-github/v41/github"
	"github.com/labstack/echo/v4"
	gitlib "github.com/morphysm/famed-github-backend/internal/client/github"
	"github.com/morphysm/famed-github-backend/internal/client/github/githubfakes"
	"github.com/morphysm/famed-github-backend/internal/famed"
	"github.com/morphysm/famed-github-backend/pkg/pointer"
	"github.com/stretchr/testify/assert"
)

func TestPostIssuesEvent(t *testing.T) {
	t.Parallel()

	famedConfig := NewTestConfig()
	testCases := []struct {
		Name            string
		Event           *github.IssuesEvent
		Events          []gitlib.IssueEvent
		PullRequest     *gitlib.PullRequest
		ExpectedComment string
		ExpectedErr     *echo.HTTPError
	}{
		{
			Name: "Close - Empty event",
			Event: &github.IssuesEvent{
				Action: pointer.String("closed"),
			},
			ExpectedComment: "",
			ExpectedErr:     &echo.HTTPError{Code: 400, Message: famed.ErrEventMissingData.Error()},
		},
		{
			Name: "Close - No Assignee",
			Event: &github.IssuesEvent{
				Action: pointer.String("closed"),
				Issue: &github.Issue{
					ID:        pointer.Int64(0),
					Title:     pointer.String("test"),
					HTMLURL:   pointer.String("TestURL"),
					Labels:    []*github.Label{{Name: pointer.String("famed")}},
					Number:    pointer.Int(0),
					CreatedAt: pointer.Time(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)),
					ClosedAt:  pointer.Time(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
				Label: &github.Label{Name: pointer.String("famed")},
				Repo: &github.Repository{
					Name:  pointer.String("test"),
					Owner: &github.User{Login: pointer.String("test")},
				},
			},
			PullRequest:     &gitlib.PullRequest{URL: "test"},
			ExpectedComment: "### Famed could not generate a reward suggestion.\nReason: The issue is missing an assignee.",
		},
		{
			Name: "Close - No Label",
			Event: &github.IssuesEvent{
				Action: pointer.String("closed"),
				Issue: &github.Issue{
					ID:        pointer.Int64(0),
					Title:     pointer.String("test"),
					HTMLURL:   pointer.String("TestURL"),
					Labels:    []*github.Label{{Name: pointer.String("famed")}},
					Number:    pointer.Int(0),
					Assignees: []*github.User{{Login: pointer.String("test")}},
					CreatedAt: pointer.Time(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)),
					ClosedAt:  pointer.Time(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
				Assignee: &github.User{Login: pointer.String("test")},
				Label:    &github.Label{Name: pointer.String("famed")},
				Repo: &github.Repository{
					Name:  pointer.String("test"),
					Owner: &github.User{Login: pointer.String("test")},
				},
			},
			PullRequest:     &gitlib.PullRequest{URL: "test"},
			ExpectedComment: "### Famed could not generate a reward suggestion.\nReason: The issue is missing a severity label.",
		},
		{
			Name: "Close - Multiple Labels",
			Event: &github.IssuesEvent{
				Action: pointer.String("closed"),
				Issue: &github.Issue{
					ID:        pointer.Int64(0),
					Title:     pointer.String("test"),
					HTMLURL:   pointer.String("TestURL"),
					Labels:    []*github.Label{{Name: pointer.String("famed")}, {Name: pointer.String("high")}, {Name: pointer.String("low")}},
					Number:    pointer.Int(0),
					Assignees: []*github.User{{Login: pointer.String("test")}},
					CreatedAt: pointer.Time(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)),
					ClosedAt:  pointer.Time(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
				Assignee: &github.User{Login: pointer.String("test")},
				Repo: &github.Repository{
					Name:  pointer.String("test"),
					Owner: &github.User{Login: pointer.String("test")},
				},
			},
			PullRequest:     &gitlib.PullRequest{URL: "test"},
			ExpectedComment: "### Famed could not generate a reward suggestion.\nReason: The issue has more than one severity label.",
		},
		{
			Name: "Close - No events",
			Event: &github.IssuesEvent{
				Action: pointer.String("closed"),
				Issue: &github.Issue{
					ID:        pointer.Int64(0),
					Title:     pointer.String("test"),
					HTMLURL:   pointer.String("TestURL"),
					Labels:    []*github.Label{{Name: pointer.String("famed")}, {Name: pointer.String("high")}},
					Number:    pointer.Int(0),
					Assignees: []*github.User{{Login: pointer.String("test")}},
					CreatedAt: pointer.Time(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)),
					ClosedAt:  pointer.Time(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
				Assignee: &github.User{Login: pointer.String("test")},
				Repo: &github.Repository{
					Name:  pointer.String("test"),
					Owner: &github.User{Login: pointer.String("test")},
				},
			},
			PullRequest:     &gitlib.PullRequest{URL: "test"},
			ExpectedComment: "### Famed could not generate a reward suggestion.\nReason: The data provided by GitHub is not sufficient to generate a reward suggestion.\nThis might be due to an assignment after the issue has been closed. Please assign assignees in the open state.",
		},
		// Commented out for DevConnect
		//{
		//	Name: "Close - No pull request",
		//	Event: &github.IssuesEvent{
		//		Action: pointer.String("closed"),
		//		Issue: &github.Issue{
		//			ID:        pointer.Int64(0),
		//			Title:     pointer.String("test"),
		//			Labels:    []*github.Label{{Name: pointer.String("famed")}, {Name: pointer.String("high")}},
		//			Number:    pointer.Int(0),
		//			Assignees: []*github.User{{Login: pointer.String("test")}},
		//			CreatedAt: pointer.Time(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)),
		//			ClosedAt:  pointer.Time(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)),
		//		},
		//		Assignee: &github.User{Login: pointer.String("test")},
		//		Repo: &github.Repository{
		//			Name:  pointer.String("test"),
		//			Owner: &github.User{Login: pointer.String("test")},
		//		},
		//	},
		//	ExpectedComment: "### Famed could not generate a reward suggestion.\nReason: The issue is missing a pull request.",
		//},
		{
			Name: "Close - Valid",
			Event: &github.IssuesEvent{
				Action: pointer.String("closed"),
				Issue: &github.Issue{
					ID:        pointer.Int64(0),
					Title:     pointer.String("test"),
					HTMLURL:   pointer.String("TestURL"),
					Labels:    []*github.Label{{Name: pointer.String("famed")}, {Name: pointer.String("high")}},
					Number:    pointer.Int(0),
					Assignees: []*github.User{{Login: pointer.String("test")}},
					CreatedAt: pointer.Time(time.Date(2021, 12, 1, 0, 0, 0, 0, time.UTC)),
					ClosedAt:  pointer.Time(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
				Assignee: &github.User{Login: pointer.String("test")},
				Repo: &github.Repository{
					Name:  pointer.String("test"),
					Owner: &github.User{Login: pointer.String("test")},
				},
			},
			PullRequest: &gitlib.PullRequest{URL: "test"},
			Events: []gitlib.IssueEvent{
				{
					Event:     "assigned",
					CreatedAt: time.Date(2021, 12, 1, 0, 0, 0, 0, time.UTC),
					Assignee:  &gitlib.User{Login: "test"},
				},
			},
			ExpectedComment: "@test - you Got Famed! 💎 Check out your new score here: https://www.famed.morphysm.com/teams/test/test\n| Contributor | Time | Reward |\n| ----------- | ----------- | ----------- |\n|test|744h0m0s|674 POINTS|",
		},
		{
			Name: "Close - Valid - Migrated",
			Event: &github.IssuesEvent{
				Action: pointer.String("closed"),
				Issue: &github.Issue{
					ID:        pointer.Int64(0),
					Title:     pointer.String("Famed Retroactive Rewards"),
					Body:      pointer.String("**UID:** CL-2021-45\n\n**Severity:** info\n\n**Type:** DoS\n\n**Affected Clients:** All clients\n\n**Summary:** Clients provide clear warnings against doing this.\n\n**Links:** \n\n**Reported:** 2021-08-28\n\n**Fixed:** 2021-08-28\n\n**Published:** 2021-12-01\n\n**Bounty Hunter:** KilianKae \n\n**Bounty Points:** 100"),
					HTMLURL:   pointer.String("TestURL"),
					Labels:    []*github.Label{{Name: pointer.String("famed")}, {Name: pointer.String("high")}},
					Number:    pointer.Int(0),
					Assignees: []*github.User{{Login: pointer.String("test")}},
					CreatedAt: pointer.Time(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
					ClosedAt:  pointer.Time(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
				Assignee: &github.User{Login: pointer.String("test")},
				Repo: &github.Repository{
					Name:  pointer.String("test"),
					Owner: &github.User{Login: pointer.String("test")},
				},
			},
			PullRequest: &gitlib.PullRequest{URL: "test"},
			Events: []gitlib.IssueEvent{
				{
					Event:     "assigned",
					CreatedAt: time.Date(2021, 12, 1, 0, 0, 0, 0, time.UTC),
					Assignee:  &gitlib.User{Login: "test"},
				},
			},
			ExpectedComment: "@test - you Got Famed! 💎 Check out your new score here: https://www.famed.morphysm.com/teams/test/test\n| Contributor | Time | Reward |\n| ----------- | ----------- | ----------- |\n|test|0s|3000 POINTS|",
		},
		{
			Name: "Close - Valid - Multiple Assignees",
			Event: &github.IssuesEvent{
				Action: pointer.String("closed"),
				Issue: &github.Issue{
					ID:        pointer.Int64(0),
					Title:     pointer.String("test"),
					HTMLURL:   pointer.String("TestURL"),
					Labels:    []*github.Label{{Name: pointer.String("famed")}, {Name: pointer.String("high")}},
					Number:    pointer.Int(0),
					Assignees: []*github.User{{Login: pointer.String("test1")}, {Login: pointer.String("test2")}},
					CreatedAt: pointer.Time(time.Date(2021, 12, 1, 0, 0, 0, 0, time.UTC)),
					ClosedAt:  pointer.Time(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
				Repo: &github.Repository{
					Name:  pointer.String("test"),
					Owner: &github.User{Login: pointer.String("testOwner")},
				},
			},
			PullRequest: &gitlib.PullRequest{URL: "test"},
			Events: []gitlib.IssueEvent{
				{
					Event:     "assigned",
					CreatedAt: time.Date(2021, 12, 1, 0, 0, 0, 0, time.UTC),
					Assignee:  &gitlib.User{Login: "test1"},
				},
				{
					Event:     "assigned",
					CreatedAt: time.Date(2021, 12, 1, 0, 0, 0, 0, time.UTC),
					Assignee:  &gitlib.User{Login: "test2"},
				},
			},
			ExpectedComment: "@test1 @test2 - you Got Famed! 💎 Check out your new score here: https://www.famed.morphysm.com/teams/testOwner/test\n| Contributor | Time | Reward |\n| ----------- | ----------- | ----------- |\n|test1|744h0m0s|337 POINTS|\n|test2|744h0m0s|337 POINTS|",
		},
		// Eligible comment
		{
			Name: "Assigned - Missing data",
			Event: &github.IssuesEvent{
				Action: pointer.String("assigned"),
				Issue: &github.Issue{
					Labels: []*github.Label{{Name: pointer.String("famed")}},
				},
				Assignee: &github.User{Login: pointer.String("test")},
			},
			ExpectedErr: &echo.HTTPError{Code: 400, Message: famed.ErrEventMissingData.Error()},
		},
		{
			Name: "Unassigned - Valid - Non present",
			Event: &github.IssuesEvent{
				Action: pointer.String("unassigned"),
				Issue: &github.Issue{
					ID:        pointer.Int64(0),
					Number:    pointer.Int(0),
					Title:     pointer.String("Test"),
					HTMLURL:   pointer.String("TestURL"),
					CreatedAt: pointer.Time(time.Date(2021, 12, 1, 0, 0, 0, 0, time.UTC)),
					Labels:    []*github.Label{{Name: pointer.String("famed")}},
				},
				Assignee: &github.User{Login: pointer.String("test")},
				Repo: &github.Repository{
					Name:  pointer.String("test"),
					Owner: &github.User{Login: pointer.String("test")},
				},
			},
			ExpectedComment: "🤖 Assignees for Issue **Test #0** are now eligible to Get Famed." +
				"\n\n❌ Add assignees to track contribution times of the issue \U0001F9B8\u200d♀️\U0001F9B9️" +
				"\n❌ Add a single severity (CVSS) label to compute the score 🏷️️" +
				//"\n❌ Link a PR when closing the issue ♻️ \U0001F9B8\u200d♀️\U0001F9B9" +
				"\n" +
				"\nHappy hacking! \U0001F9BE💙❤️️",
		},
		{
			Name: "Assigned - Valid - Assignee present",
			Event: &github.IssuesEvent{
				Action: pointer.String("assigned"),
				Issue: &github.Issue{
					ID:        pointer.Int64(0),
					Number:    pointer.Int(0),
					Title:     pointer.String("Test"),
					HTMLURL:   pointer.String("TestURL"),
					CreatedAt: pointer.Time(time.Date(2021, 12, 1, 0, 0, 0, 0, time.UTC)),
					Labels:    []*github.Label{{Name: pointer.String("famed")}},
					Assignees: []*github.User{{Login: pointer.String("test")}},
				},
				Assignee: &github.User{Login: pointer.String("test")},
				Repo: &github.Repository{
					Name:  pointer.String("test"),
					Owner: &github.User{Login: pointer.String("test")},
				},
			},
			ExpectedComment: "🤖 Assignees for Issue **Test #0** are now eligible to Get Famed." +
				"\n\n✅ Add assignees to track contribution times of the issue \U0001F9B8\u200d♀️\U0001F9B9️" +
				"\n❌ Add a single severity (CVSS) label to compute the score 🏷️️" +
				//"\n❌ Link a PR when closing the issue ♻️ \U0001F9B8\u200d♀️\U0001F9B9" +
				"\n" +
				"\nHappy hacking! \U0001F9BE💙❤️️",
		},
		{
			Name: "Assigned - Valid - Label present",
			Event: &github.IssuesEvent{
				Action: pointer.String("assigned"),
				Issue: &github.Issue{
					ID:        pointer.Int64(0),
					Number:    pointer.Int(0),
					Title:     pointer.String("Test"),
					HTMLURL:   pointer.String("TestURL"),
					CreatedAt: pointer.Time(time.Date(2021, 12, 1, 0, 0, 0, 0, time.UTC)),
					Labels:    []*github.Label{{Name: pointer.String("famed")}, {Name: pointer.String("high")}},
					Assignees: []*github.User{{Login: pointer.String("test")}},
				},
				Assignee: &github.User{Login: pointer.String("test")},
				Repo: &github.Repository{
					Name:  pointer.String("test"),
					Owner: &github.User{Login: pointer.String("test")},
				},
			},
			ExpectedComment: "🤖 Assignees for Issue **Test #0** are now eligible to Get Famed." +
				"\n\n✅ Add assignees to track contribution times of the issue \U0001F9B8\u200d♀️\U0001F9B9️" +
				"\n✅ Add a single severity (CVSS) label to compute the score 🏷️️" +
				//"\n❌ Link a PR when closing the issue ♻️ \U0001F9B8\u200d♀️\U0001F9B9" +
				"\n" +
				"\nHappy hacking! \U0001F9BE💙❤️️",
		},
		{
			Name: "Assigned - Valid - PR present",
			Event: &github.IssuesEvent{
				Action: pointer.String("assigned"),
				Issue: &github.Issue{
					ID:        pointer.Int64(0),
					Number:    pointer.Int(0),
					Title:     pointer.String("Test"),
					HTMLURL:   pointer.String("TestURL"),
					CreatedAt: pointer.Time(time.Date(2021, 12, 1, 0, 0, 0, 0, time.UTC)),
					Labels:    []*github.Label{{Name: pointer.String("famed")}},
					Assignees: []*github.User{{Login: pointer.String("test")}},
				},
				Assignee: &github.User{Login: pointer.String("test")},
				Repo: &github.Repository{
					Name:  pointer.String("test"),
					Owner: &github.User{Login: pointer.String("test")},
				},
			},
			PullRequest: &gitlib.PullRequest{URL: "test"},
			ExpectedComment: "🤖 Assignees for Issue **Test #0** are now eligible to Get Famed." +
				"\n\n✅ Add assignees to track contribution times of the issue \U0001F9B8\u200d♀️\U0001F9B9️" +
				"\n❌ Add a single severity (CVSS) label to compute the score 🏷️️" +
				//"\n✅ Link a PR when closing the issue ♻️ \U0001F9B8\u200d♀️\U0001F9B9" +
				"\n" +
				"\nHappy hacking! \U0001F9BE💙❤️️",
		},
		{
			Name: "Assigned - Valid - All present",
			Event: &github.IssuesEvent{
				Action: pointer.String("assigned"),
				Issue: &github.Issue{
					ID:        pointer.Int64(0),
					Number:    pointer.Int(0),
					Title:     pointer.String("Test"),
					HTMLURL:   pointer.String("TestURL"),
					CreatedAt: pointer.Time(time.Date(2021, 12, 1, 0, 0, 0, 0, time.UTC)),
					Labels:    []*github.Label{{Name: pointer.String("famed")}, {Name: pointer.String("high")}},
					Assignees: []*github.User{{Login: pointer.String("test")}},
				},
				Assignee: &github.User{Login: pointer.String("test")},
				Repo: &github.Repository{
					Name:  pointer.String("test"),
					Owner: &github.User{Login: pointer.String("test")},
				},
			},
			PullRequest: &gitlib.PullRequest{URL: "test"},
			ExpectedComment: "🤖 Assignees for Issue **Test #0** are now eligible to Get Famed." +
				"\n\n✅ Add assignees to track contribution times of the issue \U0001F9B8\u200d♀️\U0001F9B9️" +
				"\n✅ Add a single severity (CVSS) label to compute the score 🏷️️" +
				//"\n✅ Link a PR when closing the issue ♻️ \U0001F9B8\u200d♀️\U0001F9B9" +
				"\n\nHappy hacking! \U0001F9BE💙❤️️",
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			// GIVEN
			e := echo.New()
			b := new(bytes.Buffer)
			err := json.NewEncoder(b).Encode(testCase.Event)
			assert.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/github/webhooks/event", b)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			req.Header.Set(github.EventTypeHeader, "issues")
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			fakeInstallationClient := &githubfakes.FakeInstallationClient{}
			fakeInstallationClient.GetIssueEventsReturns(testCase.Events, nil)
			fakeInstallationClient.GetIssuePullRequestReturns(testCase.PullRequest, nil)
			cl, _ := gitlib.NewInstallationClient("", nil, nil, "", "famed", nil)
			fakeInstallationClient.ValidateWebHookEventStub = cl.ValidateWebHookEvent

			githubHandler := famed.NewHandler(nil, fakeInstallationClient, famedConfig)

			// WHEN
			err = githubHandler.PostEvent(ctx)

			// THEN
			if testCase.ExpectedComment != "" {
				assert.Equal(t, 1, fakeInstallationClient.PostCommentCallCount())
				if fakeInstallationClient.PostCommentCallCount() == 1 {
					_, _, _, _, comment := fakeInstallationClient.PostCommentArgsForCall(0)
					assert.Equal(t, testCase.ExpectedComment, comment)
				}
			} else {
				assert.Equal(t, testCase.ExpectedErr, err)
			}
		})
	}
}
