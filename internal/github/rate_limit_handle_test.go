package github_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/morphysm/famed-github-backend/internal/github"
	"github.com/morphysm/famed-github-backend/internal/repositories/github/model"
	"github.com/morphysm/famed-github-backend/internal/repositories/github/providers/providersfakes"
)

var testTime = time.Date(2022, 4, 20, 0, 0, 0, 0, time.UTC)

func TestGetRateLimits(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name             string
		Owner            string
		RateLimits       model.RateLimits
		RateLimitsError  error
		ExpectedResponse string
		ExpectedErr      error
	}{
		{
			Name:  "Valid",
			Owner: "testOwner",
			RateLimits: model.RateLimits{
				Core: model.Rate{
					Limit:     1,
					Remaining: 1,
					Reset:     testTime,
				},
				Search: model.Rate{
					Limit:     1,
					Remaining: 1,
					Reset:     testTime,
				},
			},
			ExpectedResponse: "{\"core\":{\"limit\":1,\"remaining\":1,\"reset\":\"2022-04-20T00:00:00Z\"},\"search\":{\"limit\":1,\"remaining\":1,\"reset\":\"2022-04-20T00:00:00Z\"}}\n",
		},
		{
			Name:  "No Owner",
			Owner: "",
			RateLimits: model.RateLimits{
				Core: model.Rate{
					Limit:     1,
					Remaining: 1,
					Reset:     testTime,
				},
				Search: model.Rate{
					Limit:     1,
					Remaining: 1,
					Reset:     testTime,
				},
			},
			ExpectedErr: &echo.HTTPError{
				Code:     400,
				Message:  "missing owner path parameter",
				Internal: nil,
			},
			ExpectedResponse: "",
		},
		{
			Name:             "Bad Response",
			Owner:            "testOwner",
			RateLimitsError:  errors.New("testError"),
			ExpectedErr:      errors.New("testError"),
			ExpectedResponse: "",
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			// GIVEN
			e := echo.New()

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/github/ratelimits/%s", testCase.Owner), nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.SetParamNames([]string{"owner"}...)
			ctx.SetParamValues([]string{testCase.Owner}...)

			fakeInstallationClient := &providersfakes.FakeInstallationClient{}
			fakeInstallationClient.GetRateLimitsReturns(testCase.RateLimits, testCase.RateLimitsError)

			githubHandler := github.NewHandler(fakeInstallationClient)

			// WHEN
			err := githubHandler.GetRateLimits(ctx)

			// THEN
			assert.Equal(t, testCase.ExpectedErr, err)

			if testCase.ExpectedResponse != "" {
				assert.Equal(t, 1, fakeInstallationClient.GetRateLimitsCallCount())
				assert.Equal(t, testCase.ExpectedResponse, rec.Body.String())
			}
		})
	}
}
