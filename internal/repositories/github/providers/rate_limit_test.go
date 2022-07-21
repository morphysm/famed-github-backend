package providers_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/google/go-github/v41/github"
	"github.com/stretchr/testify/assert"

	"github.com/morphysm/famed-github-backend/internal/repositories/github/model"
	"github.com/morphysm/famed-github-backend/internal/repositories/github/providers"
	"github.com/morphysm/famed-github-backend/internal/repositories/github/providers/providersfakes"
)

var testTime = time.Date(2022, 4, 20, 0, 0, 0, 0, time.Local)

func TestGetRateLimits(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                 string
		Owner                string
		GitHubResponseBody   string
		GitHubResponseStatus int
		GitHubResponseError  error
		ExpectedResponse     model.RateLimits
		ExpectErr            bool
	}{
		{
			Name:  "Valid",
			Owner: "testOwner",
			GitHubResponseBody: fmt.Sprintf(`{"resources":{
					"core": {"limit":2,"remaining":1,"reset":%d},
					"search": {"limit":3,"remaining":2,"reset":%d}
				}}`, testTime.Unix(), testTime.Unix()),
			ExpectedResponse: model.RateLimits{
				Core: model.Rate{
					Limit:     2,
					Remaining: 1,
					Reset:     testTime},
				Search: model.Rate{
					Limit:     3,
					Remaining: 2,
					Reset:     testTime}},
		},
		{
			Name:  "Missing Search",
			Owner: "testOwner",
			GitHubResponseBody: fmt.Sprintf(`{"resources":{
					"core": {"limit":2,"remaining":1,"reset":%d}
				}}`, testTime.Unix()),
			ExpectErr: true,
		},
		{
			Name:  "Missing Core",
			Owner: "testOwner",
			GitHubResponseBody: fmt.Sprintf(`{"resources":{
					"search": {"limit":2,"remaining":1,"reset":%d}
				}}`, testTime.Unix()),
			ExpectErr: true,
		},
		{
			Name:                 "GitHub Error",
			Owner:                "testOwner",
			GitHubResponseStatus: http.StatusInternalServerError,
			GitHubResponseError:  errors.New("GitHub Error"),
			ExpectErr:            true,
		},
		{
			Name:      "No Owner",
			Owner:     "",
			ExpectErr: true,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()

			// GIVEN
			fakeGitHubServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if testCase.GitHubResponseBody != "" {
					fmt.Fprint(w, testCase.GitHubResponseBody)
				}
				if testCase.GitHubResponseError != nil {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprint(w, testCase.GitHubResponseError.Error())
				}
			}))
			defer fakeGitHubServer.Close()

			fakeAppClient := &providersfakes.FakeAppClient{}
			client := fakeGitHubServer.Client()
			fakeGitHubClient, err := github.NewEnterpriseClient("", "", client)
			fakeGitHubClient.BaseURL, _ = url.Parse(fakeGitHubServer.URL + "/")
			assert.NoError(t, err)

			githubInstallationClient, err := providers.NewInstallationClient("", fakeAppClient, nil, "", "", nil)
			assert.NoError(t, err)
			githubInstallationClient.AddGitHubClient("testOwner", fakeGitHubClient)

			// WHEN
			rateLimits, err := githubInstallationClient.GetRateLimits(context.Background(), testCase.Owner)

			// THEN
			if testCase.ExpectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if testCase.ExpectedResponse != (model.RateLimits{}) {
				assert.Equal(t, testCase.ExpectedResponse, rateLimits)
			}
		})
	}
}
