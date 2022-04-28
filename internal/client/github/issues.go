package github

import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/v41/github"
	"github.com/morphysm/famed-github-backend/pkg/parse"
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

type IssueSeverity string

const (
	// Info represents a CVSS of 0
	Info IssueSeverity = "info"
	// Low represents a CVSS of 0.1-3.9
	Low IssueSeverity = "low"
	// Medium represents a CVSS of 4.0-6.9
	Medium IssueSeverity = "medium"
	// High represents a CVSS of 7.0-8.9
	High IssueSeverity = "high"
	// Critical represents a CVSS of 9.0-10.0
	Critical IssueSeverity = "critical"
)

type Issue struct {
	ID           int64
	Number       int
	HTMLURL      string
	Title        string
	CreatedAt    time.Time
	ClosedAt     *time.Time
	Assignees    []User
	Severities   []IssueSeverity
	Migrated     bool
	RedTeam      []User
	BountyPoints *int
}

// Severity returns the issue severity.
// If 0 or more than one severity are present it returns an error.
func (i *Issue) Severity() (IssueSeverity, error) {
	// Check for single severity
	if len(i.Severities) == 0 {
		return "", ErrIssueMissingSeverityLabel
	}
	if len(i.Severities) > 1 {
		return "", ErrIssueMultipleSeverityLabels
	}

	return i.Severities[0], nil
}

// GetIssuesByRepo returns all issues from a given repository.
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
	} else {
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
			log.Printf("[GetIssuesByRepo] validation error for issue with number %d: %v", issue.Number, err)
		}

		// TODO refactor this code and functions
		// Detecting migrated ethereum issues
		if isMigratedIssue(issue, compressedIssue, owner, repoName) {
			compressedIssue = c.parseMigrationIssue(ctx, owner, compressedIssue, *issue.Body)
		}

		allCompressedIssues = append(allCompressedIssues, compressedIssue)
	}

	return allCompressedIssues, nil
}

func isMigratedIssue(issue *github.Issue, compressedIssue Issue, owner string, repoName string) bool {
	return issue.Body != nil &&
		(strings.Contains(compressedIssue.Title, "Famed Retroactive Rewards") ||
			(owner == "ethereum" &&
				repoName == "public-disclosures"))
}

func validateIssue(issue *github.Issue) (Issue, error) {
	var compressedIssue Issue
	if issue == nil ||
		issue.ID == nil ||
		issue.Number == nil ||
		issue.HTMLURL == nil ||
		issue.Title == nil ||
		issue.CreatedAt == nil ||
		issue.Labels == nil {
		return Issue{}, ErrIssueMissingData
	}

	compressedIssue = Issue{
		ID:         *issue.ID,
		Number:     *issue.Number,
		HTMLURL:    *issue.HTMLURL,
		Title:      *issue.Title,
		CreatedAt:  *issue.CreatedAt,
		ClosedAt:   issue.ClosedAt,
		Severities: parseSeverities(issue.Labels),
	}

	for _, assignee := range issue.Assignees {
		assignee, err := validateUser(assignee)
		if err == nil {
			compressedIssue.Assignees = append(compressedIssue.Assignees, assignee)
		}
	}

	return compressedIssue, nil
}

// severity returns the issue severity by matching labels against CVSS
// if no matching issue severity label can be found it returns the IssueMissingLabelErr
// if multiple matching issue severity labels can be found it returns the IssueMultipleSeverityLabelsErr.
func parseSeverities(labels []*github.Label) []IssueSeverity {
	var severities []IssueSeverity
	for _, label := range labels {
		// Check if label is equal to one of the predefined severity labels.
		if (label != nil &&
			label.Name != nil) &&
			(*label.Name == string(Info) ||
				*label.Name == string(Low) ||
				*label.Name == string(Medium) ||
				*label.Name == string(High) ||
				*label.Name == string(Critical)) {
			severities = append(severities, IssueSeverity(*label.Name))
		}
	}

	return severities
}

// parseMigrationIssue returns an updated issue with migration info parsed from a GitHub issue body.
func (c *githubInstallationClient) parseMigrationIssue(ctx context.Context, owner string, issue Issue, body string) Issue {
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

	redTeam, err := c.parseRedTeam(ctx, owner, body)
	if err != nil {
		log.Printf("[parseMigrationIssue] error while parsing red teamer: %v", err)
	} else {
		issue.RedTeam = redTeam
	}

	bountyPoints, err := parseBountyPoints(body)
	if err != nil {
		log.Printf("[parseMigrationIssue] error while parsing bounty points: %v", err)
	} else {
		issue.BountyPoints = &bountyPoints
	}

	return issue
}

// parseRedTeam returns a red team parsed from a GitHub issue body.
func (c *githubInstallationClient) parseRedTeam(ctx context.Context, owner string, body string) ([]User, error) {
	var users []User

	// Parse red team from issue body
	redTeam, err := parse.FindRightOfKey(body, "Bounty Hunter:")
	if err != nil {
		return nil, err
	}

	// Split bounty hunters if two are present separated by ", "
	splitTeam := strings.Split(redTeam, ", ")
	for _, teamer := range splitTeam {
		// Read from known GitHub logins
		login := c.redTeamLogins[teamer]
		if login == "" {
			log.Printf("[parseRedTeam] no GitHub login found for red teamer %s", teamer)
			users = append(users, User{Login: teamer})
			continue
		}

		// Check if red teamer is in cache
		cachedTeamer, ok := c.cachedRedTeam.Get(login)
		if ok {
			users = append(users, cachedTeamer)
			continue
		}

		// Fetch user info
		user, err := c.GetUser(ctx, owner, login)
		if err != nil {
			log.Printf("[parseRedTeam] error while retrieving user icon for login: %s: %v", login, err)
			users = append(users, User{Login: teamer})
			continue
		}

		// Add user info to cache
		c.cachedRedTeam.Add(user)
		users = append(users, user)
	}

	return users, nil
}

// parseReportedTime returns the report time parsed from a GitHub issue body.
func parseReportedTime(body string) (time.Time, error) {
	value, err := parse.FindRightOfKey(body, "Reported:")
	if err != nil {
		return time.Time{}, err
	}

	createdAt, err := parseDate(value)
	if err != nil {
		return time.Time{}, err
	}

	return createdAt, nil
}

// parseFixTime returns the fix time parsed from a GitHub issue body.
func parseFixTime(body string) (time.Time, error) {
	value, err := parse.FindRightOfKey(body, "Fixed:")
	if err != nil {
		return time.Time{}, err
	}

	closedAt, err := parseDate(value)
	if err != nil {
		return time.Time{}, err
	}

	return closedAt, nil
}

// parseBountyPoints returns the bounty points parsed from a GitHub issue body.
func parseBountyPoints(body string) (int, error) {
	value, err := parse.FindRightOfKey(body, "Bounty Points:")
	if err != nil {
		return -1, err
	}

	bountyPoints, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return -1, err
	}

	return int(bountyPoints), nil
}

// parseDate returns a date string parsed with "YYYY-MM-DD" format to time.Time.
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
