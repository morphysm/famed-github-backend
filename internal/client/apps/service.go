package apps

import (
	"context"
	"net/http"
	"time"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v41/github"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate . Client
type Client interface {
	GetAccessTokens(ctx context.Context, installationID int64, repositoryIDs []int64) (*github.InstallationToken, error)
}

type githubAppClient struct {
	appID  int64
	client *github.Client
}

// NewClient returns a new instance of the Github client
func NewClient(baseURL string, apiKey string, appID int64) (Client, error) {
	itr, err := ghinstallation.NewAppsTransport(http.DefaultTransport, appID, []byte(apiKey))
	if err != nil {
		return nil, err
	}

	itr.BaseURL = baseURL

	// Create git client with apps transport
	client, err := github.NewEnterpriseClient(
		baseURL,
		baseURL,
		&http.Client{
			Transport: itr,
			Timeout:   time.Second * 30,
		})
	if err != nil {
		return nil, err
	}

	return &githubAppClient{
		appID:  appID,
		client: client,
	}, nil
}
