package famed

import (
	"log"
	"math"
	"time"
)

// updateRewards updates the reward for each contributor based on
// open (time when issue was opened)
// close (time issue was closed)
// k (number of times the issue was reopened)
// workLogs (time each contributor worked on the issue)
func (contributors Contributors) updateRewards(workLogs WorkLogs, open time.Time, closed time.Time, k int, severityReward float64) {
	baseReward := reward(closed.Sub(open), k)
	points := rewardToPoints(baseReward, severityReward)
	// Get the sum of work per contributor and the total sum of work
	totalWork, workSum := workLogs.Sum()

	// Divide base reward based on percentage of each contributor
	for login, contributorTotalWork := range totalWork {
		if contributorTotalWork <= 0 {
			// < is a safety measure, should not happen
			log.Printf("<= 0 contributor total work: %d\n", contributorTotalWork)
			continue
		}
		contributor := contributors[login]

		// Assign total work to contributor for issue RewardComment generation
		contributor.TotalWorkTime = contributorTotalWork

		// Calculated share of reward
		reward := points * float64(contributorTotalWork) / float64(workSum)

		// Updated reward sum
		contributor.RewardSum += reward

		// Update rewards list
		contributor.Rewards = append(contributor.Rewards, Reward{
			Date:   closed,
			Reward: reward,
		})

		// Update reward by month
		if month, ok := isInTheLast12Months(time.Now(), closed); ok {
			contributor.RewardsLastYear[month].Reward += reward
		}
	}
}

// rewardToPoints returns the base reward multiplied by the severity reward.
func rewardToPoints(baseReward float64, severityReward float64) float64 {
	return baseReward * severityReward
}

// reward returns the base reward for t (time the issue was open) and k (number of times the issue was reopened).
func reward(t time.Duration, k int) float64 {
	// 1 - t (in days) / 40 ^ 2*k+1
	reward := math.Pow(1.0-t.Hours()/(40*24), 2*float64(k)+1)
	return math.Max(0, reward)
}
