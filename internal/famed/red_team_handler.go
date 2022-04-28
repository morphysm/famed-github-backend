package famed

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/morphysm/famed-github-backend/internal/client/github"
	"github.com/morphysm/famed-github-backend/internal/config"
)

func (gH *githubHandler) GetRedTeam(c echo.Context) error {
	owner := c.Param("owner")
	if owner == "" {
		return echo.ErrBadRequest.SetInternal(ErrMissingOwnerPathParameter)
	}

	repoName := c.Param("repo_name")
	if repoName == "" {
		return echo.ErrBadRequest.SetInternal(ErrMissingRepoPathParameter)
	}

	redTeam, err := gH.generateRedTeamFromIssues(c.Request().Context(), owner, repoName)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, redTeam)
}

func (gH *githubHandler) generateRedTeamFromIssues(ctx context.Context, owner string, repo string) ([]*Contributor, error) {
	if ok := gH.githubInstallationClient.CheckInstallation(owner); !ok {
		return nil, ErrAppNotInstalled
	}

	famedLabel := gH.famedConfig.Labels[config.FamedLabel]
	issueState := github.All
	issues, err := gH.githubInstallationClient.GetIssuesByRepo(ctx, owner, repo, []string{famedLabel.Name}, &issueState)
	if err != nil {
		return nil, err
	}

	contributors := make(Contributors)
	for _, issue := range issues {
		if issue.RedTeam == nil || issue.BountyPoints == nil || issue.ClosedAt == nil {
			continue
		}

		contributors.mapIssue(issue, gH.famedConfig.Currency)
	}

	contributors.updateMeanAndDeviationOfDisclosure()
	contributors.updateAverageSeverity()
	contributors.updateMonthlyRewards()

	return contributors.toSortedSlice(), nil
}

// mapIssue maps a bug to the contributors map.
func (cs Contributors) mapIssue(issue github.Issue, currency string) {
	// Get red team contributor from map
	for _, teamer := range issue.RedTeam {
		cs.mapAssigneeIfMissing(teamer, currency)
		contributor := cs[teamer.Login]

		severity, err := severity(issue.Labels)
		if err != nil {
			return
		}

		contributor.mapIssue(issue.HTMLURL, issue.CreatedAt, *issue.ClosedAt, float64(*issue.BountyPoints)/float64(len(issue.RedTeam)), severity)
	}
}

// mapIssue maps a bug to a contributor
func (c *Contributor) mapIssue(url string, reportedDate, publishedDate time.Time, reward float64, severity config.IssueSeverity) {
	// Set reward
	c.Rewards = append(c.Rewards, Reward{
		Date:   publishedDate,
		Reward: reward,
		URL:    url,
	})

	// Updated reward sum
	c.RewardSum += reward

	// Increment fix count
	c.FixCount++
	severityCount := c.Severities[severity]
	severityCount++
	c.Severities[severity] = severityCount

	// Update times to disclosure
	c.TimeToDisclosure.Time = append(c.TimeToDisclosure.Time, publishedDate.Sub(reportedDate).Minutes())
}
