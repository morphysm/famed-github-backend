package famed

import (
	"errors"
	"log"
	"time"
)

var ErrNoWorkLogForAssignee = errors.New("no work log found for assignee")

type WorkLog struct {
	Start time.Time
	End   time.Time
}

type WorkLogs map[string][]WorkLog

// Add adds a new WorkLog.
func (wL WorkLogs) Add(login string, workLog WorkLog) {
	assigneeWorkLogs := wL[login]
	assigneeWorkLogs = append(assigneeWorkLogs, workLog)
	wL[login] = assigneeWorkLogs
}

// UpdateEnd updates the end time of the last work log of an assignee.
func (wL WorkLogs) UpdateEnd(login string, end time.Time) error {
	assigneeWorkLogs := wL[login]
	if len(assigneeWorkLogs) == 0 {
		log.Printf("[mapEventUnassigned] %v \n", ErrNoWorkLogForAssignee)
		return ErrNoWorkLogForAssignee
	}

	assigneeWorkLogs[len(assigneeWorkLogs)-1].End = end
	wL[login] = assigneeWorkLogs

	return nil
}

// Sum returns the sum of work per contributor and the total sum of work.
func (wL WorkLogs) Sum() (map[string]time.Duration, time.Duration) {
	// TotalWork maps contributor login to contributor total work
	totalWork := map[string]time.Duration{}

	// Calculate total work time of all contributors
	workSum := time.Duration(0)
	for login, workLog := range wL {
		// Calculate total work time of a contributor
		contributorTotalWork := time.Duration(0)
		for _, work := range workLog {
			contributorTotalWork += work.End.Sub(work.Start)
		}

		totalWork[login] = contributorTotalWork
		// Add work total work time of issue
		workSum += contributorTotalWork
	}

	return totalWork, workSum
}
