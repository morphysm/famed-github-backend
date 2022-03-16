package famed

import (
	"context"
	"log"

	"github.com/google/go-github/v41/github"
	"github.com/labstack/echo/v4"
	"github.com/morphysm/famed-github-backend/internal/client/installation"
	"github.com/morphysm/famed-github-backend/internal/config"
)

type Repo interface {
	Contributors(ctx context.Context) ([]*Contributor, error)

	RewardComment(ctx context.Context, issue *github.Issue) (string, error)
	RewardComments(ctx context.Context) (map[int]string, error)

	//IsIssueCloseValid(ctx context.Context, issue *github.Issue) bool
}

type repo struct {
	config             Config
	installationClient installation.Client
	owner              string
	name               string
	issues             map[int]Issue
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

//func (r *repo) IsIssueCloseValid(ctx context.Context, issue *github.Issue) bool {
//	if r.issues[*issue.Number].Error
//}

func (r *repo) RewardComment(ctx context.Context, issue *github.Issue) (string, error) {
	err := r.loadEventsForIssue(ctx, issue)
	if err != nil {
		return "", err
	}

	contributors := r.ContributorsFromIssues()
	comment := RewardComment(r.issues[*issue.Number], contributors, r.config.Currency)

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
		contributors := r.ContributorsForIssue(issueNumber)

		comments[issueNumber] = RewardComment(r.issues[issueNumber], contributors, r.config.Currency)
	}

	return comments, nil
}

func (r *repo) loadEventsForIssue(ctx context.Context, issue *github.Issue) error {
	events, err := r.installationClient.GetIssueEvents(ctx, r.owner, r.name, *issue.Number)
	if err != nil {
		return err
	}

	r.issues = make(map[int]Issue, 1)
	r.issues[*issue.Number] = Issue{
		Issue:  issue,
		Events: events,
	}

	return nil
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

	issuesEvents, errs := r.installationClient.GetIssuesEvents(ctx, r.owner, r.name, filteredIssues)
	for i, err := range errs {
		log.Printf("[loadEventsForIssue] error while requesting events for issue with number %d: %v", i, err)
	}

	// Add issues to repo
	r.issues = make(map[int]Issue, len(filteredIssues))
	for _, issue := range filteredIssues {
		wrappedIssue := Issue{Issue: issue}

		// Add error to wrapped issue
		err, ok := errs[*issue.Number]
		if ok {
			wrappedIssue.Error = err
			r.issues[*issue.Number] = wrappedIssue
			continue
		}

		// Add events to wrapped issue
		events, ok := issuesEvents[*issue.Number]
		if ok {
			wrappedIssue.Events = events
		}
		r.issues[*issue.Number] = wrappedIssue
	}

	return nil
}
