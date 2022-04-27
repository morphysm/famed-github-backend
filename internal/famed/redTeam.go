package famed

import (
	"context"
	"log"
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
		if issue.RedTeamer == nil || issue.BountyPoints == nil || issue.ClosedAt == nil {
			continue
		}

		contributors.mapBug(ctx, gH.githubInstallationClient, owner, issue, gH.famedConfig.RedTeamers, gH.famedConfig.Currency)
	}

	contributors.updateMeanAndDeviationOfDisclosure()
	contributors.updateAverageSeverity()
	contributors.updateMonthlyRewards()

	return contributors.toSortedSlice(), nil
}

// mapBug maps a bug to the contributors map.
func (cs Contributors) mapBug(ctx context.Context, client github.InstallationClient, owner string, issue github.Issue, redTeamers map[string]string, currency string) {
	// Get red team contributor from map
	contributor, ok := cs[issue.RedTeamer.Login]
	if !ok {
		contributor = &Contributor{
			Severities: make(map[config.IssueSeverity]int),
		}
		cs[issue.RedTeamer.Login] = contributor

		// Set login
		login := redTeamers[issue.RedTeamer.Login]
		if login != "" {
			contributor.Login = login

			// Get icons
			err := contributor.addUserIcon(ctx, client, owner)
			if err != nil {
				log.Printf("error while retrieving user icon for for red teamer: %s: %v", issue.RedTeamer.Login, err)
			}
		}
		if login == "" {
			log.Printf("no GitHub login found for red teamer %s", issue.RedTeamer.Login)
			contributor.Login = issue.RedTeamer.Login
		}

		// Set currency
		contributor.Currency = currency
	}

	severity, err := severity(issue.Labels)
	if err != nil {
		return
	}

	contributor.mapBug(issue.CreatedAt, *issue.ClosedAt, float64(*issue.BountyPoints), severity)
}

// mapBug maps a bug to a contributor
func (c *Contributor) mapBug(reportedDate, publishedDate time.Time, reward float64, severity config.IssueSeverity) {
	// Set reward
	c.Rewards = append(c.Rewards, Reward{Date: publishedDate, Reward: reward})

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

// addUserIcon adds a user icon to a contributor fetched from the GitHub API.
func (c *Contributor) addUserIcon(ctx context.Context, client github.InstallationClient, owner string) error {
	user, err := client.GetUser(ctx, owner, c.Login)
	if err != nil {
		return err
	}

	c.AvatarURL = user.AvatarURL
	c.HTMLURL = user.HTMLURL

	return nil
}
