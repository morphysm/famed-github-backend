package kudo

import (
	"math"
	"time"
)

type IssueSeverity string

const (
	// IssueSeverityNone represents a CVSS of 0
	IssueSeverityNone IssueSeverity = "none"
	// IssueSeverityLow represents a CVSS of 0.1-3.9
	IssueSeverityLow IssueSeverity = "low"
	// IssueSeverityMedium represents a CVSS of 4.0-6.9
	IssueSeverityMedium IssueSeverity = "medium"
	// IssueSeverityHigh represents a CVSS of 7.0-8.9
	IssueSeverityHigh IssueSeverity = "high"
	// IssueSeverityCritical represents a CVSS of 9.0-10.0
	IssueSeverityCritical IssueSeverity = "critical"
)

type WorkLog struct {
	Start time.Time
	End   time.Time
}

type Reward struct {
	Date   time.Time `json:"date"`
	Reward float64   `json:"reward"`
}

type Contributor struct {
	Login            string                `json:"login"`
	AvatarURL        *string               `json:"avatar_url,omitempty"`
	HTMLURL          *string               `json:"html_url,omitempty"`
	GravatarID       *string               `json:"gravatar_id,omitempty"`
	WorkLogs         []WorkLog             `json:"-"`
	Rewards          []Reward              `json:"rewards"`
	RewardSum        float64               `json:"reward_sum"`
	TimeToDisclosure []float64             `json:"time_to_disclosure"`
	IssueSeverities  map[IssueSeverity]int `json:"issue_severity"`
}

type Contributors []*Contributor

// UpdateReward returns the base for each contributor based on
// open (time when issue was opened)
// close (time issue was closed)
// k (number of times the issue was reopened)
// contributors (array of contributors with timeOnIssue)
func UpdateReward(contributors map[string]*Contributor, workLogs map[string][]WorkLog, open time.Time, closed time.Time, k int) map[string]*Contributor {
	baseReward := reward(closed.Sub(open), k)
	// totalWork maps contributor login to contributor total work
	totalWork := map[string]time.Duration{}

	// calculate total work time of all contributors
	workSum := time.Duration(0)
	for _, contributor := range contributors {
		// calculate total work time of a contributor
		contributorWorkLogs, ok := workLogs[contributor.Login]
		if !ok {
			continue
		}

		for _, work := range contributorWorkLogs {
			contributorTotalWork, _ := totalWork[contributor.Login]
			totalWork[contributor.Login] = contributorTotalWork + work.End.Sub(work.Start)
		}

		contributorTotalWork, _ := totalWork[contributor.Login]
		workSum += contributorTotalWork
	}

	// divide base reward based on percentage of each contributor
	for _, contributor := range contributors {
		contributorTotalWork, _ := totalWork[contributor.Login]
		reward := baseReward * float64(workSum) / float64(contributorTotalWork)

		contributor.RewardSum += reward
		contributor.Rewards = append(contributor.Rewards, Reward{
			Date:   closed,
			Reward: reward,
		})
	}

	return contributors
}

// reward returns the base reward for t (time the issue was open) and k (number of times the issue was reopened).
func reward(t time.Duration, k int) float64 {
	// 1 - t (in days) / 40 ^ 2*k+1
	return math.Pow(1.0-t.Hours()/40*24, 2*float64(k)+1)
}
