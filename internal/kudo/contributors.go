package kudo

import (
	"log"
	"math"
	"time"

	"github.com/google/go-github/v41/github"

	"github.com/morphysm/kudos-github-backend/internal/client/installation"
)

type Contributors map[string]*Contributor

type Contributor struct {
	Login            string                `json:"login"`
	AvatarURL        *string               `json:"avatar_url,omitempty"`
	HTMLURL          *string               `json:"html_url,omitempty"`
	GravatarID       *string               `json:"gravatar_id,omitempty"`
	FixCount         int                   `json:"fix_count,omitempty"`
	MonthFixCount    map[time.Month]int    `json:"month_fix_count,omitempty"`
	Rewards          []Reward              `json:"rewards"`
	RewardSum        float64               `json:"reward_sum"`
	TimeToDisclosure TimeToDisclosure      `json:"time_to_disclosure"`
	IssueSeverities  map[IssueSeverity]int `json:"issue_severity"`
	MeanSeverity     float64               `json:"mean_severity"`
}

type TimeToDisclosure struct {
	Time              []float64 `json:"time"`
	Mean              float64   `json:"mean"`
	StandardDeviation float64   `json:"standard_deviation"`
}

type WorkLog struct {
	Start time.Time
	End   time.Time
}

type Reward struct {
	Date   time.Time `json:"date"`
	Reward float64   `json:"reward"`
}

// GenerateContributors creates a contributors map based on an array of issue and a map of event arrays.
func GenerateContributors(issues []*github.Issue, eventsByIssue map[int64][]*github.IssueEvent) Contributors {
	contributors := Contributors{}
	for _, issue := range issues {
		contributors.MapIssue(issue, eventsByIssue[*issue.ID])
	}

	return contributors
}

// MapIssue updates the contributors map based on a set of events and an issue.
func (contributors Contributors) MapIssue(issue *github.Issue, events []*github.IssueEvent) {
	var (
		workLogs         = map[string][]WorkLog{}
		reopenCount      = 0
		issueCreatedAt   = *issue.CreatedAt
		issueClosedAt    = *issue.ClosedAt
		timeToDisclosure = issueClosedAt.Sub(issueCreatedAt).Minutes()
		severity         = IssueToSeverity(issue)
	)

	// Add on closed assignee
	contributors.mapAssigneeIfMissing(issue.Assignee)
	// Increment fix counter only for assignee on closed
	contributors.updateFixCounters(issue, timeToDisclosure, severity)

	for _, event := range events {
		if event.Event == nil {
			continue
		}

		switch *event.Event {
		case string(installation.IssueEventActionAssigned):
			contributors.mapEventAssigned(event, issueClosedAt, workLogs)
		case string(installation.IssueEventActionUnassigned):
			mapEventUnassigned(event, workLogs)
		case string(installation.IssueEventActionReopened):
			reopenCount++
		}
	}

	// Calculate the reward
	contributors.updateReward(workLogs, issueCreatedAt, issueClosedAt, reopenCount)
	// Calculate mean and deviation of time to disclosure
	contributors.updateMeanAndDeviationOfDisclosure()
	// Calculate average severity of fixed issues
	contributors.updateAverageSeverity()
}

// mapAssigneeIfMissing adds a contributor to the contributors' map if the contributor is missing.
func (contributors Contributors) mapAssigneeIfMissing(assignee *github.User) {
	_, ok := contributors[*assignee.Login]
	if !ok {
		contributors[*assignee.Login] = &Contributor{
			Login:            *assignee.Login,
			AvatarURL:        assignee.AvatarURL,
			HTMLURL:          assignee.HTMLURL,
			GravatarID:       assignee.GravatarID,
			Rewards:          []Reward{},
			TimeToDisclosure: TimeToDisclosure{},
			IssueSeverities:  map[IssueSeverity]int{},
			MonthFixCount:    map[time.Month]int{},
		}
	}
}

