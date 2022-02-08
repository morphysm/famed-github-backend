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
	"github.com/morphysm/famed-github-backend/internal/client/currency/currencyfakes"
	"github.com/morphysm/famed-github-backend/internal/client/installation/installationfakes"
	"github.com/morphysm/famed-github-backend/internal/famed"
	"github.com/morphysm/famed-github-backend/pkg/pointers"
	"github.com/stretchr/testify/assert"
)

func TestPostEvent(t *testing.T) {
	t.Parallel()

	rewards := map[famed.IssueSeverity]float64{
		famed.IssueSeverityNone:     0,
		famed.IssueSeverityLow:      1,
		famed.IssueSeverityMedium:   2,
		famed.IssueSeverityHigh:     3,
		famed.IssueSeverityCritical: 4,
	}
	famedConfig := famed.Config{
		Label:    "famed",
		Currency: "eth",
		Rewards:  rewards,
	}

	testCases := []struct {
		Name            string
		Event           github.IssuesEvent
		Events          []*github.IssueEvent
		ExpectedComment string
		ExpectedErr     error
	}{
		{
			Name: "Empty closed event",
			Event: github.IssuesEvent{
				Action:   pointers.String("closed"),
				Issue:    &github.Issue{},
				Assignee: &github.User{},
				Label:    &github.Label{},

				// The following fields are only populated by Webhook events.
				Changes:      &github.EditChange{},
				Repo:         &github.Repository{},
				Sender:       &github.User{},
				Installation: &github.Installation{},
			},
			ExpectedComment: "",
			ExpectedErr:     famed.ErrEventMissingData,
		},
		{
			Name: "No Assignee",
			Event: github.IssuesEvent{
				Action: pointers.String("closed"),
				Issue: &github.Issue{
					ID:        pointers.Int64(1),
					Labels:    []*github.Label{{Name: pointers.String("famed")}},
					Number:    pointers.Int(0),
					CreatedAt: pointers.Time(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)),
					ClosedAt:  pointers.Time(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
				Assignee: nil,
				Label:    &github.Label{Name: pointers.String("famed")},

				// The following fields are only populated by Webhook events.
				Changes: nil,
				Repo: &github.Repository{
					Name: pointers.String("test"),
				},
				Sender:       nil,
				Installation: nil,
			},
			ExpectedComment: "### Kudo could not generate a reward suggestion. \nReason: The issue is missing an assignee.",
			ExpectedErr:     nil,
		},
		{
			Name: "No Label",
			Event: github.IssuesEvent{
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

				// The following fields are only populated by Webhook events.
				Changes: nil,
				Repo: &github.Repository{
					Name: pointers.String("test"),
				},
				Sender:       nil,
				Installation: nil,
			},
			ExpectedComment: "### Kudo could not generate a reward suggestion. \nReason: The issue is missing a severity label.",
			ExpectedErr:     nil,
		},
		{
			Name: "Multiple Labels",
			Event: github.IssuesEvent{
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
				Label:    nil,

				// The following fields are only populated by Webhook events.
				Changes: nil,
				Repo: &github.Repository{
					Name: pointers.String("test"),
				},
				Sender:       nil,
				Installation: nil,
			},
			ExpectedComment: "### Kudo could not generate a reward suggestion. \nReason: The issue has more than one severity label.",
			ExpectedErr:     nil,
		},
		{
			Name: "No events",
			Event: github.IssuesEvent{
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
				Label:    nil,

				// The following fields are only populated by Webhook events.
				Changes: nil,
				Repo: &github.Repository{
					Name: pointers.String("test"),
				},
				Sender:       nil,
				Installation: nil,
			},
			ExpectedComment: "### Kudo could not generate a reward suggestion. \nReason: The data provided by GitHub is not sufficient to generate a reward suggestion.",
			ExpectedErr:     nil,
		},
		{
			Name: "Valid",
			Event: github.IssuesEvent{
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
				Label:    nil,

				// The following fields are only populated by Webhook events.
				Changes: nil,
				Repo: &github.Repository{
					Name: pointers.String("test"),
				},
				Sender:       nil,
				Installation: nil,
			},
			Events: []*github.IssueEvent{
				{
					Event:     pointers.String("assigned"),
					CreatedAt: pointers.Time(time.Date(2021, 12, 1, 0, 0, 0, 0, time.UTC)),
					Assignee:  &github.User{Login: pointers.String("test")},
				},
			},
			ExpectedComment: "### Kudo suggests:\n| Contributor | Time | Reward |\n| ----------- | ----------- | ----------- |\n|test|744h0m0s|0.675000 eth|",
			ExpectedErr:     nil,
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
			fakeCurrencyClient := &currencyfakes.FakeClient{}
			fakeCurrencyClient.GetUSDToETHConversionReturns(1, nil)
			githubHandler := famed.NewHandler(fakeInstallationClient, fakeCurrencyClient, nil, 0, famedConfig)

			fakeInstallationClient.GetIssueEventsReturns(testCase.Events, nil)

			// WHEN
			err = githubHandler.PostEvent(ctx)

			// THEN
			if testCase.ExpectedComment != "" {
				assert.Equal(t, 1, fakeInstallationClient.PostCommentCallCount())
				_, _, _, comment := fakeInstallationClient.PostCommentArgsForCall(0)
				assert.Equal(t, testCase.ExpectedComment, comment)
			}
			assert.Equal(t, testCase.ExpectedErr, err)
		})
	}
}
