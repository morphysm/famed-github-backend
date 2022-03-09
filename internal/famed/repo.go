package famed

import (
	"context"
	"fmt"

	"github.com/google/go-github/v41/github"
	"github.com/labstack/echo/v4"
	"github.com/morphysm/famed-github-backend/internal/client/installation"
	"github.com/morphysm/famed-github-backend/internal/config"
)

type Repo interface {
	Contributors(ctx context.Context) ([]*Contributor, error)

	ContributorComment(ctx context.Context, issue *github.Issue) (string, error)
	ContributorComments(ctx context.Context) (map[int]string, error)

	IssueStateComment(ctx context.Context, issue *github.Issue) (string, error)
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

func (r *repo) ContributorComment(ctx context.Context, issue *github.Issue) (string, error) {
	err := r.loadRateAndEventsForIssue(ctx, issue)
	if err != nil {
		return "", err
	}

	r.ContributorsForIssues()
	comment := r.comment(*issue.Number)

	return comment, nil
}

func (r *repo) ContributorComments(ctx context.Context) (map[int]string, error) {
	err := r.loadIssuesRateAndEvents(ctx)
	if err != nil {
		return nil, err
	}

	if len(r.issues) == 0 {
		return map[int]string{}, nil
	}

	comments := make(map[int]string, len(r.issues))
	for issueNumber := range r.issues {
		r.ContributorsForIssue(issueNumber)

		comments[issueNumber] = r.comment(issueNumber)
	}

	return comments, nil
}

func (r *repo) loadRateAndEventsForIssue(ctx context.Context, issue *github.Issue) error {
	r.issues = make(map[int]Issue, 1)
	r.issues[*issue.Number] = Issue{Issue: issue}

	return r.loadRateAndEvents(ctx)
}

func (r *repo) loadIssuesRateAndEvents(ctx context.Context) error {
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

	return r.loadRateAndEvents(ctx)
}

func (r *repo) loadRateAndEvents(ctx context.Context) error {
	// Get all events for each issue
	err := r.getEvents(ctx, r.owner, r.name)
	if err != nil {
		return echo.ErrBadGateway.SetInternal(err)
	}

	return nil
}

func (r *repo) IssueStateComment(ctx context.Context, issue *github.Issue) (string, error) {
	comment := fmt.Sprintf("🤖 Assignees for Issue **%s #%d** are now eligible to Get Famed.", *issue.Title, *issue.Number)

	// Check that an assignee is assigned
	comment = fmt.Sprintf("%s\n%s️", comment, assigneeComment(issue))

	// Check that a valid severity label is assigned
	comment = fmt.Sprintf("%s\n%s️", comment, severityComment(Issue{Issue: issue}))

	// Check that a PR is assigned
	comment = fmt.Sprintf("%s\n%s", comment, prComment(issue))

	// Final note
	comment = fmt.Sprintf("%s\n\nHappy hacking! 🦾💙❤️️", comment)

	return comment, nil
}

func assigneeComment(issue *github.Issue) string {
	if issue.Assignee != nil {
		return "- [x] Add assignees to track contribution times of the issue \U0001F9B8‍♀️\U0001F9B9"
	}

	return "- [ ] Add assignees to track contribution times of the issue \U0001F9B8‍♀️\U0001F9B9"
}

func severityComment(issue Issue) string {
	_, err := issue.severity()
	if err == nil {
		return "- [x] Add a severity (CVSS) label to compute the score 🏷️"
	}

	return "- [ ] Add a severity (CVSS) label to compute the score 🏷️"
}

func prComment(issue *github.Issue) string {
	if issue.PullRequestLinks != nil {
		return "- [x] Link a PR when closing the issue ♻️ \U0001F9B8‍♀️\U0001F9B9"
	}

	return "- [ ] Link a PR when closing the issue ♻️ \U0001F9B8‍♀️\U0001F9B9"
}
