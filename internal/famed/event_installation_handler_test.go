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
	"github.com/morphysm/famed-github-backend/internal/famed"
	"github.com/morphysm/famed-github-backend/pkg/pointers"
	"github.com/stretchr/testify/assert"
)

func TestPostInstallationEvent(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name        string
		Event       *github.InstallationEvent
		ExpectedErr error
	}{
		{
			Name:        "Empty installation repository event",
			Event:       &github.InstallationEvent{},
			ExpectedErr: famed.ErrEventMissingData,
		},
		{
			Name: "Valid",
			Event: &github.InstallationEvent{
				Action:       pointers.String("created"),
				Installation: &github.Installation{ID: pointers.Int64(0), Account: &github.User{Login: pointers.String("TestUser")}},
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
			req.Header.Set(github.EventTypeHeader, "installation")
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			fakeInstallationClient := &installationfakes.FakeClient{}
			fakeInstallationClient.AddInstallationReturns(nil)
			cl, _ := installation.NewClient("", nil, nil, "")
			fakeInstallationClient.ValidateWebHookEventStub = cl.ValidateWebHookEvent

			githubHandler := famed.NewHandler(nil, fakeInstallationClient, NewTestConfig())

			// WHEN
			err = githubHandler.PostEvent(ctx)

			// THEN
			if testCase.ExpectedErr == nil {
				assert.Equal(t, 1, fakeInstallationClient.AddInstallationCallCount())
				if fakeInstallationClient.AddInstallationCallCount() == 1 {
					owner, installationID := fakeInstallationClient.AddInstallationArgsForCall(0)
					assert.Equal(t, *testCase.Event.Installation.Account.Login, owner)
					assert.Equal(t, *testCase.Event.Installation.ID, installationID)
				}
			}
			assert.Equal(t, testCase.ExpectedErr, err)
		})
	}
}
