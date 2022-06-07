package model

import (
	"log"
	"time"

	"github.com/morphysm/famed-github-backend/internal/repositories/github/model"
)

func NewRedTeamFromIssues(issues []model.Issue, currency string, now time.Time) ([]*Contributor, error) {
	contributors := Contributors{}
	if len(issues) == 0 {
		return []*Contributor{}, nil
	}

	contributors.mapRedTeamFromIssues(issues, currency, now)
	contributors.updateMeanAndDeviationOfDisclosure()
	contributors.updateAverageSeverity()

	return contributors.toSortedSlice(), nil
}

// mapBlueTeamIssue maps an issue to the contributors map.
func (cs Contributors) mapRedTeamFromIssues(issues []model.Issue, currency string, now time.Time) {
	for _, issue := range issues {
		if issue.RedTeam == nil || issue.BountyPoints == nil || issue.ClosedAt == nil {
			log.Printf("[mapRedTeamFromIssues] issue with id: %d: is missing data", issue.ID)
			continue
		}

		cs.mapRedTeamFromIssue(issue, currency, now)
	}
}

// mapBlueTeamIssue maps an issue to the contributors map.
func (cs Contributors) mapRedTeamFromIssue(issue model.Issue, currency string, now time.Time) {
	// Get red team contributor from map
	for _, teamer := range issue.RedTeam {
		cs.mapAssigneeIfMissing(teamer, currency, now)
		contributor := cs[teamer.Login]

		severity, err := issue.Severity()
		if err != nil {
			log.Printf("[mapRedTeamFromIssue] error while reading severity from with id: %d: %v", issue.ID, err)
			return
		}

		contributor.mapIssue(issue.HTMLURL, issue.CreatedAt, *issue.ClosedAt, float64(*issue.BountyPoints)/float64(len(issue.RedTeam)), severity, now)
	}
}
