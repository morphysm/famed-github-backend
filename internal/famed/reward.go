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
func (cs Contributors) updateRewards(workLogs WorkLogs, open time.Time, close time.Time, k int, daysToFix int, severityReward float64) {
	baseReward := reward(close.Sub(open), k, daysToFix)
	points := rewardToPoints(baseReward, severityReward)
	// Get the sum of work per contributor and the total sum of work
	contributorsWork, workSum := workLogs.Sum()

	// Divide base reward based on percentage of each contributor
	for login, contributorTotalWork := range contributorsWork {
		if contributorTotalWork < 0 {
			// < is a safety measure, should not happen
			log.Printf("contributor total work < 0: %d\n", contributorTotalWork)
			continue
		}
		contributor := cs[login]

		// Assign total work to contributor for issue rewardComment generation
		contributor.TotalWorkTime = contributorTotalWork

		// Calculated share of reward
		// workSum can be 0 on
		var reward float64
		if workSum == 0 {
			reward = points / float64(len(contributorsWork))
		} else {
			reward = points * float64(contributorTotalWork) / float64(workSum)
		}

		// Updated reward sum
		contributor.RewardSum += reward

		// Add rewards list
		contributor.Rewards = append(contributor.Rewards, Reward{
			Date:   close,
			Reward: reward,
		})

		// Add reward by month
		if month, ok := isInTheLast12Months(time.Now(), close); ok {
			contributor.RewardsLastYear[month].Reward += reward
		}
	}
}

// rewardToPoints returns the base reward multiplied by the severity reward.
func rewardToPoints(baseReward float64, severityReward float64) float64 {
	return baseReward * severityReward
}

// reward returns the base reward for t (time the issue was open) and k (number of times the issue was reopened).
func reward(t time.Duration, k int, daysToFix int) float64 {
	// 1 - t (in days) / 40 ^ 2*k+1
	reward := math.Pow(1.0-t.Hours()/float64(daysToFix*24), 2*float64(k)+1)
	return math.Max(0, reward)
}
