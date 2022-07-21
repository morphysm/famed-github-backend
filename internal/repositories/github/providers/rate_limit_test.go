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
			GitHubResponseBody: `{"resources":{
					"core": {"limit":2,"remaining":1,"reset":1372700873},
					"search": {"limit":3,"remaining":2,"reset":1372700874},
					"graphql": {"limit":4,"remaining":3,"reset":1372700875},
					"integration_manifest": {"limit":5,"remaining":4,"reset":1372700876},
					"source_import": {"limit":6,"remaining":5,"reset":1372700877},
					"code_scanning_upload": {"limit":7,"remaining":6,"reset":1372700878},
					"actions_runner_registration": {"limit":8,"remaining":7,"reset":1372700879},
					"scim": {"limit":9,"remaining":8,"reset":1372700880}
				}}`,
			ExpectedResponse: model.RateLimits{
				Core: model.Rate{
					Limit:     2,
					Remaining: 1,
					Reset:     time.Date(2013, time.July, 1, 19, 47, 53, 0, time.Local)},
				Search: model.Rate{
					Limit:     3,
					Remaining: 2,
					Reset:     time.Date(2013, time.July, 1, 19, 47, 54, 0, time.Local)}},
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
