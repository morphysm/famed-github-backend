package model

import (
	"github.com/phuslu/log"
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

func NewIssue(issue *github.Issue, owner string, repoName string) (Issue, error) {
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
		Severities: newSeverity(issue.Labels),
	}

	for _, assignee := range issue.Assignees {
		assignee, err := NewUser(assignee)
		if err == nil {
			compressedIssue.Assignees = append(compressedIssue.Assignees, assignee)
		}
	}

	if isMigratedIssue(issue, compressedIssue, owner, repoName) {
		compressedIssue = parseMigrationIssue(compressedIssue, *issue.Body)
	}

	return compressedIssue, nil
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

func isMigratedIssue(issue *github.Issue, compressedIssue Issue, owner string, repoName string) bool {
	return issue.Body != nil &&
		(strings.Contains(compressedIssue.Title, "Famed Retroactive Rewards") ||
			(owner == "ethereum" &&
				repoName == "public-disclosures"))
}

// parseMigrationIssue returns an updated issue with migration info parsed from a GitHub issue body.
func parseMigrationIssue(issue Issue, body string) Issue {
	issue.Migrated = true

	createdAt, err := parseReportedTime(body)
	if err != nil {
		log.Error().Err(err).Msg("[parseMigrationIssue] error while parsing reported time")
	} else {
		issue.CreatedAt = createdAt
	}

	closedAt, err := parseFixTime(body)
	if err != nil {
		log.Error().Err(err).Msg("[parseMigrationIssue] error while parsing fix time")
	} else {
		issue.ClosedAt = &closedAt
	}

	bountyPoints, err := parseBountyPoints(body)
	if err != nil {
		log.Error().Err(err).Msg("[parseMigrationIssue] error while parsing bounty points")
	} else {
		issue.BountyPoints = &bountyPoints
	}

	return issue
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
