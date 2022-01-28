package installation

import (
	"context"

	"github.com/google/go-github/v41/github"
	"golang.org/x/oauth2"

	"github.com/morphysm/kudos-github-backend/internal/client/apps"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate . Client
type Client interface {
	GetIssuesByRepo(ctx context.Context, repoName string, labels []string, state IssueState) ([]*github.Issue, error)
	GetIssueEvents(ctx context.Context, repoName string, issueNumber int) ([]*github.IssueEvent, error)
	PostComment(ctx context.Context, repoName string, issueNumber int, comment string) (*github.IssueComment, error)
}

type githubInstallationClient struct {
	baseURL string
	owner   string
	client  *github.Client
}

// NewClient returns a new instance of the GitHub client
func NewClient(baseURL string, client apps.Client, installationID int64, repoIDs []int64, owner string) (Client, error) {
	ts := NewGithubTokenSource(client, installationID, repoIDs)
	oAuthClient := oauth2.NewClient(context.Background(), ts)

	apiClient, err := github.NewEnterpriseClient(baseURL, baseURL, oAuthClient)
	if err != nil {
		return nil, err
	}

	return &githubInstallationClient{
		baseURL: baseURL,
		owner:   owner,
		client:  apiClient,
	}, nil
}
