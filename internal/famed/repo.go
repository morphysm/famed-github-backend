package famed

import (
	"context"

	"github.com/google/go-github/v41/github"
	"github.com/labstack/echo/v4"
	"github.com/morphysm/famed-github-backend/internal/client/installation"
	"github.com/morphysm/famed-github-backend/internal/config"
)

type Repo interface {
	Contributors(ctx context.Context) ([]*Contributor, error)

	RewardComment(ctx context.Context, issue *github.Issue) (string, error)
	RewardComments(ctx context.Context) (map[int]string, error)
}

type repo struct {
	config             Config
	installationClient installation.Client
	owner              string
	name               string
	issues             map[int]Issue
	contributors       Contributors
}

type Config struct {
	Currency string
	Rewards  map[config.IssueSeverity]float64
	Labels   map[string]installation.Label
	BotLogin string
}

// NewRepo returns a new instance of the famed repo representation.
func NewRepo(config Config, installationClient installation.Client, owner string, name string) Repo {
	return &repo{
		config:             config,
		installationClient: installationClient,
		owner:              owner,
		name:               name,
	}
}

func (r *repo) Contributors(ctx context.Context) ([]*Contributor, error) {
	err := r.loadIssuesAndEvents(ctx)
	if err != nil {
		return nil, err
	}

	if len(r.issues) == 0 {
		return []*Contributor{}, nil
	}

	// Use issues to generate contributor list
	contributors := r.contributorsArray()

	return contributors, nil
}

func (r *repo) RewardComment(ctx context.Context, issue *github.Issue) (string, error) {
	err := r.loadEventsForIssue(ctx, issue)
	if err != nil {
		return "", err
	}

	r.ContributorsForIssues()
	comment := RewardComment(r.issues[*issue.Number], r.contributors, r.config.Currency)

	return comment, nil
}

func (r *repo) RewardComments(ctx context.Context) (map[int]string, error) {
	err := r.loadIssuesAndEvents(ctx)
	if err != nil {
		return nil, err
	}

	if len(r.issues) == 0 {
		return map[int]string{}, nil
	}

	comments := make(map[int]string, len(r.issues))
	for issueNumber := range r.issues {
		r.ContributorsForIssue(issueNumber)

		comments[issueNumber] = RewardComment(r.issues[issueNumber], r.contributors, r.config.Currency)
	}

	return comments, nil
}

func (r *repo) loadEventsForIssue(ctx context.Context, issue *github.Issue) error {
	r.issues = make(map[int]Issue, 1)
	r.issues[*issue.Number] = Issue{Issue: issue}

	return r.getEvents(ctx)
}

func (r *repo) loadIssuesAndEvents(ctx context.Context) error {
	// Get all issues filtered by label and closed state
	famedLabel := r.config.Labels[config.FamedLabel]
	issuesResponse, err := r.installationClient.GetIssuesByRepo(ctx, r.owner, r.name, []string{famedLabel.Name}, installation.Closed)
	if err != nil {
		return echo.ErrBadGateway.SetInternal(err)
	}

	filteredIssues := filterIssues(issuesResponse)
	if len(filteredIssues) == 0 {
		return nil
	}

	r.issues = make(map[int]Issue, len(filteredIssues))
	for _, issue := range filteredIssues {
		r.issues[*issue.Number] = Issue{Issue: issue}
	}

	return r.getEvents(ctx)
}
