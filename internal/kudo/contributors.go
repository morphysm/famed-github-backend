package kudo

import (
	"log"
	"time"

	"github.com/google/go-github/v41/github"

	"github.com/morphysm/kudos-github-backend/internal/client/installation"
)

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
			TimeToDisclosure: []float64{},
			IssueSeverities:  map[IssueSeverity]int{},
		}
	}

	// Increment severity counter
	counterSeverities, _ := contributor.IssueSeverities[severity]
	contributor.IssueSeverities[severity] = counterSeverities + 1

	// Append time to disclosure
	contributor.TimeToDisclosure = append(contributor.TimeToDisclosure, timeToDisclosure)

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
