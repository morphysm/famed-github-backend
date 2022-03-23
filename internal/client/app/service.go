package app

import (
	"context"
	"net/http"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v41/github"
	libHttp "github.com/morphysm/famed-github-backend/pkg/http"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate . Client
type Client interface {
	GetInstallations(ctx context.Context) ([]*github.Installation, error)
	GetAccessToken(ctx context.Context, installationID int64) (*github.InstallationToken, error)
	GetRateLimits(ctx context.Context) (*github.RateLimits, error)
}

type githubAppClient struct {
	appID  int64
	client *github.Client
}

// NewClient returns a new instance of the GitHub client
func NewClient(baseURL string, apiKey string, appID int64) (Client, error) {
	itr, err := ghinstallation.NewAppsTransport(http.DefaultTransport, appID, []byte(apiKey))
	if err != nil {
		return nil, err
	}

	itr.BaseURL = baseURL
	loggingClient := libHttp.AddLogging(&http.Client{
		Transport: itr,
	})

	// Create git client with app transport
	client, err := github.NewEnterpriseClient(
		baseURL,
		baseURL,
		loggingClient)
	if err != nil {
		return nil, err
	}

	return &githubAppClient{
		appID:  appID,
		client: client,
	}, nil
}
