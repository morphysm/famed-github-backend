package famed_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	gitlib "github.com/morphysm/famed-github-backend/internal/client/github"
	"github.com/morphysm/famed-github-backend/internal/client/github/githubfakes"
	"github.com/morphysm/famed-github-backend/internal/famed"
	"github.com/morphysm/famed-github-backend/pkg/pointer"
	"github.com/stretchr/testify/assert"
)

func TestRedTeam(t *testing.T) {
	t.Parallel()

	open := time.Date(2022, 4, 4, 0, 0, 0, 0, time.UTC)
	closed := open.Add(24 * time.Hour)
	famedConfig := NewTestConfig()
	//redTeam := map[string]string{"TestUser": "TestLogin"}

	testCases := []struct {
		Name             string
		Owner            string
		RepoName         string
		AppInstalled     bool
		Issues           []gitlib.Issue
		ExpectedResponse string
		ExpectedErr      error
	}{
		{
			Name:         "Valid - One Issue",
			Owner:        "testOwner",
			RepoName:     "testRepo",
			AppInstalled: true,
			Issues: []gitlib.Issue{{
				HTMLURL:      "TestURL",
				Severities:   []gitlib.IssueSeverity{gitlib.IssueSeverity("low")},
				CreatedAt:    open,
				ClosedAt:     &closed,
				Migrated:     true,
				RedTeam:      []gitlib.User{{Login: "testUser"}},
				BountyPoints: pointer.Int(975),
			}},
			ExpectedResponse: "[{\"login\":\"testUser\",\"avatarUrl\":\"\",\"htmlUrl\":\"\",\"fixCount\":1,\"rewards\":[{\"date\":\"2022-04-05T00:00:00Z\",\"reward\":975,\"url\":\"TestURL\"}],\"rewardSum\":975,\"currency\":\"POINTS\",\"rewardsLastYear\":[{\"month\":\"4.2022\",\"reward\":975},{\"month\":\"3.2022\",\"reward\":0},{\"month\":\"2.2022\",\"reward\":0},{\"month\":\"1.2022\",\"reward\":0},{\"month\":\"12.2021\",\"reward\":0},{\"month\":\"11.2021\",\"reward\":0},{\"month\":\"10.2021\",\"reward\":0},{\"month\":\"9.2021\",\"reward\":0},{\"month\":\"8.2021\",\"reward\":0},{\"month\":\"7.2021\",\"reward\":0},{\"month\":\"6.2021\",\"reward\":0},{\"month\":\"5.2021\",\"reward\":0}],\"timeToDisclosure\":{\"time\":[1440],\"mean\":1440,\"standardDeviation\":0},\"severities\":{\"low\":1},\"meanSeverity\":2}]\n",
		},
		{
			Name:         "Valid - One Issue - Two in RedTeam",
			Owner:        "testOwner",
			RepoName:     "testRepo",
			AppInstalled: true,
			Issues: []gitlib.Issue{{
				HTMLURL:      "TestURL",
				Severities:   []gitlib.IssueSeverity{gitlib.IssueSeverity("low")},
				CreatedAt:    open,
				ClosedAt:     &closed,
				Migrated:     true,
				RedTeam:      []gitlib.User{{Login: "testUser1"}, {Login: "testUser2"}},
				BountyPoints: pointer.Int(1000),
			}},
			ExpectedResponse: "[{\"login\":\"testUser1\",\"avatarUrl\":\"\",\"htmlUrl\":\"\",\"fixCount\":1,\"rewards\":[{\"date\":\"2022-04-05T00:00:00Z\",\"reward\":500,\"url\":\"TestURL\"}],\"rewardSum\":500,\"currency\":\"POINTS\",\"rewardsLastYear\":[{\"month\":\"4.2022\",\"reward\":500},{\"month\":\"3.2022\",\"reward\":0},{\"month\":\"2.2022\",\"reward\":0},{\"month\":\"1.2022\",\"reward\":0},{\"month\":\"12.2021\",\"reward\":0},{\"month\":\"11.2021\",\"reward\":0},{\"month\":\"10.2021\",\"reward\":0},{\"month\":\"9.2021\",\"reward\":0},{\"month\":\"8.2021\",\"reward\":0},{\"month\":\"7.2021\",\"reward\":0},{\"month\":\"6.2021\",\"reward\":0},{\"month\":\"5.2021\",\"reward\":0}],\"timeToDisclosure\":{\"time\":[1440],\"mean\":1440,\"standardDeviation\":0},\"severities\":{\"low\":1},\"meanSeverity\":2},{\"login\":\"testUser2\",\"avatarUrl\":\"\",\"htmlUrl\":\"\",\"fixCount\":1,\"rewards\":[{\"date\":\"2022-04-05T00:00:00Z\",\"reward\":500,\"url\":\"TestURL\"}],\"rewardSum\":500,\"currency\":\"POINTS\",\"rewardsLastYear\":[{\"month\":\"4.2022\",\"reward\":500},{\"month\":\"3.2022\",\"reward\":0},{\"month\":\"2.2022\",\"reward\":0},{\"month\":\"1.2022\",\"reward\":0},{\"month\":\"12.2021\",\"reward\":0},{\"month\":\"11.2021\",\"reward\":0},{\"month\":\"10.2021\",\"reward\":0},{\"month\":\"9.2021\",\"reward\":0},{\"month\":\"8.2021\",\"reward\":0},{\"month\":\"7.2021\",\"reward\":0},{\"month\":\"6.2021\",\"reward\":0},{\"month\":\"5.2021\",\"reward\":0}],\"timeToDisclosure\":{\"time\":[1440],\"mean\":1440,\"standardDeviation\":0},\"severities\":{\"low\":1},\"meanSeverity\":2}]\n",
		},
		{
			Name:         "Valid - Two Issues - Same RedTeam",
			Owner:        "testOwner",
			RepoName:     "testRepo",
			AppInstalled: true,
			Issues: []gitlib.Issue{
				{
					HTMLURL:      "TestURL",
					CreatedAt:    open,
					ClosedAt:     &closed,
					Severities:   []gitlib.IssueSeverity{gitlib.IssueSeverity("low")},
					Migrated:     true,
					RedTeam:      []gitlib.User{{Login: "testUser"}},
					BountyPoints: pointer.Int(975),
				},
				{
					HTMLURL:      "TestURL",
					CreatedAt:    open,
					ClosedAt:     &closed,
					Severities:   []gitlib.IssueSeverity{gitlib.IssueSeverity("low")},
					Migrated:     true,
					RedTeam:      []gitlib.User{{Login: "testUser"}},
					BountyPoints: pointer.Int(975),
				},
			},
			ExpectedResponse: "[{\"login\":\"testUser\",\"avatarUrl\":\"\",\"htmlUrl\":\"\",\"fixCount\":2,\"rewards\":[{\"date\":\"2022-04-05T00:00:00Z\",\"reward\":975,\"url\":\"TestURL\"},{\"date\":\"2022-04-05T00:00:00Z\",\"reward\":975,\"url\":\"TestURL\"}],\"rewardSum\":1950,\"currency\":\"POINTS\",\"rewardsLastYear\":[{\"month\":\"4.2022\",\"reward\":1950},{\"month\":\"3.2022\",\"reward\":0},{\"month\":\"2.2022\",\"reward\":0},{\"month\":\"1.2022\",\"reward\":0},{\"month\":\"12.2021\",\"reward\":0},{\"month\":\"11.2021\",\"reward\":0},{\"month\":\"10.2021\",\"reward\":0},{\"month\":\"9.2021\",\"reward\":0},{\"month\":\"8.2021\",\"reward\":0},{\"month\":\"7.2021\",\"reward\":0},{\"month\":\"6.2021\",\"reward\":0},{\"month\":\"5.2021\",\"reward\":0}],\"timeToDisclosure\":{\"time\":[1440,1440],\"mean\":1440,\"standardDeviation\":0},\"severities\":{\"low\":2},\"meanSeverity\":2}]\n",
		},
		{
			Name:         "Invalid - Missing RedTeam",
			Owner:        "testOwner",
			RepoName:     "testRepo",
			AppInstalled: true,
			Issues: []gitlib.Issue{
				{
					HTMLURL:      "TestURL",
					CreatedAt:    open,
					ClosedAt:     &closed,
					Severities:   []gitlib.IssueSeverity{gitlib.IssueSeverity("low")},
					Migrated:     true,
					RedTeam:      nil,
					BountyPoints: pointer.Int(975),
				},
			},
			ExpectedResponse: "[]\n",
		},
		{
			Name:         "Invalid - Missing BountyPoints",
			Owner:        "testOwner",
			RepoName:     "testRepo",
			AppInstalled: true,
			Issues: []gitlib.Issue{
				{
					HTMLURL:      "TestURL",
					CreatedAt:    open,
					ClosedAt:     &closed,
					Assignees:    nil,
					Severities:   []gitlib.IssueSeverity{gitlib.IssueSeverity("low")},
					Migrated:     true,
					RedTeam:      []gitlib.User{{Login: "testUser"}},
					BountyPoints: nil,
				},
			},
			ExpectedResponse: "[]\n",
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			// GIVEN
			e := echo.New()
			b := new(bytes.Buffer)

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

			githubHandler := famed.NewHandler(nil, fakeInstallationClient, famedConfig)

			// WHEN
			err := githubHandler.GetRedTeam(ctx)

			// THEN
			assert.Equal(t, testCase.ExpectedErr, err)

			if testCase.ExpectedResponse != "" {
				assert.Equal(t, 1, fakeInstallationClient.CheckInstallationCallCount())
				assert.Equal(t, 1, fakeInstallationClient.GetIssuesByRepoCallCount())
				assert.Equal(t, testCase.ExpectedResponse, rec.Body.String())
			}
		})
	}
}
