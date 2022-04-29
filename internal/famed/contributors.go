package famed

import (
	"log"
	"math"
	"sort"
	"time"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"

	"github.com/morphysm/famed-github-backend/internal/client/github"
)

type contributors map[string]*contributor

type contributor struct {
	Login            string                       `json:"login"`
	AvatarURL        string                       `json:"avatarUrl"`
	HTMLURL          string                       `json:"htmlUrl"`
	FixCount         int                          `json:"fixCount"`
	Rewards          []rewardEvent                `json:"rewards"`
	RewardSum        float64                      `json:"rewardSum"`
	Currency         string                       `json:"currency"`
	RewardsLastYear  rewardsLastYear              `json:"rewardsLastYear"`
	TimeToDisclosure timeToDisclosure             `json:"timeToDisclosure"`
	Severities       map[github.IssueSeverity]int `json:"severities"`
	MeanSeverity     float64                      `json:"meanSeverity"`
	// For issue rewardComment generation
	TotalWorkTime time.Duration `json:"-"`
}

type timeToDisclosure struct {
	Time              []float64 `json:"time"`
	Mean              float64   `json:"mean"`
	StandardDeviation float64   `json:"standardDeviation"`
}

type rewardEvent struct {
	Date   time.Time `json:"date"`
	Reward float64   `json:"reward"`
	URL    string    `json:"url"`
}

type boardOptions struct {
	currency  string
	rewards   map[github.IssueSeverity]float64
	daysToFix int
}

// contributorsArray generates a contributor list based on a list of issues
func contributorsArray(issues map[int]enrichedIssue, options boardOptions) []*contributor {
	// Generate the contributors from the issues and events
	contributors := ContributorsFromIssues(issues, options)
	// Transformation of contributors map to contributors array
	contributorsArray := contributors.toSortedSlice()

	return contributorsArray
}

// ContributorsFromIssue returns a contributors map generated from the repo's internal issue with issueID
// and its corresponding events.
func ContributorsFromIssue(issue enrichedIssue, options boardOptions) (contributors, error) {
	contributors := contributors{}
	// Map issue to contributors
	err := contributors.MapIssue(issue, options)
	if err != nil {
		log.Printf("[contributors] error while mapping issue with ID: %d, error: %v", issue.ID, err)
		return contributors, err
	}

	return contributors, nil
}

// ContributorsFromIssues returns a contributors map generated from the repo's internal issues corresponding events.
func ContributorsFromIssues(issues map[int]enrichedIssue, options boardOptions) contributors {
	// Map issues and events to contributors
	contributors := issuesAndEventsToContributors(issues, options)
	// Calculate mean and deviation of time to disclosure
	contributors.updateMeanAndDeviationOfDisclosure()
	// Calculate average severity of fixed issues
	contributors.updateAverageSeverity()

	return contributors
}

