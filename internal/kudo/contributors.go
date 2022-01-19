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
	Login            string                 `json:"login"`
	AvatarURL        *string                `json:"avatar_url,omitempty"`
	HTMLURL          *string                `json:"html_url,omitempty"`
	GravatarID       *string                `json:"gravatar_id,omitempty"`
	FixCount         int                    `json:"fix_count,omitempty"`
	Rewards          []Reward               `json:"rewards"`
	RewardSum        float64                `json:"reward_sum"`
	RewardUnit       string                 `json:"reward_unit"`
	MonthlyRewards   map[time.Month]float64 `json:"monthly_rewards,omitempty"`
	TimeToDisclosure TimeToDisclosure       `json:"time_to_disclosure"`
	Severities       map[IssueSeverity]int  `json:"severities"`
	MeanSeverity     float64                `json:"mean_severity"`
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
func GenerateContributors(issues []*github.Issue, eventsByIssue map[int64][]*github.IssueEvent, rewardUnit string, rewards map[IssueSeverity]float64, usdToEthRate float64) Contributors {
	contributors := Contributors{}
	for _, issue := range issues {
		contributors.MapIssue(issue, eventsByIssue[*issue.ID], rewardUnit, rewards, usdToEthRate)
	}

	return contributors
}

// MapIssue updates the contributors map based on a set of events and an issue.
// TODO investigate if different data handling for rewards works
func (contributors Contributors) MapIssue(issue *github.Issue, events []*github.IssueEvent, rewardUnit string, rewards map[IssueSeverity]float64, usdToEthRate float64) {
	var (
		workLogs         = map[string][]WorkLog{}
		reopenCount      = 0
		issueCreatedAt   = *issue.CreatedAt
		issueClosedAt    = *issue.ClosedAt
		timeToDisclosure = issueClosedAt.Sub(issueCreatedAt).Minutes()
	)

	// Read severity from issue
	severity, err := IssueToSeverity(issue)
	if err != nil {
		log.Printf("[MapIssue] no valid label found for issue with ID: %d and label error: %v", issue.ID, err)
		return
	}

	// Get severity reward from config
	severityReward := rewards[severity]

	// Add on closed assignee
	contributors.mapAssigneeIfMissing(issue.Assignee, rewardUnit)
	// Increment fix counter only for assignee on closed
	contributors.updateFixCounters(issue, timeToDisclosure, severity)

	for _, event := range events {
		if event.Event == nil {
			continue
		}

		switch *event.Event {
		case string(installation.IssueEventActionAssigned):
			contributors.mapEventAssigned(event, issueClosedAt, workLogs, rewardUnit)
		case string(installation.IssueEventActionUnassigned):
			mapEventUnassigned(event, workLogs)
		case string(installation.IssueEventActionReopened):
			reopenCount++
		}
	}

	// Calculate the reward
	contributors.updateReward(workLogs, issueCreatedAt, issueClosedAt, reopenCount, severityReward, usdToEthRate)
	// Calculate mean and deviation of time to disclosure
	contributors.updateMeanAndDeviationOfDisclosure()
	// Calculate average severity of fixed issues
	contributors.updateAverageSeverity()
}

// mapAssigneeIfMissing adds a contributor to the contributors' map if the contributor is missing.
func (contributors Contributors) mapAssigneeIfMissing(assignee *github.User, rewardUnit string) {
	_, ok := contributors[*assignee.Login]
	if !ok {
		contributors[*assignee.Login] = &Contributor{
			Login:            *assignee.Login,
			AvatarURL:        assignee.AvatarURL,
			HTMLURL:          assignee.HTMLURL,
			GravatarID:       assignee.GravatarID,
			Rewards:          []Reward{},
			RewardUnit:       rewardUnit,
			TimeToDisclosure: TimeToDisclosure{},
			Severities:       map[IssueSeverity]int{},
			MonthlyRewards:   map[time.Month]float64{},
		}
	}
}

// updateFixCounters updates the fix counters of the contributor who is assigned to the issue in the contributors' map.
func (contributors Contributors) updateFixCounters(issue *github.Issue, timeToDisclosure float64, severity IssueSeverity) {
	contributor, _ := contributors[*issue.Assignee.Login]

	// Increment fix count
	contributor.FixCount++

	// Increment severity counter
	counterSeverities := contributor.Severities[severity]
	contributor.Severities[severity] = counterSeverities + 1

	// Append time to disclosure
	contributor.TimeToDisclosure.Time = append(contributor.TimeToDisclosure.Time, timeToDisclosure)
}

// mapEventAssigned handles an assigned event, updating the contributor map.
func (contributors Contributors) mapEventAssigned(event *github.IssueEvent, issueClosedAt time.Time, workLogs map[string][]WorkLog, rewardUnit string) {
	if event.Assignee == nil || event.Assignee.Login == nil || event.CreatedAt == nil {
		return
	}

	if event.CreatedAt.After(issueClosedAt) {
		return
	}

	contributors.mapAssigneeIfMissing(event.Assignee, rewardUnit)

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
		log.Printf("[mapEventUnassigned] no work log on event unassigned of issue with id %d \n", event.Issue.ID)
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

		// Calculate mean
		var totalTime, sd float64
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

		contributor.MeanSeverity = (2*float64(contributor.Severities[IssueSeverityLow]) +
			5.5*float64(contributor.Severities[IssueSeverityMedium]) +
			9*float64(contributor.Severities[IssueSeverityHigh]) +
			9.5*float64(contributor.Severities[IssueSeverityCritical])) / float64(contributor.FixCount)
	}
}