package model

import (
	"github.com/phuslu/log"
	"time"

	"github.com/morphysm/famed-github-backend/internal/repositories/github/model"
)

// NewBlueTeamFromIssues generates a blue team from a given slice of enriched issues.
func NewBlueTeamFromIssues(issues map[int]model.EnrichedIssue, options BoardOptions) []*Contributor {
	// Map issues and events to contributors
	contributors := issuesToBlueTeam(issues, options)
	// Calculate mean and deviation of time to disclosure
	contributors.updateMeanAndDeviationOfDisclosure()
	// Calculate average severity of fixed issues
	contributors.updateAverageSeverity()
	// Transformation of contributors map to contributors array
	contributorsArray := contributors.toSortedSlice()
	return contributorsArray
}

// NewBlueTeamFromIssue returns a contributors map generated from the repo's internal issue with issueID
// and its corresponding events.
func NewBlueTeamFromIssue(issue model.EnrichedIssue, options BoardOptions) ([]*Contributor, error) {
	contributors := Contributors{}
	// Map issue to contributors
	err := contributors.mapBlueTeamIssue(issue, options)
	if err != nil {
		log.Error().Err(err).Msgf("[contributors] error while mapping issue with ID: %d", issue.ID)
		return nil, err
	}
	// Transformation of contributors map to contributors array
	contributorsArray := contributors.toSortedSlice()
	return contributorsArray, nil
}

func issuesToBlueTeam(issues map[int]model.EnrichedIssue, options BoardOptions) Contributors {
	contributors := Contributors{}
	for issueID, issue := range issues {
		// Map issue to contributors
		err := contributors.mapBlueTeamIssue(issue, options)
		if err != nil {
			log.Error().Err(err).Msgf("[issuesToBlueTeam] error while mapping issue with ID: %d", issueID)
			issues[issueID] = issue
		}
	}

	return contributors
}

// mapBlueTeamIssue updates the contributors map based on a set of events and an issue.
func (cs Contributors) mapBlueTeamIssue(issue model.EnrichedIssue, boardOptions BoardOptions) error {
	// Check if issue has closed at timestamp
	if issue.ClosedAt == nil {
		return ErrIssueMissingClosedAt
	}
	issueClosedAt := *issue.ClosedAt
	timeToDisclosure := issueClosedAt.Sub(issue.CreatedAt).Minutes()

	severity, err := issue.Severity()
	if err != nil {
		log.Error().Err(err).Msgf("[mapBlueTeamIssue] error while reading severity from with id: %d", issue.ID)
		return err
	}

	var workLogs WorkLogs
	var reopenCount int
	if !issue.Migrated {
		workLogs, reopenCount = cs.mapBlueTeamEvents(issue.Events, issueClosedAt, severity, timeToDisclosure, boardOptions.Currency, boardOptions.Now)
	}
	if issue.Migrated {
		for _, assignee := range issue.Assignees {
			cs.mapAssigneeIfMissing(assignee, boardOptions.Currency, boardOptions.Now)
			workLogs = WorkLogs{}
			workLogs.Add(assignee.Login, WorkLog{issue.CreatedAt, issueClosedAt})
			cs.incrementFixCounters(assignee.Login, timeToDisclosure, severity)
		}
	}

	// Calculate the reward
	cs.UpdateRewards(issue.HTMLURL, workLogs, issue.CreatedAt, issueClosedAt, reopenCount, severity, boardOptions)

	return nil
}

// mapBlueTeamEvents maps issue events to the contributors
func (cs Contributors) mapBlueTeamEvents(events []model.IssueEvent, issueClosedAt time.Time, severity model.IssueSeverity, timeToDisclosure float64, currency string, now time.Time) (WorkLogs, int) {
	// areIncremented tracks contributors that have had their fix counters incremented
	var (
		workLogs       = WorkLogs{}
		areIncremented = make(map[string]bool)
		reopenCount    = 0
	)

	// Iterate through issue events and map events if event type is of interest
	for _, event := range events {
		switch event.Event {
		case string(model.IssueEventActionAssigned):
			if event.Assignee == nil {
				log.Warn().Msgf("[mapBlueTeamIssue] event assigned is missing for event with ID: %d", event.ID)
				continue
			}
			if event.CreatedAt.After(issueClosedAt) {
				continue
			}

			cs.mapEventAssigned(event, issueClosedAt, workLogs, currency, now)

			// Increment fix count if not yet done
			if isIncremented := areIncremented[event.Assignee.Login]; !isIncremented {
				cs.incrementFixCounters(event.Assignee.Login, timeToDisclosure, severity)
				areIncremented[event.Assignee.Login] = true
			}
		case string(model.IssueEventActionUnassigned):
			mapEventUnassigned(event, workLogs)
		case string(model.IssueEventActionReopened):
			reopenCount++
		}
	}

	return workLogs, reopenCount
}

// mapEventAssigned handles an assigned event, updating the contributor map.
func (cs Contributors) mapEventAssigned(event model.IssueEvent, issueClosedAt time.Time, workLogs WorkLogs, currency string, now time.Time) {
	cs.mapAssigneeIfMissing(*event.Assignee, currency, now)

	// Append work log
	workLogs.Add(event.Assignee.Login, WorkLog{event.CreatedAt, issueClosedAt})
}

// mapEventUnassigned handles an unassigned event, updating the work log of the unassigned contributor.
func mapEventUnassigned(event model.IssueEvent, workLogs WorkLogs) {
	err := workLogs.UpdateEnd(event.Assignee.Login, event.CreatedAt)
	if err != nil {
		log.Error().Err(err).Msgf("[mapEventUnassigned] error on map of event with id %d \n", event.ID)
	}
}
