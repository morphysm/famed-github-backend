package installation

import (
	"context"
	"strings"

	"github.com/google/go-github/v41/github"
	"github.com/morphysm/famed-github-backend/internal/client/app"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate . Client
type Client interface {
	GetIssuesByRepo(ctx context.Context, owner string, repoName string, labels []string, state IssueState) ([]*github.Issue, error)
	ReopenIssue(ctx context.Context, owner string, repoName string, issueNumber int) error

	GetIssuePullRequest(ctx context.Context, owner string, repoName string, issueNumber int) (*PullRequest, error)

	GetIssueEvents(ctx context.Context, owner string, repoName string, issueNumber int) ([]*github.IssueEvent, error)
	GetIssuesEvents(ctx context.Context, owner string, repoName string, issues []*github.Issue) (map[int][]*github.IssueEvent, map[int]error)

	GetComments(ctx context.Context, owner string, repoName string, issueNumber int) ([]*github.IssueComment, error)
	PostComment(ctx context.Context, owner string, repoName string, issueNumber int, comment string) error
	UpdateComment(ctx context.Context, owner string, repoName string, commentID int64, comment string) error

	PostLabel(ctx context.Context, owner string, repo string, label Label) error
	PostLabels(ctx context.Context, owner string, repositories []*github.Repository, labels map[string]Label) []error

	AddInstallation(owner string, installationID int64) error
	CheckInstallation(owner string) bool
}

// safeClientMap represents a map from owner to client.
// The map is wrapped to avoid any capitalization errors.
type safeClientMap struct {
	m   map[string]*github.Client
	qlM map[string]*githubv4.Client
}

// newSafeClientMap returns a new safeClientMap.
func newSafeClientMap() safeClientMap {
	return safeClientMap{
		m:   make(map[string]*github.Client),
		qlM: make(map[string]*githubv4.Client),
	}
}

// add adds an owner client pair to the safeClientMap.
func (s safeClientMap) add(owner string, client *github.Client) {
	s.m[strings.ToLower(owner)] = client
}

// get gets an owner client pair from the safeClientMap.
func (s safeClientMap) get(owner string) (*github.Client, bool) {
	client, ok := s.m[strings.ToLower(owner)]
	return client, ok
}

// add adds an GraphQL owner client pair to the safeClientMap.
func (s safeClientMap) addGql(owner string, client *githubv4.Client) {
	s.qlM[strings.ToLower(owner)] = client
}

// get gets an GraphQL owner client pair from the safeClientMap.
func (s safeClientMap) getGql(owner string) (*githubv4.Client, bool) {
	client, ok := s.qlM[strings.ToLower(owner)]
	return client, ok
}

// githubInstallationClient represents all GitHub installation clients
type githubInstallationClient struct {
	baseURL   string
	appClient app.Client
	clients   safeClientMap
}

// NewClient returns a new instance of the GitHub client
func NewClient(baseURL string, appClient app.Client, installations map[string]int64) (Client, error) {
	client := &githubInstallationClient{
		baseURL:   baseURL,
		appClient: appClient,
		clients:   newSafeClientMap(),
	}

	for owner, installationID := range installations {
		err := client.AddInstallation(owner, installationID)
		if err != nil {
			return nil, err
		}
	}

	return client, nil
}

// AddInstallation adds a new installation to the githubInstallationClient.
func (c *githubInstallationClient) AddInstallation(owner string, installationID int64) error {
	ts := NewGithubTokenSource(c.appClient, installationID)
	oAuthClient := oauth2.NewClient(context.Background(), ts)

	client, err := github.NewEnterpriseClient(c.baseURL, c.baseURL, oAuthClient)
	if err != nil {
		return err
	}

	c.clients.add(owner, client)

	// GraphQL client for missing "pull_requests" field workaround https://github.community/t/get-referenced-pull-request-from-issue/14027
	gQLClient := githubv4.NewClient(oAuthClient)
	c.clients.addGql(owner, gQLClient)

	return nil
}

// CheckInstallation checks if an installations is present in the githubInstallationClient.
func (c *githubInstallationClient) CheckInstallation(owner string) bool {
	_, ok := c.clients.get(owner)
	return ok
}
