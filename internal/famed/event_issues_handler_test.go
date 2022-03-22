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
	"github.com/morphysm/famed-github-backend/internal/client/installation"
	"github.com/morphysm/famed-github-backend/internal/client/installation/installationfakes"
	"github.com/morphysm/famed-github-backend/internal/config"
	"github.com/morphysm/famed-github-backend/internal/famed"
	"github.com/morphysm/famed-github-backend/pkg/pointers"
	"github.com/stretchr/testify/assert"
)

func TestPostIssuesEvent(t *testing.T) {
	t.Parallel()

	rewards := map[config.IssueSeverity]float64{
		config.CVSSNone:     0,
		config.CVSSLow:      1,
		config.CVSSMedium:   2,
		config.CVSSHigh:     3,
		config.CVSSCritical: 4,
	}
	famedConfig := famed.Config{
		Labels:   map[string]installation.Label{"famed": {Name: "famed"}},
		Currency: "eth",
		Rewards:  rewards,
	}

	testCases := []struct {
		Name            string
		Event           *github.IssuesEvent
		Events          []*github.IssueEvent
		PullRequest     *installation.PullRequest
		ExpectedComment string
		ExpectedErr     error
	}{
		{
			Name: "Closed - Empty event",
			Event: &github.IssuesEvent{
				Action: pointers.String("closed"),
			},
			ExpectedComment: "",
			ExpectedErr:     famed.ErrEventMissingData,
		},
		{
			Name: "Closed - No Assignee",
			Event: &github.IssuesEvent{
				Action: pointers.String("closed"),
				Issue: &github.Issue{
					ID:        pointers.Int64(1),
					Labels:    []*github.Label{{Name: pointers.String("famed")}},
					Number:    pointers.Int(0),
					CreatedAt: pointers.Time(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)),
					ClosedAt:  pointers.Time(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
				Label: &github.Label{Name: pointers.String("famed")},
				Repo: &github.Repository{
					Name:  pointers.String("test"),
					Owner: &github.User{Login: pointers.String("test")},
				},
			},
			PullRequest:     &installation.PullRequest{URL: "test"},
			ExpectedComment: "### Famed could not generate a reward suggestion. \nReason: The issue is missing an assignee.",
		},
		{
			Name: "Closed - No Label",
			Event: &github.IssuesEvent{
				Action: pointers.String("closed"),
				Issue: &github.Issue{
					ID:        pointers.Int64(1),
					Labels:    []*github.Label{{Name: pointers.String("famed")}},
					Number:    pointers.Int(0),
					Assignee:  &github.User{Login: pointers.String("test")},
					CreatedAt: pointers.Time(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)),
					ClosedAt:  pointers.Time(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
				Assignee: &github.User{Login: pointers.String("test")},
				Label:    &github.Label{Name: pointers.String("famed")},
				Repo: &github.Repository{
					Name:  pointers.String("test"),
					Owner: &github.User{Login: pointers.String("test")},
				},
			},
			PullRequest:     &installation.PullRequest{URL: "test"},
			ExpectedComment: "### Famed could not generate a reward suggestion. \nReason: The issue is missing a severity label.",
		},
		{
			Name: "Closed - Multiple Labels",
			Event: &github.IssuesEvent{
				Action: pointers.String("closed"),
				Issue: &github.Issue{
					ID:        pointers.Int64(1),
					Labels:    []*github.Label{{Name: pointers.String("famed")}, {Name: pointers.String("high")}, {Name: pointers.String("low")}},
					Number:    pointers.Int(0),
					Assignee:  &github.User{Login: pointers.String("test")},
					CreatedAt: pointers.Time(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)),
					ClosedAt:  pointers.Time(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
				Assignee: &github.User{Login: pointers.String("test")},
				Repo: &github.Repository{
					Name:  pointers.String("test"),
					Owner: &github.User{Login: pointers.String("test")},
				},
			},
			PullRequest:     &installation.PullRequest{URL: "test"},
			ExpectedComment: "### Famed could not generate a reward suggestion. \nReason: The issue has more than one severity label.",
		},
		{
			Name: "Closed - No events",
			Event: &github.IssuesEvent{
				Action: pointers.String("closed"),
				Issue: &github.Issue{
					ID:        pointers.Int64(1),
					Labels:    []*github.Label{{Name: pointers.String("famed")}, {Name: pointers.String("high")}},
					Number:    pointers.Int(0),
					Assignee:  &github.User{Login: pointers.String("test")},
					CreatedAt: pointers.Time(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)),
					ClosedAt:  pointers.Time(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
				Assignee: &github.User{Login: pointers.String("test")},
				Repo: &github.Repository{
					Name:  pointers.String("test"),
					Owner: &github.User{Login: pointers.String("test")},
				},
			},
			PullRequest:     &installation.PullRequest{URL: "test"},
			ExpectedComment: "### Famed could not generate a reward suggestion. \nReason: The data provided by GitHub is not sufficient to generate a reward suggestion.",
		},
		{
			Name: "Closed - No pull request",
			Event: &github.IssuesEvent{
				Action: pointers.String("closed"),
				Issue: &github.Issue{
					ID:        pointers.Int64(1),
					Labels:    []*github.Label{{Name: pointers.String("famed")}, {Name: pointers.String("high")}},
					Number:    pointers.Int(0),
					Assignee:  &github.User{Login: pointers.String("test")},
					CreatedAt: pointers.Time(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)),
					ClosedAt:  pointers.Time(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
				Assignee: &github.User{Login: pointers.String("test")},
				Repo: &github.Repository{
					Name:  pointers.String("test"),
					Owner: &github.User{Login: pointers.String("test")},
				},
			},
			ExpectedComment: "### Famed could not generate a reward suggestion. \nReason: The issue is missing a pull request.",
		},
		{
			Name: "Closed - Valid",
			Event: &github.IssuesEvent{
				Action: pointers.String("closed"),
				Issue: &github.Issue{
					ID:        pointers.Int64(1),
					Labels:    []*github.Label{{Name: pointers.String("famed")}, {Name: pointers.String("high")}},
					Number:    pointers.Int(0),
					Assignee:  &github.User{Login: pointers.String("test")},
					CreatedAt: pointers.Time(time.Date(2021, 12, 1, 0, 0, 0, 0, time.UTC)),
					ClosedAt:  pointers.Time(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
				Assignee: &github.User{Login: pointers.String("test")},
				Repo: &github.Repository{
					Name:  pointers.String("test"),
					Owner: &github.User{Login: pointers.String("test")},
				},
			},
			PullRequest: &installation.PullRequest{URL: "test"},
			Events: []*github.IssueEvent{
				{
					Event:     pointers.String("assigned"),
					CreatedAt: pointers.Time(time.Date(2021, 12, 1, 0, 0, 0, 0, time.UTC)),
					Assignee:  &github.User{Login: pointers.String("test")},
				},
			},
			ExpectedComment: "### Famed suggests:\n| Contributor | Time | Reward |\n| ----------- | ----------- | ----------- |\n|test|744h0m0s|0.675000 eth|",
		},
		// Eligible comment
		{
			Name: "Assigned - Missing data",
			Event: &github.IssuesEvent{
				Action: pointers.String("assigned"),
				Issue: &github.Issue{
					Labels: []*github.Label{{Name: pointers.String("famed")}},
				},
				Assignee: &github.User{Login: pointers.String("test")},
			},
			ExpectedErr: famed.ErrEventMissingData,
		},
		{
			Name: "Assigned - Valid - Non present",
			Event: &github.IssuesEvent{
				Action: pointers.String("assigned"),
				Issue: &github.Issue{
					ID:     pointers.Int64(1),
					Number: pointers.Int(0),
					Title:  pointers.String("Test"),
					Labels: []*github.Label{{Name: pointers.String("famed")}},
				},
				Assignee: &github.User{Login: pointers.String("test")},
				Repo: &github.Repository{
					Name:  pointers.String("test"),
					Owner: &github.User{Login: pointers.String("test")},
				},
			},
			ExpectedComment: "🤖 Assignees for WrappedIssue **Test #0** are now eligible to Get Famed." +
				"\n- ❌ Add assignees to track contribution times of the issue \U0001F9B8\u200d♀️\U0001F9B9️" +
				"\n- ❌ Add a severity (CVSS) label to compute the score 🏷️️" +
				"\n- ❌ Link a PR when closing the issue ♻️ \U0001F9B8\u200d♀️\U0001F9B9" +
				"\n" +
				"\nHappy hacking! \U0001F9BE💙❤️️",
		},
		{
			Name: "Assigned - Valid - Assignee present",
			Event: &github.IssuesEvent{
				Action: pointers.String("assigned"),
				Issue: &github.Issue{
					ID:       pointers.Int64(1),
					Number:   pointers.Int(0),
					Title:    pointers.String("Test"),
					Labels:   []*github.Label{{Name: pointers.String("famed")}},
					Assignee: &github.User{Login: pointers.String("test")},
				},
				Assignee: &github.User{Login: pointers.String("test")},
				Repo: &github.Repository{
					Name:  pointers.String("test"),
					Owner: &github.User{Login: pointers.String("test")},
				},
			},
			ExpectedComment: "🤖 Assignees for WrappedIssue **Test #0** are now eligible to Get Famed." +
				"\n- ✅ Add assignees to track contribution times of the issue \U0001F9B8\u200d♀️\U0001F9B9️" +
				"\n- ❌ Add a severity (CVSS) label to compute the score 🏷️️" +
				"\n- ❌ Link a PR when closing the issue ♻️ \U0001F9B8\u200d♀️\U0001F9B9" +
				"\n" +
				"\nHappy hacking! \U0001F9BE💙❤️️",
		},
		{
			Name: "Assigned - Valid - Label present",
			Event: &github.IssuesEvent{
				Action: pointers.String("assigned"),
				Issue: &github.Issue{
					ID:     pointers.Int64(1),
					Number: pointers.Int(0),
					Title:  pointers.String("Test"),
					Labels: []*github.Label{{Name: pointers.String("famed")}, {Name: pointers.String("high")}},
				},
				Assignee: &github.User{Login: pointers.String("test")},
				Repo: &github.Repository{
					Name:  pointers.String("test"),
					Owner: &github.User{Login: pointers.String("test")},
				},
			},
			ExpectedComment: "🤖 Assignees for WrappedIssue **Test #0** are now eligible to Get Famed." +
				"\n- ❌ Add assignees to track contribution times of the issue \U0001F9B8\u200d♀️\U0001F9B9️" +
				"\n- ✅ Add a severity (CVSS) label to compute the score 🏷️️" +
				"\n- ❌ Link a PR when closing the issue ♻️ \U0001F9B8\u200d♀️\U0001F9B9" +
				"\n" +
				"\nHappy hacking! \U0001F9BE💙❤️️",
		},
		{
			Name: "Assigned - Valid - PR present",
			Event: &github.IssuesEvent{
				Action: pointers.String("assigned"),
				Issue: &github.Issue{
					ID:     pointers.Int64(1),
					Number: pointers.Int(0),
					Title:  pointers.String("Test"),
					Labels: []*github.Label{{Name: pointers.String("famed")}},
				},
				Assignee: &github.User{Login: pointers.String("test")},
				Repo: &github.Repository{
					Name:  pointers.String("test"),
					Owner: &github.User{Login: pointers.String("test")},
				},
			},
			PullRequest: &installation.PullRequest{URL: "test"},
			ExpectedComment: "🤖 Assignees for WrappedIssue **Test #0** are now eligible to Get Famed." +
				"\n- ❌ Add assignees to track contribution times of the issue \U0001F9B8\u200d♀️\U0001F9B9️" +
				"\n- ❌ Add a severity (CVSS) label to compute the score 🏷️️" +
				"\n- ✅ Link a PR when closing the issue ♻️ \U0001F9B8\u200d♀️\U0001F9B9" +
				"\n" +
				"\nHappy hacking! \U0001F9BE💙❤️️",
		},
		{
			Name: "Assigned - Valid - All present",
			Event: &github.IssuesEvent{
				Action: pointers.String("assigned"),
				Issue: &github.Issue{
					ID:       pointers.Int64(1),
					Number:   pointers.Int(0),
					Title:    pointers.String("Test"),
					Labels:   []*github.Label{{Name: pointers.String("famed")}, {Name: pointers.String("high")}},
					Assignee: &github.User{Login: pointers.String("test")},
				},
				Assignee: &github.User{Login: pointers.String("test")},
				Repo: &github.Repository{
					Name:  pointers.String("test"),
					Owner: &github.User{Login: pointers.String("test")},
				},
			},
			PullRequest: &installation.PullRequest{URL: "test"},
			ExpectedComment: "🤖 Assignees for WrappedIssue **Test #0** are now eligible to Get Famed." +
				"\n- ✅ Add assignees to track contribution times of the issue \U0001F9B8\u200d♀️\U0001F9B9️" +
				"\n- ✅ Add a severity (CVSS) label to compute the score 🏷️️" +
				"\n- ✅ Link a PR when closing the issue ♻️ \U0001F9B8\u200d♀️\U0001F9B9" +
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

			fakeInstallationClient := &installationfakes.FakeClient{}
			fakeInstallationClient.GetIssueEventsReturns(testCase.Events, nil)
			fakeInstallationClient.GetIssuePullRequestReturns(testCase.PullRequest, nil)

			githubHandler := famed.NewHandler(nil, fakeInstallationClient, nil, famedConfig)

			// WHEN
			err = githubHandler.PostEvent(ctx)

			// THEN
			if testCase.ExpectedComment != "" {
				assert.Equal(t, 1, fakeInstallationClient.PostCommentCallCount())
				_, _, _, _, comment := fakeInstallationClient.PostCommentArgsForCall(0)
				assert.Equal(t, testCase.ExpectedComment, comment)
			}
			assert.Equal(t, testCase.ExpectedErr, err)
		})
	}
}
