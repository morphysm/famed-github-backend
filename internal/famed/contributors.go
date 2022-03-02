package famed

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/google/go-github/v41/github"
	"github.com/morphysm/famed-github-backend/internal/client/installation"
	"github.com/morphysm/famed-github-backend/internal/config"
)

type Contributors map[string]*Contributor

type Contributor struct {
	Login            string                       `json:"login"`
	AvatarURL        *string                      `json:"avatarUrl,omitempty"`
	HTMLURL          *string                      `json:"htmlUrl,omitempty"`
	GravatarID       *string                      `json:"gravatarId,omitempty"`
	FixCount         int                          `json:"fixCount,omitempty"`
	Rewards          []Reward                     `json:"rewards"`
	RewardSum        float64                      `json:"rewardSum"`
	Currency         string                       `json:"currency"`
	RewardsLastYear  RewardsLastYear              `json:"rewardsLastYear,omitempty"`
	TimeToDisclosure TimeToDisclosure             `json:"timeToDisclosure"`
	Severities       map[config.IssueSeverity]int `json:"severities"`
	MeanSeverity     float64                      `json:"meanSeverity"`
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
	rewards      map[config.IssueSeverity]float64
	usdToEthRate float64
}

// ContributorsForIssue returns a contributors map generated from the repo's internal issue with issueID
// and its corresponding events.
func (r *repo) ContributorsForIssue(issueNumber int) Contributors {
	r.contributors = Contributors{}
	issue := r.issues[issueNumber]
	// Map issue to contributors
	err := r.contributors.MapIssue(issue, BoardOptions{
		currency:     r.config.Currency,
		rewards:      r.config.Rewards,
		usdToEthRate: r.ethRate,
	})
	if err != nil {
		log.Printf("[contributors] error while mapping issue with ID: %d, error: %v", issue.Issue.ID, err)
		issue.Error = err
		r.issues[issueNumber] = issue
	}

	return r.contributors
}

// ContributorsForIssues returns a contributors map generated from the repo's internal issues corresponding events.
func (r *repo) ContributorsForIssues() Contributors {
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
		// areIncremented tracks contributors that have had their fix counters incremented
		areIncremented   = make(map[string]bool)
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

	// Iterate through issue events and map events if event type is of interest
	for _, event := range issue.Events {
		switch *event.Event {
		case string(installation.IssueEventActionAssigned):
			if eventValid, err := isIssueUnAssignedEventDataValid(event); !eventValid {
				log.Printf("[MapIssue] event assigned data is invalid for event with ID: %d, err: %v", event.ID, err)
				continue
			}
			if event.CreatedAt.After(issueClosedAt) {
				continue
			}

			contributors.mapEventAssigned(event, issueClosedAt, workLogs, boardOptions.currency)

			// Increment fix count if not yet done
			if isIncremented := areIncremented[*event.Assignee.Login]; !isIncremented {
				contributors.incrementFixCounters(event.Assignee, timeToDisclosure, severity)
				areIncremented[*event.Assignee.Login] = true
			}
		case string(installation.IssueEventActionUnassigned):
			if eventValid, err := isIssueUnAssignedEventDataValid(event); !eventValid {
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

// mapEventAssigned handles an assigned event, updating the contributor map.
func (contributors Contributors) mapEventAssigned(event *github.IssueEvent, issueClosedAt time.Time, workLogs WorkLogs, currency string) {
	contributors.mapAssigneeIfMissing(event.Assignee, currency)

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

// mapAssigneeIfMissing adds a contributor to the contributors' map if the contributor is missing.
func (contributors Contributors) mapAssigneeIfMissing(assignee *github.User, currency string) {
	_, ok := contributors[*assignee.Login]
	if !ok {
		contributors[*assignee.Login] = &Contributor{
			Login:            *assignee.Login,
			AvatarURL:        assignee.AvatarURL,
			HTMLURL:          assignee.HTMLURL,
			GravatarID:       assignee.GravatarID,
			Rewards:          []Reward{},
			Currency:         currency,
			TimeToDisclosure: TimeToDisclosure{},
			Severities:       map[config.IssueSeverity]int{},
			RewardsLastYear:  newRewardsLastYear(time.Now()),
		}
	}
}

// updateFixCounters updates the fix counters of the contributor who is assigned to the issue in the contributors' map.
func (contributors Contributors) incrementFixCounters(assignee *github.User, timeToDisclosure float64, severity config.IssueSeverity) {
	contributor, _ := contributors[*assignee.Login]
	if contributor == nil {
		fmt.Println(*assignee.Login)
		fmt.Println(contributor)
	}

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
