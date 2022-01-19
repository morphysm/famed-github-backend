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
func (contributors Contributors) updateReward(workLogs map[string][]WorkLog, open time.Time, closed time.Time, k int, severityReward float64, usdToEthRate float64) {
	baseReward := reward(closed.Sub(open), k)
	ethReward := rewardToEth(baseReward, severityReward, usdToEthRate)
	// TotalWork maps contributor login to contributor total work
	totalWork := map[string]time.Duration{}

	// Calculate total work time of all contributors
	workSum := time.Duration(0)
	for login, workLog := range workLogs {
		// Calculate total work time of a contributor
		for _, work := range workLog {
			contributorTotalWork := totalWork[login]
			totalWork[login] = contributorTotalWork + work.End.Sub(work.Start)
		}

		// Update work total work time of issue
		contributorTotalWork := totalWork[login]
		workSum += contributorTotalWork
	}

	// Divide base reward based on percentage of each contributor
	for login, contributorTotalWork := range totalWork {
		if contributorTotalWork == 0 {
			continue
		}
		contributor := contributors[login]

		// < is a safety measure, should not happen
		if contributorTotalWork < 0 {
			log.Printf("negative contributor total work: %d\n", contributorTotalWork)
		}

		// Calculated share of reward
		reward := ethReward * float64(contributorTotalWork) / float64(workSum)

		// Updated reward sum
		contributor.RewardSum += reward

		// Update rewards list
		contributor.Rewards = append(contributor.Rewards, Reward{
			Date:   closed,
			Reward: reward,
		})

		// Update reward by month
		monthCount := contributor.MonthlyRewards[closed.Month()]
		monthCount += reward
		contributor.MonthlyRewards[closed.Month()] = monthCount
	}
}

// rewardToEth returns the base reward multiplied by the severity reward and changed to eth.
func rewardToEth(baseReward float64, severityReward float64, usdToEthRate float64) float64 {
	return baseReward * severityReward * usdToEthRate
}

// reward returns the base reward for t (time the issue was open) and k (number of times the issue was reopened).
func reward(t time.Duration, k int) float64 {
	// 1 - t (in days) / 40 ^ 2*k+1
	reward := math.Pow(1.0-t.Hours()/(40*24), 2*float64(k)+1)
	return math.Max(0, reward)
}
