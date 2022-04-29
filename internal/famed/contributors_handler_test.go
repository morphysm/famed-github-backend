package famed_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/go-github/v41/github"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	gitlib "github.com/morphysm/famed-github-backend/internal/client/github"
	"github.com/morphysm/famed-github-backend/internal/client/github/githubfakes"
	"github.com/morphysm/famed-github-backend/internal/famed"
	"github.com/morphysm/famed-github-backend/pkg/pointer"
)

func TestGetContributors(t *testing.T) {
	t.Parallel()

	open := time.Date(2022, 4, 4, 0, 0, 0, 0, time.UTC)
	closed := open.Add(24 * time.Hour)
	famedConfig := NewTestConfig()

	testCases := []struct {
		Name             string
		Owner            string
		RepoName         string
		AppInstalled     bool
		Issues           []gitlib.Issue
		Event            *github.IssuesEvent
		Events           []gitlib.IssueEvent
		PullRequest      *gitlib.PullRequest
		ExpectedResponse string
		ExpectedErr      error
	}{
		{
			Name:         "Valid - One Issue",
			Owner:        "testOwner",
			RepoName:     "testRepo",
			AppInstalled: true,
			Issues: []gitlib.Issue{{
				ID:         0,
				Number:     0,
				HTMLURL:    "TestURL",
				Title:      "TestIssue",
				CreatedAt:  open,
				ClosedAt:   &closed,
				Assignees:  []gitlib.User{{Login: "testUser"}},
				Severities: []gitlib.IssueSeverity{gitlib.IssueSeverity("low")},
				Migrated:   false,
			}},
			Event: &github.IssuesEvent{
				Action: pointer.String("closed"),
				Issue: &github.Issue{
					ID:        pointer.Int64(0),
					Title:     pointer.String("testUser"),
					Labels:    []*github.Label{{Name: pointer.String("famed")}, {Name: pointer.String("high")}},
					Number:    pointer.Int(0),
					Assignees: []*github.User{{Login: pointer.String("testUser")}},
					CreatedAt: &open,
					ClosedAt:  &closed,
				},
				Assignee: &github.User{Login: pointer.String("testUser")},
				Repo: &github.Repository{
					Name:  pointer.String("testUser"),
					Owner: &github.User{Login: pointer.String("testUser")},
				},
			},
			PullRequest: &gitlib.PullRequest{URL: "testUser"},
			Events: []gitlib.IssueEvent{
				{
					Event:     "assigned",
					CreatedAt: time.Date(2021, 12, 1, 0, 0, 0, 0, time.UTC),
					Assignee:  &gitlib.User{Login: "testUser"},
				},
			},
			ExpectedResponse: "[{\"login\":\"testUser\",\"avatarUrl\":\"\",\"htmlUrl\":\"\",\"fixCount\":1,\"rewards\":[{\"date\":\"2022-04-05T00:00:00Z\",\"reward\":975,\"url\":\"TestURL\"}],\"rewardSum\":975,\"currency\":\"POINTS\",\"rewardsLastYear\":[{\"month\":\"4.2022\",\"reward\":975},{\"month\":\"3.2022\",\"reward\":0},{\"month\":\"2.2022\",\"reward\":0},{\"month\":\"1.2022\",\"reward\":0},{\"month\":\"12.2021\",\"reward\":0},{\"month\":\"11.2021\",\"reward\":0},{\"month\":\"10.2021\",\"reward\":0},{\"month\":\"9.2021\",\"reward\":0},{\"month\":\"8.2021\",\"reward\":0},{\"month\":\"7.2021\",\"reward\":0},{\"month\":\"6.2021\",\"reward\":0},{\"month\":\"5.2021\",\"reward\":0}],\"timeToDisclosure\":{\"time\":[1440],\"mean\":1440,\"standardDeviation\":0},\"severities\":{\"low\":1},\"meanSeverity\":2}]\n",
		},
		{
			Name:         "Valid - Two Issues",
			Owner:        "testOwner",
			RepoName:     "testRepo",
			AppInstalled: true,
			Issues: []gitlib.Issue{
				{
					ID:         0,
					Number:     0,
					HTMLURL:    "TestURL",
					Title:      "TestIssue",
					CreatedAt:  open,
					ClosedAt:   &closed,
					Assignees:  []gitlib.User{{Login: "testUser"}},
					Severities: []gitlib.IssueSeverity{gitlib.IssueSeverity("low")},
					Migrated:   false,
				},
				{
					ID:         1,
					Number:     1,
					HTMLURL:    "TestURL",
					Title:      "TestIssue",
					CreatedAt:  open,
					ClosedAt:   &closed,
					Assignees:  []gitlib.User{{Login: "testUser"}},
					Severities: []gitlib.IssueSeverity{gitlib.IssueSeverity("low")},
					Migrated:   false,
				},
			},
			Event: &github.IssuesEvent{
				Action: pointer.String("closed"),
				Issue: &github.Issue{
					ID:        pointer.Int64(0),
					Title:     pointer.String("testUser"),
					Labels:    []*github.Label{{Name: pointer.String("famed")}, {Name: pointer.String("high")}},
					Number:    pointer.Int(0),
					Assignees: []*github.User{{Login: pointer.String("testUser")}},
					CreatedAt: &open,
					ClosedAt:  &closed,
				},
				Assignee: &github.User{Login: pointer.String("testUser")},
				Repo: &github.Repository{
					Name:  pointer.String("testUser"),
					Owner: &github.User{Login: pointer.String("testUser")},
				},
			},
			PullRequest: &gitlib.PullRequest{URL: "testUser"},
			Events: []gitlib.IssueEvent{
				{
					Event:     "assigned",
					CreatedAt: time.Date(2021, 12, 1, 0, 0, 0, 0, time.UTC),
					Assignee:  &gitlib.User{Login: "testUser"},
				},
			},
			ExpectedResponse: "[{\"login\":\"testUser\",\"avatarUrl\":\"\",\"htmlUrl\":\"\",\"fixCount\":2,\"rewards\":[{\"date\":\"2022-04-05T00:00:00Z\",\"reward\":975,\"url\":\"TestURL\"},{\"date\":\"2022-04-05T00:00:00Z\",\"reward\":975,\"url\":\"TestURL\"}],\"rewardSum\":1950,\"currency\":\"POINTS\",\"rewardsLastYear\":[{\"month\":\"4.2022\",\"reward\":1950},{\"month\":\"3.2022\",\"reward\":0},{\"month\":\"2.2022\",\"reward\":0},{\"month\":\"1.2022\",\"reward\":0},{\"month\":\"12.2021\",\"reward\":0},{\"month\":\"11.2021\",\"reward\":0},{\"month\":\"10.2021\",\"reward\":0},{\"month\":\"9.2021\",\"reward\":0},{\"month\":\"8.2021\",\"reward\":0},{\"month\":\"7.2021\",\"reward\":0},{\"month\":\"6.2021\",\"reward\":0},{\"month\":\"5.2021\",\"reward\":0}],\"timeToDisclosure\":{\"time\":[1440,1440],\"mean\":1440,\"standardDeviation\":0},\"severities\":{\"low\":2},\"meanSeverity\":2}]\n",
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

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/github/repos/%s/%s/contributors", testCase.Owner, testCase.RepoName), b)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.SetParamNames([]string{"owner", "repo_name"}...)
			ctx.SetParamValues([]string{testCase.Owner, testCase.RepoName}...)

			fakeInstallationClient := &githubfakes.FakeInstallationClient{}
			fakeInstallationClient.CheckInstallationReturns(testCase.AppInstalled)
			// TODO testUser for error
			fakeInstallationClient.GetIssuesByRepoReturns(testCase.Issues, nil)
			fakeInstallationClient.GetIssueEventsReturns(testCase.Events, nil)
			fakeInstallationClient.GetIssuePullRequestReturns(testCase.PullRequest, nil)

			githubHandler := famed.NewHandler(nil, fakeInstallationClient, famedConfig)

			// WHEN
			err = githubHandler.GetContributors(ctx)

			// THEN
			assert.Equal(t, testCase.ExpectedErr, err)

			if testCase.ExpectedResponse != "" {
				assert.Equal(t, 1, fakeInstallationClient.CheckInstallationCallCount())
				assert.Equal(t, 1, fakeInstallationClient.GetIssuesByRepoCallCount())
				assert.Equal(t, len(testCase.Issues), fakeInstallationClient.GetIssuePullRequestCallCount())
				assert.Equal(t, testCase.ExpectedResponse, rec.Body.String())
			}
		})
	}
}
