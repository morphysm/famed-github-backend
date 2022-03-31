package github

import (
	"context"
	"errors"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/google/go-github/v41/github"
	"github.com/morphysm/famed-github-backend/pkg/pointer"
)

type IssueState string

const (
	All        IssueState = "all"
	Opened     IssueState = "opened"
	Closed     IssueState = "closed"
	Reopened   IssueState = "reopened"
	Edited     IssueState = "edited"
	Assigned   IssueState = "assigned"
	Unassigned IssueState = "unassigned"
	Labeled    IssueState = "labeled"
	Unlabeled  IssueState = "unlabeled"
)

type Issue struct {
	ID        int64
	Number    int
	Title     string
	CreatedAt time.Time
	ClosedAt  *time.Time
	Assignee  *User
	Labels    []Label
	Migrated  bool
}

type User struct {
	Login     string
	AvatarURL string
	HTMLURL   string
}

func (c *githubInstallationClient) GetIssuesByRepo(ctx context.Context, owner string, repoName string, labels []string, state *IssueState) ([]Issue, error) {
	var (
		client, _           = c.clients.get(owner)
		allIssues           []*github.Issue
		allCompressedIssues []Issue
		listOptions         = &github.IssueListByRepoOptions{
			Labels: labels,
			ListOptions: github.ListOptions{
				Page:    1,
				PerPage: 30,
			},
		}
	)

	if state != nil {
		listOptions.State = string(*state)
	}

	if state == nil {
		listOptions.State = string(All)
	}

	for {
		issues, resp, err := client.Issues.ListByRepo(ctx, owner, repoName, listOptions)
		if err != nil {
			return allCompressedIssues, err
		}
		allIssues = append(allIssues, issues...)
		if resp.NextPage == 0 {
			break
		}
		listOptions.Page = resp.NextPage
	}

	for _, issue := range allIssues {
		compressedIssue, err := validateIssue(issue)
		if err != nil {
			log.Printf("validation error for issue with number %d: %v", issue.Number, err)
		}

		allCompressedIssues = append(allCompressedIssues, compressedIssue)
	}

	return allCompressedIssues, nil
}

func validateIssue(issue *github.Issue) (Issue, error) {
	var compressedIssue Issue
	if issue == nil ||
		issue.ID == nil ||
		issue.Number == nil ||
		issue.Title == nil ||
		issue.CreatedAt == nil {
		return compressedIssue, ErrIssueMissingData
	}

	compressedIssue = Issue{
		ID:        *issue.ID,
		Number:    *issue.Number,
		Title:     *issue.Title,
		CreatedAt: *issue.CreatedAt,
		ClosedAt:  issue.ClosedAt,
	}

	assignee, err := validateUser(issue.Assignee)
	if err == nil {
		compressedIssue.Assignee = &assignee
	}

	if issue.Labels != nil {
		for _, label := range issue.Labels {
			if label.Name == nil {
				continue
			}
			compressedLabel := Label{Name: *label.Name}
			compressedIssue.Labels = append(compressedIssue.Labels, compressedLabel)
		}
	}

	// TODO refactor this code and functions
	if issue.Body != nil && strings.Contains(compressedIssue.Title, "Famed Retroactive Rewards") {
		compressedIssue = parseMigrationIssue(compressedIssue, *issue.Body)
	}

	return compressedIssue, nil
}

func validateUser(user *github.User) (User, error) {
	if user == nil ||
		user.Login == nil {
		return User{}, ErrUserMissingData

	}

	return User{
		Login:     *user.Login,
		AvatarURL: pointer.ToString(user.AvatarURL),
		HTMLURL:   pointer.ToString(user.HTMLURL),
	}, nil
}

func parseMigrationIssue(issue Issue, body string) Issue {
	issue.Migrated = true

	createdAt, err := parseReportedTime(body)
	if err != nil {
		log.Printf("[parseMigrationIssue] error while parsing reported time: %v", err)
	} else {
		issue.CreatedAt = createdAt
	}

	closedAt, err := parseFixTime(body)
	if err != nil {
		log.Printf("[parseMigrationIssue] error while parsing fix time: %v", err)
	} else {
		issue.ClosedAt = &closedAt
	}

	return issue
}

func parseReportedTime(body string) (time.Time, error) {
	r, err := regexp.Compile(`\*\*Reported:\*\*\s*([^\n\r]*)`)
	if err != nil {
		return time.Time{}, err
	}

	matches := r.FindStringSubmatch(body)
	if err != nil {
		return time.Time{}, errors.New("no matches found for reported time")
	}

	createdAt, err := parseDate(matches[1])
	if err != nil {
		return time.Time{}, err
	}

	return createdAt, nil
}

func parseFixTime(body string) (time.Time, error) {
	r, err := regexp.Compile(`\*\*Fixed:\*\*\s*([^\n\r]*)`)
	if err != nil {
		return time.Time{}, err
	}

	matches := r.FindStringSubmatch(body)
	if err != nil {
		return time.Time{}, errors.New("no matches found for fix time")
	}

	createdAt, err := parseDate(matches[1])
	if err != nil {
		return time.Time{}, err
	}

	return createdAt, nil
}

func parseDate(data string) (time.Time, error) {
	const layout = "2006-01-02"

	date, err := time.Parse(layout, data)
	if err != nil {
		return time.Time{}, err
	}

	return date, nil
}

// GetIssuePullRequest returns a pull request if a linked pull request for the given issue can be found.
// This is a workaround for the missing "pull_request" field in the event and issue objects provided by the REST GitHub API.
// https://github.community/t/get-referenced-pull-request-from-issue/14027
func (c *githubInstallationClient) GetIssuePullRequest(ctx context.Context, owner string, repoName string, issueNumber int) (*PullRequest, error) {
	allTimelineItemsConnected, err := c.getConnectedEvents(ctx, owner, repoName, issueNumber)
	if err != nil {
		return nil, err
	}

	var lastConnectedEvent *issueTimelineConnectionItem
	for _, node := range allTimelineItemsConnected {
		if node.ConnectedEvent.Subject.PullRequest.URL != "" &&
			(lastConnectedEvent == nil || lastConnectedEvent.ConnectedEvent.CreatedAt.Before(node.ConnectedEvent.CreatedAt)) {
			// Last pull request connected event
			tmpN := node
			lastConnectedEvent = &tmpN
		}
	}

	if lastConnectedEvent == nil {
		return nil, nil
	}

	allTimelineItemsDisconnected, err := c.getDisconnectedEvents(ctx, owner, repoName, issueNumber)
	if err != nil {
		return nil, err
	}

	for _, node := range allTimelineItemsDisconnected {
		if node.DisconnectedEvent.Subject.PullRequest.URL != "" &&
			lastConnectedEvent.ConnectedEvent.CreatedAt.Before(node.DisconnectedEvent.CreatedAt) {
			// Pull request disconnected after last connected event
			return nil, nil
		}
	}

	return &lastConnectedEvent.ConnectedEvent.Subject.PullRequest, nil
}
