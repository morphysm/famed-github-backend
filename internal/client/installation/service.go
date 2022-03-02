package installation

import (
	"context"

	"github.com/google/go-github/v41/github"
	"github.com/morphysm/famed-github-backend/internal/client/app"
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

	AddInstallation(owner string, installationID int64) error
	CheckInstallation(owner string) bool
}

type githubInstallationClient struct {
	baseURL   string
	appClient app.Client
	clients   map[string]*github.Client
}

// NewClient returns a new instance of the GitHub client
func NewClient(baseURL string, appClient app.Client, installations map[string]int64) (Client, error) {
	clients := make(map[string]*github.Client)

	for owner, installationID := range installations {
		ts := NewGithubTokenSource(appClient, installationID)
		oAuthClient := oauth2.NewClient(context.Background(), ts)

		installationClient, err := github.NewEnterpriseClient(baseURL, baseURL, oAuthClient)
		if err != nil {
			return nil, err
		}

		clients[owner] = installationClient
	}

	return &githubInstallationClient{
		baseURL:   baseURL,
		appClient: appClient,
		clients:   clients,
	}, nil
}

func (c *githubInstallationClient) AddInstallation(owner string, installationID int64) error {
	ts := NewGithubTokenSource(c.appClient, installationID)
	oAuthClient := oauth2.NewClient(context.Background(), ts)

	client, err := github.NewEnterpriseClient(c.baseURL, c.baseURL, oAuthClient)
	if err != nil {
		return err
	}

	c.clients[owner] = client
	return nil
}

func (c *githubInstallationClient) CheckInstallation(owner string) bool {
	_, ok := c.clients[owner]
	return ok
}
