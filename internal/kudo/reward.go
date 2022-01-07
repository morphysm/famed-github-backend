package kudo

import (
	"math"
	"time"
)

// TODO rename
type Work struct {
	Start time.Time
	End   time.Time
}

type Contributor struct {
	Name      string `json:"name"`
	Work      []Work `json:"-"`
	totalWork time.Duration
	Reward    float64 `json:"reward"`
}

type Contributors []*Contributor

// Reward returns the base for each contributor based on
// open (time when issue was opened)
// close (time issue was closed)
// k (number of times the issue was reopened)
// contributors (array of contributors with timeOnIssue)
func (c Contributors) Reward(open time.Time, closed time.Time, k int) {
	baseReward := reward(closed.Sub(open), k)

	// calculate total work time of all contributors
	workSum := time.Duration(0)
	for _, contributor := range c {
		// calculate total work time of a contributor
		for _, work := range contributor.Work {
			contributor.totalWork = contributor.totalWork + work.End.Sub(work.Start)
		}

		workSum += contributor.totalWork
	}

	// divide base reward based on percentage of each contributor
	for _, contributor := range c {
		contributor.Reward = baseReward * float64(workSum) / float64(contributor.totalWork)
	}
}

// reward returns the base reward for t (time the issue was open) and k (number of times the issue was reopened).
func reward(t time.Duration, k int) float64 {
	// 1 - t (in days) / 40 ^ 2*k+1
	return math.Pow(1.0-t.Hours()/40*24, 2*float64(k)+1)
}
