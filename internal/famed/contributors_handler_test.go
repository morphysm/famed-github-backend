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

	"github.com/morphysm/famed-github-backend/internal/famed"
	model2 "github.com/morphysm/famed-github-backend/internal/famed/model"
	model "github.com/morphysm/famed-github-backend/internal/respositories/github/model"
	"github.com/morphysm/famed-github-backend/internal/respositories/github/providers/providersfakes"
	"github.com/morphysm/famed-github-backend/pkg/pointer"
)

func Now() time.Time {
	var nowTime = time.Date(2022, 4, 20, 0, 0, 0, 0, time.UTC)
	return nowTime
}

func NewTestConfig() model2.Config {
	rewards := map[model.IssueSeverity]float64{
		model.Info:     0,
		model.Low:      1000,
		model.Medium:   2000,
		model.High:     3000,
		model.Critical: 4000,
	}
	labels := map[string]model.Label{
		"famed": {
			Name:        "famed",
			Color:       "testColor",
			Description: "testDescription",
		},
	}

	return model2.NewFamedConfig("POINTS",
		rewards,
		labels,
		40,
		"b",
	)
}

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
		Issues           []model.Issue
		Event            *github.IssuesEvent
		Events           []model.IssueEvent
		PullRequest      *string
		ExpectedResponse string
		ExpectedErr      error
	}{
		{
			Name:         "Valid - One issue",
			Owner:        "testOwner",
			RepoName:     "testRepo",
			AppInstalled: true,
			Issues: []model.Issue{{
				ID:         0,
				Number:     0,
				HTMLURL:    "TestURL",
				Title:      "TestIssue",
				CreatedAt:  open,
				ClosedAt:   &closed,
				Assignees:  []model.User{{Login: "testUser"}},
				Severities: []model.IssueSeverity{model.IssueSeverity("low")},
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
			PullRequest: pointer.String("testUser"),
			Events: []model.IssueEvent{
				{
					Event:     "assigned",
					CreatedAt: time.Date(2021, 12, 1, 0, 0, 0, 0, time.UTC),
					Assignee:  &model.User{Login: "testUser"},
				},
			},
			ExpectedResponse: "[{\"login\":\"testUser\",\"avatarUrl\":\"\",\"htmlUrl\":\"\",\"fixCount\":1,\"rewards\":[{\"date\":\"2022-04-05T00:00:00Z\",\"reward\":975,\"url\":\"TestURL\"}],\"rewardSum\":975,\"currency\":\"POINTS\",\"rewardsLastYear\":[{\"month\":\"4.2022\",\"reward\":975},{\"month\":\"3.2022\",\"reward\":0},{\"month\":\"2.2022\",\"reward\":0},{\"month\":\"1.2022\",\"reward\":0},{\"month\":\"12.2021\",\"reward\":0},{\"month\":\"11.2021\",\"reward\":0},{\"month\":\"10.2021\",\"reward\":0},{\"month\":\"9.2021\",\"reward\":0},{\"month\":\"8.2021\",\"reward\":0},{\"month\":\"7.2021\",\"reward\":0},{\"month\":\"6.2021\",\"reward\":0},{\"month\":\"5.2021\",\"reward\":0}],\"timeToDisclosure\":{\"time\":[1440],\"mean\":1440,\"standardDeviation\":0},\"severities\":{\"low\":1},\"meanSeverity\":2}]\n",
		},
		{
			Name:         "Valid - Two Issues",
			Owner:        "testOwner",
			RepoName:     "testRepo",
			AppInstalled: true,
			Issues: []model.Issue{
				{
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
				{
					ID:         1,
					Number:     1,
					HTMLURL:    "TestURL",
					Title:      "TestIssue",
					CreatedAt:  open,
					ClosedAt:   &closed,
					Assignees:  []model.User{{Login: "testUser"}},
					Severities: []model.IssueSeverity{model.IssueSeverity("low")},
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
			PullRequest: pointer.String("testUser"),
			Events: []model.IssueEvent{
				{
					Event:     "assigned",
					CreatedAt: time.Date(2021, 12, 1, 0, 0, 0, 0, time.UTC),
					Assignee:  &model.User{Login: "testUser"},
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

			fakeInstallationClient := &providersfakes.FakeInstallationClient{}
			fakeInstallationClient.CheckInstallationReturns(testCase.AppInstalled)
			enrichedIssues := make(map[int]model.EnrichedIssue, len(testCase.Issues))
			for _, issue := range testCase.Issues {
				enrichedIssues[issue.Number] = model.NewEnrichIssue(issue, testCase.PullRequest, testCase.Events)
			}
			fakeInstallationClient.GetEnrichedIssuesReturns(enrichedIssues, nil)

			githubHandler := famed.NewHandler(nil, fakeInstallationClient, famedConfig, Now)

			// WHEN
			err = githubHandler.GetBlueTeam(ctx)

			// THEN
			assert.Equal(t, testCase.ExpectedErr, err)

			if testCase.ExpectedResponse != "" {
				assert.Equal(t, 1, fakeInstallationClient.GetEnrichedIssuesCallCount())
				assert.Equal(t, testCase.ExpectedResponse, rec.Body.String())
			}
		})
	}
}
