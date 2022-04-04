package famed

import (
	"log"
	"math"
	"sort"
	"time"

	"github.com/morphysm/famed-github-backend/internal/client/github"
	"github.com/morphysm/famed-github-backend/internal/config"
)

type Contributors map[string]*Contributor

type Contributor struct {
	Login            string                       `json:"login"`
	AvatarURL        string                       `json:"avatarUrl"`
	HTMLURL          string                       `json:"htmlUrl"`
	FixCount         int                          `json:"fixCount"`
	Rewards          []Reward                     `json:"rewards"`
	RewardSum        float64                      `json:"rewardSum"`
	Currency         string                       `json:"currency"`
	RewardsLastYear  RewardsLastYear              `json:"rewardsLastYear"`
	TimeToDisclosure TimeToDisclosure             `json:"timeToDisclosure"`
	Severities       map[config.IssueSeverity]int `json:"severities"`
	MeanSeverity     float64                      `json:"meanSeverity"`
	// For issue rewardComment generation
	TotalWorkTime time.Duration
}

type TimeToDisclosure struct {
	Time              []float64 `json:"time"`
	Mean              float64   `json:"mean"`
	StandardDeviation float64   `json:"standardDeviation"`
}

type Reward struct {
	Date   time.Time `json:"date"`
	Reward float64   `json:"reward"`
}

type BoardOptions struct {
	currency string
	rewards  map[config.IssueSeverity]float64
}

// contributorsArray generates a contributor list based on a list of issues
func contributorsArray(issues map[int]WrappedIssue, options BoardOptions) []*Contributor {
	// Generate the contributors from the issues and events
	contributors := ContributorsFromIssues(issues, options)
	// Transformation of contributors map to contributors array
	contributorsArray := contributors.toSortedSlice()
	// Sort contributors array by total rewards
	sortContributors(contributorsArray)

	return contributorsArray
}

// ContributorsFromIssue returns a contributors map generated from the repo's internal issue with issueID
// and its corresponding events.
func ContributorsFromIssue(issue WrappedIssue, options BoardOptions) (Contributors, error) {
	contributors := Contributors{}
	// Map issue to contributors
	err := contributors.MapIssue(issue, options)
	if err != nil {
		log.Printf("[contributors] error while mapping issue with ID: %d, error: %v", issue.Issue.ID, err)
		return contributors, err
	}

	return contributors, nil
}

// ContributorsFromIssues returns a contributors map generated from the repo's internal issues corresponding events.
func ContributorsFromIssues(issues map[int]WrappedIssue, options BoardOptions) Contributors {
	// Map issues and events to contributors
	contributors := issuesAndEventsToContributors(issues, options)
	// Calculate mean and deviation of time to disclosure
	contributors.updateMeanAndDeviationOfDisclosure()
	// Calculate average severity of fixed issues
	contributors.updateAverageSeverity()

	return contributors
}

func issuesAndEventsToContributors(issues map[int]WrappedIssue, options BoardOptions) Contributors {
	contributors := Contributors{}
	for issueID, issue := range issues {
		// Map issue to contributors
		err := contributors.MapIssue(issue, options)
		if err != nil {
			log.Printf("[issuesAndEventsToContributors] error while mapping issue with ID: %d, error: %v", issueID, err)
			issues[issueID] = issue
		}
	}

	return contributors
}

// MapIssue updates the contributors map based on a set of events and an issue.
func (contributors Contributors) MapIssue(issue WrappedIssue, boardOptions BoardOptions) error {
	// Check if issue has closed at timestamp
	if issue.Issue.ClosedAt == nil {
		return ErrIssueMissingClosedAt
	}
	issueClosedAt := *issue.Issue.ClosedAt
	timeToDisclosure := issueClosedAt.Sub(issue.Issue.CreatedAt).Minutes()

	// Read severity from issue
	severity, err := severity(issue.Issue.Labels)
	if err != nil {
		log.Printf("[MapIssue] no valid label found for issue with ID: %d and label error: %v", issue.Issue.ID, err)
		return err
	}

	var workLogs WorkLogs
	var reopenCount int
	if !issue.Issue.Migrated {
		workLogs, reopenCount = contributors.mapEvents(issue.Events, issueClosedAt, severity, timeToDisclosure, boardOptions.currency)
	}
	if issue.Issue.Migrated {
		contributors.mapAssigneeIfMissing(*issue.Issue.Assignee, boardOptions.currency)
		workLogs = WorkLogs{}
		workLogs.Add(issue.Issue.Assignee.Login, WorkLog{issue.Issue.CreatedAt, issueClosedAt})
		contributors.incrementFixCounters(issue.Issue.Assignee.Login, timeToDisclosure, severity)
	}

	// Get severity reward from config
	severityReward := boardOptions.rewards[severity]

	// Calculate the reward
	contributors.updateRewards(workLogs, issue.Issue.CreatedAt, issueClosedAt, reopenCount, severityReward)

	return nil
}

