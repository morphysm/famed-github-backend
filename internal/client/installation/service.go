package installation

import (
	"context"

	"github.com/google/go-github/v41/github"
	"golang.org/x/oauth2"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate . Client
type Client interface {
	GetRepos(ctx context.Context) ([]*github.Repository, error)
	GetRepoLabels(ctx context.Context, repoID string) ([]*github.Label, error)
	GetRepoEvents(ctx context.Context, repoID string) ([]*github.Event, error)

	GetIssuesByRepo(ctx context.Context, repoName string, labels []string, state IssueState) ([]*github.Issue, error)
	GetIssueEvents(ctx context.Context, repoName string, issueNumber int) ([]*github.IssueEvent, error)
	PostComment(ctx context.Context, repoName string, issueNumber int, comment string) (*github.IssueComment, error)
}

type githubInstallationClient struct {
	baseURL        string
	installationID int
	owner          string
	client         *github.Client
}

// NewClient returns a new instance of the Github client
func NewClient(baseURL string, token *github.InstallationToken, owner string, installationID int) (Client, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token.GetToken()},
	)
	oAuthClient := oauth2.NewClient(context.Background(), ts)

	apiClient, err := github.NewEnterpriseClient(baseURL, baseURL, oAuthClient)
	if err != nil {
		return nil, err
	}

	return &githubInstallationClient{
		baseURL:        baseURL,
		installationID: installationID,
		owner:          owner,
		client:         apiClient,
	}, nil
}
