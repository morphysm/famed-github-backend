package kudo

import (
	"log"
	"math"
	"time"

	"github.com/google/go-github/v41/github"

	"github.com/morphysm/kudos-github-backend/internal/client/installation"
)

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

func EventsToContributors(contributors map[string]*Contributor, events []*github.IssueEvent, issueCreatedAt time.Time, issueClosedAt time.Time, severity IssueSeverity) map[string]*Contributor {
	var (
		workLogs         = map[string][]WorkLog{}
		timeToDisclosure = issueClosedAt.Sub(issueCreatedAt).Minutes()
	)

	if contributors == nil {
		contributors = map[string]*Contributor{}
	}

	for _, event := range events {
		if event.Event == nil {
			continue
		}

		switch *event.Event {
		case string(installation.IssueEventActionAssigned):
			handleEventAssigned(contributors, event, issueClosedAt, timeToDisclosure, severity, workLogs)
		case string(installation.IssueEventActionUnassigned):
			handleEventUnassigned(event, workLogs)
		}
	}

	// Calculate the reward // TODO incorrect because it does not count multiple issues
	contributors = updateReward(contributors, workLogs, issueCreatedAt, issueClosedAt, 0)
	contributors = updateMeanAndDeviationOfDisclosure(contributors)
	contributors = updateAverageSeverity(contributors)

	return contributors
}

func handleEventAssigned(contributors map[string]*Contributor, event *github.IssueEvent, issueClosedAt time.Time, timeToDisclosure float64, severity IssueSeverity, workLogs map[string][]WorkLog) {
	if event.Assignee == nil || event.Assignee.Login == nil || event.CreatedAt == nil {
		return
	}

	if event.CreatedAt.After(issueClosedAt) {
		return
	}

	contributor, ok := contributors[*event.Assignee.Login]
	if !ok {
		contributor = &Contributor{
			Login:            *event.Assignee.Login,
			AvatarURL:        event.Assignee.AvatarURL,
			HTMLURL:          event.Assignee.HTMLURL,
			GravatarID:       event.Assignee.GravatarID,
			Rewards:          []Reward{},
			TimeToDisclosure: TimeToDisclosure{},
			IssueSeverities:  map[IssueSeverity]int{},
			MonthFixCount:    map[time.Month]int{},
		}
	}

	// Increment fix count
	contributor.FixCount++
	monthCount, _ := contributor.MonthFixCount[issueClosedAt.Month()]
	monthCount++
	contributor.MonthFixCount[issueClosedAt.Month()] = monthCount

	// Increment severity counter
	counterSeverities, _ := contributor.IssueSeverities[severity]
	contributor.IssueSeverities[severity] = counterSeverities + 1

	// Append time to disclosure
	contributor.TimeToDisclosure.Time = append(contributor.TimeToDisclosure.Time, timeToDisclosure)

	// Append work log
	// TODO check if work end works like this
	work := WorkLog{Start: *event.CreatedAt, End: issueClosedAt}
	assigneeWorkLogs, _ := workLogs[*event.Assignee.Login]
	assigneeWorkLogs = append(assigneeWorkLogs, work)
	workLogs[*event.Assignee.Login] = assigneeWorkLogs

	contributors[*event.Assignee.Login] = contributor
}

func handleEventUnassigned(event *github.IssueEvent, workLogs map[string][]WorkLog) {
	if event.Assignee == nil || event.Assignee.Login == nil || event.CreatedAt == nil {
		return
	}

	// Append work log
	assigneeWorkLogs, _ := workLogs[*event.Assignee.Login]
	if assigneeWorkLogs == nil {
		log.Printf("no work log on event unassigned")
		return
	}

	assigneeWorkLogs[len(assigneeWorkLogs)-1].End = *event.CreatedAt
	workLogs[*event.Assignee.Login] = assigneeWorkLogs
}

// TODO make the functions more efficient
// TODO can we only pas TimeToDisclosure?
func updateMeanAndDeviationOfDisclosure(contributors map[string]*Contributor) map[string]*Contributor {
	for _, contributor := range contributors {
		if contributor.FixCount == 0 {
			continue
		}

		var totalTime, sd float64
		// Calculate mean
		for _, time := range contributor.TimeToDisclosure.Time {
			totalTime += time
		}

		contributor.TimeToDisclosure.Mean = totalTime / float64(contributor.FixCount)

		// Calculate standard deviation
		for _, time := range contributor.TimeToDisclosure.Time {
			// The use of Pow math function func Pow(x, y float64) float64
			sd += math.Pow(time-contributor.TimeToDisclosure.Mean, 2)
		}

		contributor.TimeToDisclosure.StandardDeviation = math.Sqrt(sd / float64(contributor.FixCount))
	}

	return contributors
}

func updateAverageSeverity(contributors map[string]*Contributor) map[string]*Contributor {
	for _, contributor := range contributors {
		if contributor.FixCount == 0 {
			continue
		}

		contributor.MeanSeverity = (2*float64(contributor.IssueSeverities[IssueSeverityLow]) +
			5.5*float64(contributor.IssueSeverities[IssueSeverityMedium]) +
			9*float64(contributor.IssueSeverities[IssueSeverityHigh]) +
			9.5*float64(contributor.IssueSeverities[IssueSeverityCritical])) / float64(contributor.FixCount)
	}

	return contributors
}
