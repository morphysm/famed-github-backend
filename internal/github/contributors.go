package github

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/google/go-github/v41/github"
	"github.com/labstack/echo/v4"

	"github.com/morphysm/kudos-github-backend/internal/client/installation"
	"github.com/morphysm/kudos-github-backend/internal/kudo"
)

func (gH *githubHandler) GetContributors(c echo.Context) error {
	repoName := c.Param("repo_name")
	if repoName == "" {
		return echo.ErrBadRequest.SetInternal(errors.New("missing repo name path parameter"))
	}

	// Get all issues in repo
	issuesResponse, err := gH.githubInstallationClient.GetIssuesByRepo(c.Request().Context(), repoName, []string{gH.kudoLabel}, installation.Closed)
	if err != nil {
		return echo.ErrBadGateway.SetInternal(err)
	}

	// Use issues to generate contributor list
	contributors, err := gH.issuesToContributors(c.Request().Context(), issuesResponse, repoName)
	if err != nil {
		return echo.ErrBadGateway.SetInternal(err)
	}

	return c.JSON(http.StatusOK, contributors)
}

// TODO test if issues are returned in chronological order
func (gH *githubHandler) issuesToContributors(ctx context.Context, issues []*github.Issue, repoName string) (map[string]*kudo.Contributor, error) {
	contributors := map[string]*kudo.Contributor{}

	for _, issue := range issues {
		if issue.ID == nil || issue.CreatedAt == nil || issue.ClosedAt == nil {
			continue
		}

		eventsResp, err := gH.githubInstallationClient.GetIssueEvents(ctx, repoName, *issue.Number)
		if err != nil {
			return nil, err
		}

		severity := issueSeverity(issue)

		contributors = eventsToContributors(contributors, eventsResp, *issue.CreatedAt, *issue.ClosedAt, severity)
	}
	// TODO is this ordered by time of occurrence?

	return contributors, nil
}

// TODO how do we handle multiple CVSS
// issueSeverity returns the issue severity by matching labels against CVSS
// if no matching issue severity label can be found it returns issue severity none
func issueSeverity(issue *github.Issue) kudo.IssueSeverity {
	if issue.Labels == nil {
		return kudo.IssueSeverityNone
	}

	for _, label := range issue.Labels {
		if label.Name == nil {
			continue
		}

		switch *label.Name {
		case string(kudo.IssueSeverityLow):
			return kudo.IssueSeverityLow
		case string(kudo.IssueSeverityMedium):
			return kudo.IssueSeverityMedium
		case string(kudo.IssueSeverityHigh):
			return kudo.IssueSeverityHigh
		case string(kudo.IssueSeverityCritical):
			return kudo.IssueSeverityCritical
		}
	}

	return kudo.IssueSeverityNone
}

func eventsToContributors(contributors map[string]*kudo.Contributor, events []*github.IssueEvent, issueCreatedAt time.Time, issueClosedAt time.Time, severity kudo.IssueSeverity) map[string]*kudo.Contributor {
	var (
		timeToDisclosure = issueClosedAt.Sub(issueCreatedAt).Minutes()
		workLogs         = map[string][]kudo.WorkLog{}
	)

	for _, event := range events {
		if event.Event == nil {
			continue
		}

		switch *event.Event {
		case string(installation.IssueEventActionAssigned):
			if event.Assignee == nil || event.Assignee.Login == nil || event.CreatedAt == nil {
				break
			}

			contributor, ok := contributors[*event.Assignee.Login]
			if !ok {
				contributor = &kudo.Contributor{
					Login:            *event.Assignee.Login,
					AvatarURL:        event.Assignee.AvatarURL,
					HTMLURL:          event.Assignee.HTMLURL,
					GravatarID:       event.Assignee.GravatarID,
					Rewards:          []kudo.Reward{},
					TimeToDisclosure: []float64{timeToDisclosure},
					IssueSeverities:  map[kudo.IssueSeverity]int{},
				}
			}

			// Increment severity counter
			counterSeverities, _ := contributor.IssueSeverities[severity]
			contributor.IssueSeverities[severity] = counterSeverities + 1

			// Append work log
			// TODO check if work end works like this
			work := kudo.WorkLog{Start: *event.CreatedAt, End: issueClosedAt}
			assigneeWorkLogs, _ := workLogs[*event.Assignee.Login]
			assigneeWorkLogs = append(assigneeWorkLogs, work)
			workLogs[*event.Assignee.Login] = assigneeWorkLogs

			contributors[*event.Assignee.Login] = contributor
		case string(installation.IssueEventActionUnassigned):
			if event.Assignee == nil || event.Assignee.Login == nil || event.CreatedAt == nil {
				break
			}

			// Append work log
			assigneeWorkLogs, _ := workLogs[*event.Assignee.Login]
			assigneeWorkLogs[len(assigneeWorkLogs)-1].End = *event.CreatedAt
			workLogs[*event.Assignee.Login] = assigneeWorkLogs
		}
	}

	// Calculate the reward // TODO incorrect because it does not count multiple issues
	contributors = kudo.UpdateReward(contributors, workLogs, issueCreatedAt, issueClosedAt, 0)

	return contributors
}
