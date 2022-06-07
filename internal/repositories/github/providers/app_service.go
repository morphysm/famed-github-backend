package providers

import (
	"context"
	"net/http"

	"github.com/awnumar/memguard"
	"github.com/google/go-github/v41/github"

	"github.com/morphysm/famed-github-backend/internal/repositories/github/model"
	libHttp "github.com/morphysm/famed-github-backend/pkg/http"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate . AppClient
type AppClient interface {
	GetInstallations(ctx context.Context) ([]model.Installation, error)
	GetAccessToken(ctx context.Context, installationID int64) (*github.InstallationToken, error)
}

// githubAppClient represents a GitHub app client.
type githubAppClient struct {
	appID  int64
	client *github.Client
}

// NewAppClient returns a new instance of the GitHub client
func NewAppClient(baseURL string, appID int64, keyEnclave *memguard.Enclave) (AppClient, error) {
	transport := NewAppsTransport(baseURL, http.DefaultTransport, appID, keyEnclave)
	loggingClient := libHttp.AddLogging(&http.Client{
		Transport: transport,
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