// updateFixCounters updates the fix counters of the contributor who is assigned to the issue in the contributors' map.
func (contributors Contributors) updateFixCounters(issue *github.Issue, timeToDisclosure float64, severity IssueSeverity) {
	contributor, _ := contributors[*issue.Assignee.Login]

	// Increment fix count
	contributor.FixCount++
	monthCount := contributor.MonthFixCount[issue.ClosedAt.Month()]
	monthCount++
	contributor.MonthFixCount[issue.ClosedAt.Month()] = monthCount

	// Increment severity counter
	counterSeverities := contributor.IssueSeverities[severity]
	contributor.IssueSeverities[severity] = counterSeverities + 1

	// Append time to disclosure
	contributor.TimeToDisclosure.Time = append(contributor.TimeToDisclosure.Time, timeToDisclosure)
}

// mapEventAssigned handles an assigned event, updating the contributor map.
func (contributors Contributors) mapEventAssigned(event *github.IssueEvent, issueClosedAt time.Time, workLogs map[string][]WorkLog) {
	if event.Assignee == nil || event.Assignee.Login == nil || event.CreatedAt == nil {
		return
	}

	if event.CreatedAt.After(issueClosedAt) {
		return
	}

	contributors.mapAssigneeIfMissing(event.Assignee)

	// Append work log
	// TODO check if work end works like this
	work := WorkLog{Start: *event.CreatedAt, End: issueClosedAt}
	assigneeWorkLogs := workLogs[*event.Assignee.Login]
	assigneeWorkLogs = append(assigneeWorkLogs, work)
	workLogs[*event.Assignee.Login] = assigneeWorkLogs
}

// mapEventUnassigned handles an unassigned event, updating the work log of the unassigned contributor.
func mapEventUnassigned(event *github.IssueEvent, workLogs map[string][]WorkLog) {
	if event.Assignee == nil || event.Assignee.Login == nil || event.CreatedAt == nil {
		return
	}

	// Append work log
	assigneeWorkLogs := workLogs[*event.Assignee.Login]
	if len(assigneeWorkLogs) == 0 {
		log.Printf("no work log on event unassigned")
		return
	}

	assigneeWorkLogs[len(assigneeWorkLogs)-1].End = *event.CreatedAt
	workLogs[*event.Assignee.Login] = assigneeWorkLogs
}

// updateMeanAndDeviationOfDisclosure updates the mean and deviation of the time to disclosure of all contributors.
func (contributors Contributors) updateMeanAndDeviationOfDisclosure() {
	for _, contributor := range contributors {
		if contributor.FixCount == 0 {
			continue
		}

		var totalTime, sd float64
		// Calculate mean
		for _, timeToDisclosure := range contributor.TimeToDisclosure.Time {
			totalTime += timeToDisclosure
		}

		contributor.TimeToDisclosure.Mean = totalTime / float64(contributor.FixCount)

		// Calculate standard deviation
		for _, timeToDisclosure := range contributor.TimeToDisclosure.Time {
			sd += math.Pow(timeToDisclosure-contributor.TimeToDisclosure.Mean, 2) //nolint:gomnd
		}

		contributor.TimeToDisclosure.StandardDeviation = math.Sqrt(sd / float64(contributor.FixCount))
	}
}

// updateAverageSeverity updates the average severity field of all contributors.
func (contributors Contributors) updateAverageSeverity() {
	for _, contributor := range contributors {
		if contributor.FixCount == 0 {
			continue
		}

		contributor.MeanSeverity = (2*float64(contributor.IssueSeverities[IssueSeverityLow]) +
			5.5*float64(contributor.IssueSeverities[IssueSeverityMedium]) +
			9*float64(contributor.IssueSeverities[IssueSeverityHigh]) +
			9.5*float64(contributor.IssueSeverities[IssueSeverityCritical])) / float64(contributor.FixCount)
	}
}
