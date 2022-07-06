package providers

import (
	"context"
	"net/http"
	"strings"
	"sync"

	"github.com/google/go-github/v41/github"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"

	"github.com/morphysm/famed-github-backend/internal/repositories/github/model"
	libHttp "github.com/morphysm/famed-github-backend/pkg/http"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate . InstallationClient
type InstallationClient interface {
	GetRateLimit(ctx context.Context, owner string) (model.RateLimit, error)

	GetUser(ctx context.Context, owner string, login string) (model.User, error)

	GetReposByOwner(ctx context.Context, owner string) ([]string, error)

	GetIssuesByRepo(ctx context.Context, owner string, repoName string, options IssueListByRepoOptions) ([]model.Issue, error)
	GetEnrichedIssues(ctx context.Context, owner string, repoName string) (map[int]model.EnrichedIssue, error)
	EnrichIssues(ctx context.Context, owner string, repoName string, issues []model.Issue) map[int]model.EnrichedIssue
	EnrichIssue(ctx context.Context, owner string, repoName string, issues model.Issue) model.EnrichedIssue

	GetIssuePullRequest(ctx context.Context, owner string, repoName string, issueNumber int) (*string, error)

	GetIssueEvents(ctx context.Context, owner string, repoName string, issueNumber int) ([]model.IssueEvent, error)
	ValidateWebHookEvent(request *http.Request) (interface{}, error)

	GetComments(ctx context.Context, owner string, repoName string, issueNumber int) ([]model.IssueComment, error)
	PostComment(ctx context.Context, owner string, repoName string, issueNumber int, comment string) error
	UpdateComment(ctx context.Context, owner string, repoName string, commentID int64, comment string) error
	DeleteComment(ctx context.Context, owner string, repoName string, commentID int64) error

	PostLabel(ctx context.Context, owner string, repoName string, label model.Label) error
	PostLabels(ctx context.Context, owner string, repoNames []string, labels map[string]model.Label) []error

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

type safeUserMap struct {
	sync.RWMutex
	wrappedUsers map[string]model.User
}

func newSafeUserMap() *safeUserMap {
	return &safeUserMap{
		wrappedUsers: make(map[string]model.User),
	}
}

func (s *safeUserMap) Add(user model.User) {
	s.Lock()
	defer s.Unlock()
	s.wrappedUsers[user.Login] = user
}

func (s *safeUserMap) Get(login string) (model.User, bool) {
	s.Lock()
	defer s.Unlock()
	user, ok := s.wrappedUsers[login]
	return user, ok
}

// githubInstallationClient represents all GitHub github clients
type githubInstallationClient struct {
	baseURL       string
	webhookSecret string
	appClient     AppClient
	clients       safeClientMap
	famedLabel    string
	// TODO replace by cache eg. redis
	redTeamLogins map[string]string
	cachedRedTeam *safeUserMap
}

// NewInstallationClient returns a new instance of the GitHub client
func NewInstallationClient(baseURL string, appClient AppClient, installations map[string]int64, webhookSecret string, famedLabel string, redTeamLogins map[string]string) (InstallationClient, error) {
	client := &githubInstallationClient{
		baseURL:       baseURL,
		webhookSecret: webhookSecret,
		appClient:     appClient,
		clients:       newSafeClientMap(),
		famedLabel:    famedLabel,
		redTeamLogins: redTeamLogins,
		cachedRedTeam: newSafeUserMap(),
	}

	for owner, installationID := range installations {
		err := client.AddInstallation(owner, installationID)
		if err != nil {
			return nil, err
		}
	}

	return client, nil
}

// AddInstallation adds a new GitHub to the githubInstallationClient.
func (c *githubInstallationClient) AddInstallation(owner string, installationID int64) error {
	ts := NewGithubTokenSource(c.appClient, installationID)
	oAuthClient := oauth2.NewClient(context.Background(), ts)
	loggingClient := libHttp.AddLogging(oAuthClient)

	client, err := github.NewEnterpriseClient(c.baseURL, c.baseURL, loggingClient)
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
