package famed

import (
	"context"

	"github.com/google/go-github/v41/github"
	"github.com/labstack/echo/v4"
	"github.com/morphysm/famed-github-backend/internal/client/currency"
	"github.com/morphysm/famed-github-backend/internal/client/installation"
)

type Repo interface {
	GetContributors(ctx context.Context) ([]*Contributor, error)
	GetComment(ctx context.Context, issue *github.Issue) (string, error)
	GetComments(ctx context.Context) (map[int64]string, error)
}

type repo struct {
	config             Config
	installationClient installation.Client
	currencyClient     currency.Client
	name               string
	issues             map[int64]Issue
	ethRate            float64
	contributors       Contributors
}

// NewRepo returns a new instance of the famed repo representation.
func NewRepo(config Config, installationClient installation.Client, currencyClient currency.Client, name string) Repo {
	return &repo{
		config:             config,
		installationClient: installationClient,
		currencyClient:     currencyClient,
		name:               name,
	}
}

func (r *repo) GetContributors(ctx context.Context) ([]*Contributor, error) {
	err := r.loadIssuesRateAndEvents(ctx)
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

func (r *repo) GetComment(ctx context.Context, issue *github.Issue) (string, error) {
	err := r.loadRateAndEventsForIssue(ctx, issue)
	if err != nil {
		return "", err
	}

	r.Contributors()
	comment := r.comment(*issue.ID)

	return comment, nil
}

func (r *repo) GetComments(ctx context.Context) (map[int64]string, error) {
	err := r.loadIssuesRateAndEvents(ctx)
	if err != nil {
		return nil, err
	}

	if len(r.issues) == 0 {
		return map[int64]string{}, nil
	}

	comments := make(map[int64]string, len(r.issues))
	for issueID := range r.issues {
		r.ContributorsForIssue(issueID)
		comments[issueID] = r.comment(issueID)
	}

	return comments, nil
}

func (r *repo) loadRateAndEventsForIssue(ctx context.Context, issue *github.Issue) error {
	r.issues = make(map[int64]Issue, 1)
	r.issues[*issue.ID] = Issue{Issue: issue}

	return r.loadRateAndEvents(ctx)
}

func (r *repo) loadIssuesRateAndEvents(ctx context.Context) error {
	// Get all issues filtered by label and closed state
	issuesResponse, err := r.installationClient.GetIssuesByRepo(ctx, r.name, []string{r.config.Label}, installation.Closed)
	if err != nil {
		return echo.ErrBadGateway.SetInternal(err)
	}

	filteredIssues := filterIssues(issuesResponse)
	if len(filteredIssues) == 0 {
		return nil
	}

	r.issues = make(map[int64]Issue, len(filteredIssues))
	for _, issue := range filteredIssues {
		r.issues[*issue.ID] = Issue{Issue: issue}
	}

	return r.loadRateAndEvents(ctx)
}

func (r *repo) loadRateAndEvents(ctx context.Context) error {
	ethRate, err := r.currencyClient.GetUSDToETHConversion(ctx)
	if err != nil {
		return echo.ErrBadGateway.SetInternal(err)
	}

	r.ethRate = ethRate

	// Get all events for each issue
	err = r.getEvents(ctx, r.name)
	if err != nil {
		return echo.ErrBadGateway.SetInternal(err)
	}

	return nil
}
