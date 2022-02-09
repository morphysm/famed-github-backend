package famed

import (
	"log"
	"math"
	"time"

	"github.com/google/go-github/v41/github"
	"github.com/morphysm/famed-github-backend/internal/client/installation"
)

type Contributors map[string]*Contributor

type Contributor struct {
	Login            string                `json:"login"`
	AvatarURL        *string               `json:"avatarUrl,omitempty"`
	HTMLURL          *string               `json:"htmlUrl,omitempty"`
	GravatarID       *string               `json:"gravatarId,omitempty"`
	FixCount         int                   `json:"fixCount,omitempty"`
	Rewards          []Reward              `json:"rewards"`
	RewardSum        float64               `json:"rewardSum"`
	Currency         string                `json:"currency"`
	RewardsLastYear  RewardsLastYear       `json:"rewardsLastYear,omitempty"`
	TimeToDisclosure TimeToDisclosure      `json:"timeToDisclosure"`
	Severities       map[IssueSeverity]int `json:"severities"`
	MeanSeverity     float64               `json:"meanSeverity"`
	// For issue comment generation
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
	currency     string
	rewards      map[IssueSeverity]float64
	usdToEthRate float64
}

// contributors creates a contributors map based on an array of issue and a map of event arrays.
func (r *repo) Contributors() Contributors {
	// Map issues and events to contributors
	r.issuesAndEventsToContributors()
	// Calculate mean and deviation of time to disclosure
	r.contributors.updateMeanAndDeviationOfDisclosure()
	// Calculate average severity of fixed issues
	r.contributors.updateAverageSeverity()

	return r.contributors
}

func (r *repo) issuesAndEventsToContributors() {
	r.contributors = Contributors{}
	for issueID, issue := range r.issues {
		// Map issue to contributors
		err := r.contributors.MapIssue(issue, BoardOptions{
			currency:     r.config.Currency,
			rewards:      r.config.Rewards,
			usdToEthRate: r.ethRate,
		})
		if err != nil {
			log.Printf("[contributors] error while mapping issue with ID: %d, error: %v", issue.Issue.ID, err)
			issue.Error = err
			r.issues[issueID] = issue
		}
	}
}

// MapIssue updates the contributors map based on a set of events and an issue.
// TODO investigate if different data handling for rewards works
func (contributors Contributors) MapIssue(issue Issue, boardOptions BoardOptions) error {
	var (
		workLogs         = WorkLogs{}
		reopenCount      = 0
		issueCreatedAt   = *issue.Issue.CreatedAt
		issueClosedAt    = *issue.Issue.ClosedAt
		timeToDisclosure = issueClosedAt.Sub(issueCreatedAt).Minutes()
	)

	// Read severity from issue
	severity, err := issue.severity()
	if err != nil {
		log.Printf("[MapIssue] no valid label found for issue with ID: %d and label error: %v", issue.Issue.ID, err)
		return err
	}

	// Get severity reward from config
	severityReward := boardOptions.rewards[severity]

	// Increment fix count for all assignees assigned to the closed issue
	for _, assignee := range issue.Issue.Assignees {
		// Add on closed assignee
		contributors.mapAssigneeIfMissing(assignee, boardOptions.currency)
		// Increment fix counter only for assignees on closed
		contributors.updateFixCounters(assignee, timeToDisclosure, severity)
	}

	// Iterate through issue events and map events if event type is of interest
	for _, event := range issue.Events {
		switch *event.Event {
		case string(installation.IssueEventActionAssigned):
			if _, err = isIssueUnAssignedEventDataValid(event); err != nil {
				log.Printf("[MapIssue] event assigned data is invalid for event with ID: %d, err: %v", event.ID, err)
				continue
			}
			contributors.mapEventAssigned(event, issueClosedAt, workLogs, boardOptions.currency)
		case string(installation.IssueEventActionUnassigned):
			if _, err = isIssueUnAssignedEventDataValid(event); err != nil {
				log.Printf("[MapIssue] event unassigened data is invalid for event with ID: %d, err: %v", event.ID, err)
				continue
			}
			mapEventUnassigned(event, workLogs)
		case string(installation.IssueEventActionReopened):
			reopenCount++
		}
	}

	// Calculate the reward
	contributors.updateReward(workLogs, issueCreatedAt, issueClosedAt, reopenCount, severityReward, boardOptions.usdToEthRate)

	return nil
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
			Currency:         rewardUnit,
			TimeToDisclosure: TimeToDisclosure{},
			Severities:       map[IssueSeverity]int{},
			RewardsLastYear:  newRewardsLastYear(time.Now()),
		}
	}
}

// updateFixCounters updates the fix counters of the contributor who is assigned to the issue in the contributors' map.
func (contributors Contributors) updateFixCounters(assignee *github.User, timeToDisclosure float64, severity IssueSeverity) {
	contributor, _ := contributors[*assignee.Login]

	// Increment fix count
	contributor.FixCount++
	// Increment severity counter
	counterSeverities := contributor.Severities[severity]
	contributor.Severities[severity] = counterSeverities + 1
	// Append time to disclosure
	contributor.TimeToDisclosure.Time = append(contributor.TimeToDisclosure.Time, timeToDisclosure)
}

// mapEventAssigned handles an assigned event, updating the contributor map.
func (contributors Contributors) mapEventAssigned(event *github.IssueEvent, issueClosedAt time.Time, workLogs WorkLogs, rewardUnit string) {
	if event.CreatedAt.After(issueClosedAt) {
		return
	}

	contributors.mapAssigneeIfMissing(event.Assignee, rewardUnit)

	// Append work log
	workLogs.Add(*event.Assignee.Login, WorkLog{*event.CreatedAt, issueClosedAt})
}

// mapEventUnassigned handles an unassigned event, updating the work log of the unassigned contributor.
func mapEventUnassigned(event *github.IssueEvent, workLogs WorkLogs) {
	err := workLogs.UpdateEnd(*event.Assignee.Login, *event.CreatedAt)
	if err != nil {
		log.Printf("[mapEventUnassigned] %v on map of event with id %d \n", err, event.ID)
	}
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
