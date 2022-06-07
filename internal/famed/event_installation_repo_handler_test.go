package famed_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-github/v41/github"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/morphysm/famed-github-backend/internal/config"
	"github.com/morphysm/famed-github-backend/internal/famed"
	model2 "github.com/morphysm/famed-github-backend/internal/famed/model"
	model "github.com/morphysm/famed-github-backend/internal/repositories/github/model"
	"github.com/morphysm/famed-github-backend/internal/repositories/github/providers"
	"github.com/morphysm/famed-github-backend/internal/repositories/github/providers/providersfakes"
	"github.com/morphysm/famed-github-backend/pkg/pointer"
)

//nolint:funlen
func TestPostInstallationRepositoriesEvent(t *testing.T) {
	t.Parallel()

	labels := map[string]model.Label{
		config.FamedLabelKey:   {Name: config.FamedLabelKey, Color: "TestColor", Description: "TestDescription"},
		string(model.Info):     {Name: string(model.Info), Color: "TestColor", Description: "TestDescription"},
		string(model.Low):      {Name: string(model.Low), Color: "TestColor", Description: "TestDescription"},
		string(model.Medium):   {Name: string(model.Medium), Color: "TestColor", Description: "TestDescription"},
		string(model.High):     {Name: string(model.High), Color: "TestColor", Description: "TestDescription"},
		string(model.Critical): {Name: string(model.Critical), Color: "TestColor", Description: "TestDescription"},
	}
	famedConfig := model2.Config{
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
			ExpectedErr: &echo.HTTPError{Code: 400, Message: model2.ErrEventMissingData.Error()},
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

			fakeInstallationClient := &providersfakes.FakeInstallationClient{}
			fakeInstallationClient.PostLabelReturns(nil)
			cl, _ := providers.NewInstallationClient("", nil, nil, "", "famed", nil)
			fakeInstallationClient.ValidateWebHookEventStub = cl.ValidateWebHookEvent

			githubHandler := famed.NewHandler(nil, fakeInstallationClient, famedConfig, Now)

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
