package github

import (
	"context"
	"net/http"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v41/github"

	libHttp "github.com/morphysm/famed-github-backend/pkg/http"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate . AppClient
type AppClient interface {
	GetInstallations(ctx context.Context) ([]Installation, error)
	GetAccessToken(ctx context.Context, installationID int64) (*github.InstallationToken, error)
}

// githubAppClient represents a GitHub app client.
type githubAppClient struct {
	appID  int64
	client *github.Client
}

// NewAppClient returns a new instance of the GitHub client
func NewAppClient(baseURL string, apiKey string, appID int64) (AppClient, error) {
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
