package kudo

import (
	"log"
	"math"
	"time"
)

// updateReward returns the base for each contributor based on
// open (time when issue was opened)
// close (time issue was closed)
// k (number of times the issue was reopened)
// contributors (array of contributors with timeOnIssue)
func updateReward(contributors map[string]*Contributor, workLogs map[string][]WorkLog, open time.Time, closed time.Time, k int) map[string]*Contributor {
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
			contributorTotalWork := totalWork[contributor.Login]
			totalWork[contributor.Login] = contributorTotalWork + work.End.Sub(work.Start)
		}

		contributorTotalWork := totalWork[contributor.Login]
		workSum += contributorTotalWork
	}

	// divide base reward based on percentage of each contributor
	for _, contributor := range contributors {
		contributorTotalWork := totalWork[contributor.Login]
		// < is a safety measure, should not happen
		if contributorTotalWork == 0 {
			continue
		}

		if contributorTotalWork < 0 {
			log.Printf("negative contributor total work: %d\n", contributorTotalWork)
		}

		// calculated share of reward
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
	reward := math.Pow(1.0-t.Hours()/(40*24), 2*float64(k)+1)
	return math.Max(0, reward)
}
