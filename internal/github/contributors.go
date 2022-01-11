package github

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/morphysm/kudos-github-backend/internal/client/installation"
	"github.com/morphysm/kudos-github-backend/internal/kudo"
)

func (gH *githubHandler) GetContributors(c echo.Context) error {
	repoName := c.Param("repo_name")
	if repoName == "" {
		return echo.ErrBadRequest.SetInternal(errors.New("missing repo name path parameter"))
	}

	issuesResponse, err := gH.githubInstallationClient.GetIssues(c.Request().Context(), repoName, gH.kudoLabel, installation.Closed)
	if err != nil {
		return echo.ErrBadGateway.SetInternal(err)
	}

	// TODO improve & handle no issues
	contributors, err := gH.issuesToContributors(c.Request().Context(), issuesResponse, repoName)
	if err != nil {
		return echo.ErrBadGateway.SetInternal(err)
	}

	return c.JSON(http.StatusOK, contributors)
}

// https://docs.github.com/en/developers/webhooks-and-events/events/issue-event-types
const (
	eventLabeled  = "labeled"
	eventAssigned = "assigned"
	eventAdded    = "added_to_project"
	eventClosed   = "closed"
)

// TODO test if issues are returned in chronological order
func (gH *githubHandler) issuesToContributors(ctx context.Context, issues installation.IssueResponse, repoName string) (kudo.Contributors, error) {
	var contributors kudo.Contributors

	for _, issue := range issues {
		// TODO add different assignment times
		if issue.ID == nil || issue.CreatedAt == nil || issue.ClosedAt == nil {
			continue
		}

		eventsResp, err := gH.githubInstallationClient.GetIssueEvents(ctx, repoName, *issue.ID)
		if err != nil {
			return contributors, err
		}

		for _, event := range eventsResp {
			if event.Event == nil {
				continue
			}
			//if *event.Event == string(installation.ActionAssigned) {
			//	if len(contributors) == 0 {
			//		contributors = kudo.Contributors{
			//			{Name: *event..Login,
			//				Work: []kudo.Work{{Start: *issue.CreatedAt, End: *issue.ClosedAt}},
			//			}}
			//	}
			//}
		}

		//TODO check for existence of assignee and Login etc.
		//contributors = []*kudo.Contributor{
		//	{Name: *issue.Assignee.Login,
		//		// TODO generate work logs from events
		//		Work: []kudo.Work{{Start: *issue.CreatedAt, End: *issue.ClosedAt}},
		//	}}

		// Calculate the reward
		contributors.Reward(*issue.CreatedAt, *issue.ClosedAt, 0)
	}

	return contributors, nil
}
