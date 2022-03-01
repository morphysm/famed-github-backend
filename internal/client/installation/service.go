package installation

import (
	"context"

	"github.com/google/go-github/v41/github"
	"github.com/morphysm/famed-github-backend/internal/client/apps"
	"golang.org/x/oauth2"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate . Client
type Client interface {
	GetIssuesByRepo(ctx context.Context, owner string, repoName string, labels []string, state IssueState) ([]*github.Issue, error)
	GetIssueEvents(ctx context.Context, owner string, repoName string, issueNumber int) ([]*github.IssueEvent, error)
	GetComments(ctx context.Context, owner string, repoName string, issueNumber int) ([]*github.IssueComment, error)
	PostComment(ctx context.Context, owner string, repoName string, issueNumber int, comment string) error
	PostLabel(ctx context.Context, owner string, repo string, label Label) error
}

type githubInstallationClient struct {
	baseURL       string
	installations map[string]int64
	clients       map[string]*github.Client
}

// NewClient returns a new instance of the GitHub client
func NewClient(baseURL string, client apps.Client, installations map[string]int64) (Client, error) {
	clients := make(map[string]*github.Client)

	for owner, installationID := range installations {
		ts := NewGithubTokenSource(client, installationID)
		oAuthClient := oauth2.NewClient(context.Background(), ts)

		installationClient, err := github.NewEnterpriseClient(baseURL, baseURL, oAuthClient)
		if err != nil {
			return nil, err
		}

		clients[owner] = installationClient
	}

	return &githubInstallationClient{
		baseURL:       baseURL,
		installations: installations,
		clients:       clients,
	}, nil
}
