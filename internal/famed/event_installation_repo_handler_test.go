package famed_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-github/v41/github"
	"github.com/labstack/echo/v4"
	gitLib "github.com/morphysm/famed-github-backend/internal/client/github"
	"github.com/morphysm/famed-github-backend/internal/client/github/githubfakes"
	"github.com/morphysm/famed-github-backend/internal/config"
	"github.com/morphysm/famed-github-backend/internal/famed"
	"github.com/morphysm/famed-github-backend/pkg/pointer"
	"github.com/stretchr/testify/assert"
)

//nolint:funlen
func TestPostInstallationRepositoriesEvent(t *testing.T) {
	t.Parallel()

	labels := map[string]gitLib.Label{
		config.FamedLabelKey:    {Name: config.FamedLabelKey, Color: "TestColor", Description: "TestDescription"},
		string(gitLib.Info):     {Name: string(gitLib.Info), Color: "TestColor", Description: "TestDescription"},
		string(gitLib.Low):      {Name: string(gitLib.Low), Color: "TestColor", Description: "TestDescription"},
		string(gitLib.Medium):   {Name: string(gitLib.Medium), Color: "TestColor", Description: "TestDescription"},
		string(gitLib.High):     {Name: string(gitLib.High), Color: "TestColor", Description: "TestDescription"},
		string(gitLib.Critical): {Name: string(gitLib.Critical), Color: "TestColor", Description: "TestDescription"},
	}
	famedConfig := famed.Config{
		Labels: labels,
	}

	testCases := []struct {
		Name          string
		Event         *github.InstallationRepositoriesEvent
		ExpectedRepos []string
		ExpectedErr   *echo.HTTPError
	}{
		{
			Name:        "Empty github event",
			Event:       &github.InstallationRepositoriesEvent{},
			ExpectedErr: &echo.HTTPError{Code: 400, Message: famed.ErrEventMissingData.Error()},
		},
		{
			Name: "Valid",
			Event: &github.InstallationRepositoriesEvent{
				Action:            pointer.String("added"),
				RepositoriesAdded: []*github.Repository{{Name: pointer.String("TestRepo1")}},
				Installation:      &github.Installation{Account: &github.User{Login: pointer.String("TestUser")}},
			},
			ExpectedRepos: []string{"TestRepo1"},
			ExpectedErr:   nil,
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
			req.Header.Set(github.EventTypeHeader, "installation_repositories")
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			fakeInstallationClient := &githubfakes.FakeInstallationClient{}
			fakeInstallationClient.PostLabelReturns(nil)
			cl, _ := gitLib.NewInstallationClient("", nil, nil, "", "famed", nil)
			fakeInstallationClient.ValidateWebHookEventStub = cl.ValidateWebHookEvent

			githubHandler := famed.NewHandler(nil, fakeInstallationClient, famedConfig)

			// WHEN
			err = githubHandler.PostEvent(ctx)

			// THEN
			if testCase.ExpectedErr == nil {
				assert.Equal(t, 1, fakeInstallationClient.PostLabelsCallCount())
				if fakeInstallationClient.PostLabelsCallCount() == 1 {
					_, owner, repos, labels := fakeInstallationClient.PostLabelsArgsForCall(0)
					assert.Equal(t, *testCase.Event.Installation.Account.Login, owner)
					assert.Equal(t, testCase.ExpectedRepos, repos)
					assert.Equal(t, famedConfig.Labels, labels)
				}
			} else {
				assert.Equal(t, testCase.ExpectedErr, err)
			}
		})
	}
}
