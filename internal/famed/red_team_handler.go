package famed

import (
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
		return echo.NewHTTPError(http.StatusBadRequest, ErrMissingOwnerPathParameter.Error())
	}

	repoName := c.Param("repo_name")
	if repoName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, ErrMissingRepoPathParameter.Error())
	}

	if ok := gH.githubInstallationClient.CheckInstallation(owner); !ok {
		return echo.NewHTTPError(http.StatusBadRequest, ErrAppNotInstalled.Error())
	}

	famedLabel := gH.famedConfig.Labels[config.FamedLabelKey]
	issueState := github.All
	issues, err := gH.githubInstallationClient.GetIssuesByRepo(c.Request().Context(), owner, repoName, []string{famedLabel.Name}, &issueState)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadGateway, err.Error())
	}

	redTeam, err := generateRedTeamFromIssues(issues, gH.famedConfig.Currency)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, redTeam)
}

func generateRedTeamFromIssues(issues []github.Issue, currency string) ([]*Contributor, error) {
	contributors := Contributors{}
	if len(issues) == 0 {
		return []*Contributor{}, nil
	}

	for _, issue := range issues {
		if issue.RedTeam == nil || issue.BountyPoints == nil || issue.ClosedAt == nil {
			continue
		}

		contributors.mapIssue(issue, currency)
	}

	contributors.updateMeanAndDeviationOfDisclosure()
	contributors.updateAverageSeverity()

	return contributors.toSortedSlice(), nil
}

// mapIssue maps an issue to the contributors map.
func (cs Contributors) mapIssue(issue github.Issue, currency string) {
	// Get red team contributor from map
	for _, teamer := range issue.RedTeam {
		cs.mapAssigneeIfMissing(teamer, currency)
		contributor := cs[teamer.Login]

		severity, err := issue.Severity()
		if err != nil {
			log.Printf("[MapIssue] error while reading severity from with id: %d: %v", issue.ID, err)
			return
		}

		contributor.mapIssue(issue.HTMLURL, issue.CreatedAt, *issue.ClosedAt, float64(*issue.BountyPoints)/float64(len(issue.RedTeam)), severity)
	}
}

// mapIssue maps an issue to a contributor.
func (c *Contributor) mapIssue(url string, reportedDate, publishedDate time.Time, reward float64, severity github.IssueSeverity) {
	// Set reward
	c.updateReward(url, publishedDate, reward)

	// Increment fix count
	c.incrementFixCounters(publishedDate.Sub(reportedDate).Minutes(), severity)
}