// mapEvents
func (contributors Contributors) mapEvents(events []github.IssueEvent, issueClosedAt time.Time, severity config.IssueSeverity, timeToDisclosure float64, currency string) (WorkLogs, int) {
	// areIncremented tracks contributors that have had their fix counters incremented
	var (
		workLogs       = WorkLogs{}
		areIncremented = make(map[string]bool)
		reopenCount    = 0
	)

	// Iterate through issue events and map events if event type is of interest
	for _, event := range events {
		switch event.Event {
		case string(github.IssueEventActionAssigned):
			if event.Assignee == nil {
				log.Printf("[MapIssue] event assigned is missing for event with ID: %d", event.ID)
				continue
			}
			if event.CreatedAt.After(issueClosedAt) {
				continue
			}

			contributors.mapEventAssigned(event, issueClosedAt, workLogs, currency)

			// Increment fix count if not yet done
			if isIncremented := areIncremented[event.Assignee.Login]; !isIncremented {
				contributors.incrementFixCounters(event.Assignee.Login, timeToDisclosure, severity)
				areIncremented[event.Assignee.Login] = true
			}
		case string(github.IssueEventActionUnassigned):
			mapEventUnassigned(event, workLogs)
		case string(github.IssueEventActionReopened):
			reopenCount++
		}
	}

	return workLogs, reopenCount
}

// mapEventAssigned handles an assigned event, updating the contributor map.
func (contributors Contributors) mapEventAssigned(event github.IssueEvent, issueClosedAt time.Time, workLogs WorkLogs, currency string) {
	contributors.mapAssigneeIfMissing(*event.Assignee, currency)

	// Append work log
	workLogs.Add(event.Assignee.Login, WorkLog{event.CreatedAt, issueClosedAt})
}

// mapEventUnassigned handles an unassigned event, updating the work log of the unassigned contributor.
func mapEventUnassigned(event github.IssueEvent, workLogs WorkLogs) {
	err := workLogs.UpdateEnd(event.Assignee.Login, event.CreatedAt)
	if err != nil {
		log.Printf("[mapEventUnassigned] %v on map of event with id %d \n", err, event.ID)
	}
}

// mapAssigneeIfMissing adds a contributor to the contributors' map if the contributor is missing.
func (contributors Contributors) mapAssigneeIfMissing(assignee github.User, currency string) {
	_, ok := contributors[assignee.Login]
	if !ok {
		contributors[assignee.Login] = &Contributor{
			Login:            assignee.Login,
			AvatarURL:        assignee.AvatarURL,
			HTMLURL:          assignee.HTMLURL,
			Rewards:          []Reward{},
			Currency:         currency,
			TimeToDisclosure: TimeToDisclosure{},
			Severities:       map[config.IssueSeverity]int{},
			RewardsLastYear:  newRewardsLastYear(time.Now()),
		}
	}
}

// updateFixCounters updates the fix counters of the contributor who is assigned to the issue in the contributors' map.
func (contributors Contributors) incrementFixCounters(login string, timeToDisclosure float64, severity config.IssueSeverity) {
	contributor := contributors[login]

	// Increment fix count
	contributor.FixCount++
	// Increment severity counter
	counterSeverities := contributor.Severities[severity]
	contributor.Severities[severity] = counterSeverities + 1
	// Append time to disclosure
	contributor.TimeToDisclosure.Time = append(contributor.TimeToDisclosure.Time, timeToDisclosure)
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

		contributor.MeanSeverity = (2*float64(contributor.Severities[config.CVSSLow]) +
			5.5*float64(contributor.Severities[config.CVSSMedium]) +
			9*float64(contributor.Severities[config.CVSSHigh]) +
			9.5*float64(contributor.Severities[config.CVSSCritical])) / float64(contributor.FixCount)
	}
}

func (contributors Contributors) toSortedSlice() []*Contributor {
	contributorsSlice := contributors.toSlice()
	sortContributors(contributorsSlice)
	return contributorsSlice
}

// mapToSlice transforms the contributors map to a contributors slice.
func (contributors Contributors) toSlice() []*Contributor {
	contributorsSlice := make([]*Contributor, 0)
	for _, contributor := range contributors {
		contributorsSlice = append(contributorsSlice, contributor)
	}

	return contributorsSlice
}

// sortContributors sorts the contributors by descending reward sum.
func sortContributors(contributors []*Contributor) {
	sort.SliceStable(contributors, func(i, j int) bool {
		return contributors[i].RewardSum > contributors[j].RewardSum
	})
}

// updateMonthlyRewards maps the rewards of each contributor to a monthly timeframe for the past year.
func (contributors Contributors) updateMonthlyRewards() {
	now := time.Now()
	for _, contributor := range contributors {
		contributor.RewardsLastYear = newRewardsLastYear(now)

		for _, reward := range contributor.Rewards {
			if month, ok := isInTheLast12Months(now, reward.Date); ok {
				contributor.RewardsLastYear[month].Reward += reward.Reward
			}
		}
	}
}
