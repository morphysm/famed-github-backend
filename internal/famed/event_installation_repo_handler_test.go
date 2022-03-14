package famed_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-github/v41/github"
	"github.com/labstack/echo/v4"
	"github.com/morphysm/famed-github-backend/internal/client/installation"
	"github.com/morphysm/famed-github-backend/internal/client/installation/installationfakes"
	"github.com/morphysm/famed-github-backend/internal/config"
	"github.com/morphysm/famed-github-backend/internal/famed"
	"github.com/morphysm/famed-github-backend/pkg/pointers"
	"github.com/stretchr/testify/assert"
)

//nolint:funlen
func TestPostInstallationRepositoriesEvent(t *testing.T) {
	t.Parallel()

	labels := map[string]installation.Label{
		config.FamedLabel:           {Name: config.FamedLabel, Color: "TestColor", Description: "TestDescription"},
		string(config.CVSSNone):     {Name: string(config.CVSSNone), Color: "TestColor", Description: "TestDescription"},
		string(config.CVSSLow):      {Name: string(config.CVSSLow), Color: "TestColor", Description: "TestDescription"},
		string(config.CVSSMedium):   {Name: string(config.CVSSMedium), Color: "TestColor", Description: "TestDescription"},
		string(config.CVSSHigh):     {Name: string(config.CVSSHigh), Color: "TestColor", Description: "TestDescription"},
		string(config.CVSSCritical): {Name: string(config.CVSSCritical), Color: "TestColor", Description: "TestDescription"},
	}
	famedConfig := famed.Config{
		Labels: labels,
	}

	testCases := []struct {
		Name        string
		Event       *github.InstallationRepositoriesEvent
		ExpectedErr error
	}{
		{
			Name:        "Empty installation event",
			Event:       &github.InstallationRepositoriesEvent{},
			ExpectedErr: famed.ErrEventMissingData,
		},
		{
			Name: "Valid",
			Event: &github.InstallationRepositoriesEvent{
				Action:            pointers.String("added"),
				RepositoriesAdded: []*github.Repository{{Name: pointers.String("TestRepo1")}},
				Installation:      &github.Installation{Account: &github.User{Login: pointers.String("TestUser")}},
			},
			ExpectedErr: nil,
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

			fakeInstallationClient := &installationfakes.FakeClient{}
			fakeInstallationClient.PostLabelReturns(nil)

			githubHandler := famed.NewHandler(nil, fakeInstallationClient, nil, famedConfig)

			// WHEN
			err = githubHandler.PostEvent(ctx)

			// THEN
			if testCase.ExpectedErr == nil {
				_, owner, repos, labels := fakeInstallationClient.PostLabelsArgsForCall(0)
				assert.Equal(t, *testCase.Event.Installation.Account.Login, owner)
				assert.Equal(t, testCase.Event.RepositoriesAdded, repos)
				assert.Equal(t, famedConfig.Labels, labels)
			}
			assert.Equal(t, testCase.ExpectedErr, err)
		})
	}
}