func issuesAndEventsToContributors(issues map[int]enrichedIssue, options boardOptions) contributors {
	contributors := contributors{}
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
func (cs contributors) MapIssue(issue enrichedIssue, boardOptions boardOptions) error {
	// Check if issue has closed at timestamp
	if issue.ClosedAt == nil {
		return ErrIssueMissingClosedAt
	}
	issueClosedAt := *issue.ClosedAt
	timeToDisclosure := issueClosedAt.Sub(issue.CreatedAt).Minutes()

	severity, err := issue.Severity()
	if err != nil {
		log.Printf("[MapIssue] error while reading severity from with id: %d: %v", issue.ID, err)
		return err
	}

	var workLogs WorkLogs
	var reopenCount int
	if !issue.Migrated {
		workLogs, reopenCount = cs.mapEvents(issue.Events, issueClosedAt, severity, timeToDisclosure, boardOptions.currency)
	}
	if issue.Migrated {
		for _, assignee := range issue.Assignees {
			cs.mapAssigneeIfMissing(assignee, boardOptions.currency)
			workLogs = WorkLogs{}
			workLogs.Add(assignee.Login, WorkLog{issue.CreatedAt, issueClosedAt})
			cs.incrementFixCounters(assignee.Login, timeToDisclosure, severity)
		}
	}

	// Get severity reward from config
	severityReward := boardOptions.rewards[severity]

	// Calculate the reward
	cs.updateRewards(issue.HTMLURL, workLogs, issue.CreatedAt, issueClosedAt, reopenCount, boardOptions.daysToFix, severityReward)

	return nil
}

// mapEvents maps issue events to the contributors
func (cs contributors) mapEvents(events []github.IssueEvent, issueClosedAt time.Time, severity github.IssueSeverity, timeToDisclosure float64, currency string) (WorkLogs, int) {
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

			cs.mapEventAssigned(event, issueClosedAt, workLogs, currency)

			// Increment fix count if not yet done
			if isIncremented := areIncremented[event.Assignee.Login]; !isIncremented {
				cs.incrementFixCounters(event.Assignee.Login, timeToDisclosure, severity)
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
func (cs contributors) mapEventAssigned(event github.IssueEvent, issueClosedAt time.Time, workLogs WorkLogs, currency string) {
	cs.mapAssigneeIfMissing(*event.Assignee, currency)

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
func (cs contributors) mapAssigneeIfMissing(assignee github.User, currency string) {
	_, ok := cs[assignee.Login]
	if !ok {
		cs[assignee.Login] = &contributor{
			Login:            assignee.Login,
			AvatarURL:        assignee.AvatarURL,
			HTMLURL:          assignee.HTMLURL,
			Rewards:          []rewardEvent{},
			Currency:         currency,
			TimeToDisclosure: timeToDisclosure{},
			Severities:       map[github.IssueSeverity]int{},
			RewardsLastYear:  newRewardsLastYear(time.Now()),
		}
	}
}

// updateFixCounters updates the fix counters of the contributor who is assigned to the issue in the contributors' map.
func (cs contributors) incrementFixCounters(login string, timeToDisclosure float64, severity github.IssueSeverity) {
	contributor := cs[login]
	contributor.incrementFixCounters(timeToDisclosure, severity)
}

// updateFixCounters updates the fix counters of the contributor who is assigned to the issue in the contributors' map.
func (c *contributor) incrementFixCounters(timeToDisclosure float64, severity github.IssueSeverity) {
	// Increment fix count
	c.FixCount++
	// Increment severity counter
	counterSeverities := c.Severities[severity]
	c.Severities[severity] = counterSeverities + 1
	// Append time to disclosure
	c.TimeToDisclosure.Time = append(c.TimeToDisclosure.Time, timeToDisclosure)
}

// updateMeanAndDeviationOfDisclosure updates the mean and deviation of the time to disclosure of all contributors.
func (cs contributors) updateMeanAndDeviationOfDisclosure() {
	for _, contributor := range cs {
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
func (cs contributors) updateAverageSeverity() {
	for _, contributor := range cs {
		if contributor.FixCount == 0 {
			continue
		}

		contributor.MeanSeverity = (2*float64(contributor.Severities[github.Low]) +
			5.5*float64(contributor.Severities[github.Medium]) +
			9*float64(contributor.Severities[github.High]) +
			9.5*float64(contributor.Severities[github.Critical])) / float64(contributor.FixCount)
	}
}

func (cs contributors) toSortedSlice() []*contributor {
	contributorsSlice := cs.toSlice()
	sortContributors(contributorsSlice)
	return contributorsSlice
}

// mapToSlice transforms the contributors map to a contributors slice.
func (cs contributors) toSlice() []*contributor {
	contributorsSlice := make([]*contributor, 0)
	for _, contributor := range cs {
		contributorsSlice = append(contributorsSlice, contributor)
	}

	return contributorsSlice
}

// sortContributors sorts the contributors by descending reward sum.
func sortContributors(contributors []*contributor) {
	c := collate.New(language.Und, collate.IgnoreCase)
	sort.SliceStable(contributors, func(i, j int) bool {
		if contributors[i].RewardSum == contributors[j].RewardSum {
			return c.CompareString(contributors[i].Login, contributors[j].Login) == -1
		}
		return contributors[i].RewardSum > contributors[j].RewardSum
	})
}
