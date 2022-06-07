package model

import (
	"math"
	"time"

	"github.com/morphysm/famed-github-backend/internal/repositories/github/model"
)

type Contributor struct {
	Login            string                      `json:"login"`
	AvatarURL        string                      `json:"avatarUrl"`
	HTMLURL          string                      `json:"htmlUrl"`
	FixCount         int                         `json:"fixCount"`
	Rewards          []RewardEvent               `json:"rewards"`
	RewardSum        float64                     `json:"rewardSum"`
	Currency         string                      `json:"currency"`
	RewardsLastYear  RewardsLastYear             `json:"rewardsLastYear"`
	TimeToDisclosure TimeToDisclosure            `json:"timeToDisclosure"`
	Severities       map[model.IssueSeverity]int `json:"severities"`
	MeanSeverity     float64                     `json:"meanSeverity"`
	// For issue rewardComment generation
	TotalWorkTime time.Duration `json:"-"`
}

type TimeToDisclosure struct {
	Time              []float64 `json:"time"`
	Mean              float64   `json:"mean"`
	StandardDeviation float64   `json:"standardDeviation"`
}

type RewardEvent struct {
	Date   time.Time `json:"date"`
	Reward float64   `json:"reward"`
	URL    string    `json:"url"`
}

func newContributor(assignee model.User, currency string, now time.Time) *Contributor {
	return &Contributor{
		Login:            assignee.Login,
		AvatarURL:        assignee.AvatarURL,
		HTMLURL:          assignee.HTMLURL,
		Rewards:          []RewardEvent{},
		Currency:         currency,
		TimeToDisclosure: TimeToDisclosure{},
		Severities:       map[model.IssueSeverity]int{},
		RewardsLastYear:  NewRewardsLastYear(now),
	}
}

// mapIssue maps an issue to a contributor.
func (c *Contributor) mapIssue(url string, reportedDate, publishedDate time.Time, reward float64, severity model.IssueSeverity, now time.Time) {
	// Set reward
	c.updateReward(url, now, publishedDate, reward)

	// Increment fix count
	c.incrementFixCounters(publishedDate.Sub(reportedDate).Minutes(), severity)
}

// updateFixCounters updates the fix counters of the contributor who is assigned to the issue in the contributors' map.
func (c *Contributor) incrementFixCounters(timeToDisclosure float64, severity model.IssueSeverity) {
	// Increment fix count
	c.FixCount++
	// Increment severity counter
	counterSeverities := c.Severities[severity]
	c.Severities[severity] = counterSeverities + 1
	// Append time to disclosure
	c.TimeToDisclosure.Time = append(c.TimeToDisclosure.Time, timeToDisclosure)
}

func (c *Contributor) updateMeanAndDeviationOfDisclosure() {
	if c.FixCount == 0 {
		return
	}

	// Calculate mean
	var totalTime, sd float64
	for _, timeToDisclosure := range c.TimeToDisclosure.Time {
		totalTime += timeToDisclosure
	}

	c.TimeToDisclosure.Mean = totalTime / float64(c.FixCount)

	// Calculate standard deviation
	for _, timeToDisclosure := range c.TimeToDisclosure.Time {
		sd += math.Pow(timeToDisclosure-c.TimeToDisclosure.Mean, 2) //nolint:gomnd
	}

	c.TimeToDisclosure.StandardDeviation = math.Sqrt(sd / float64(c.FixCount))
}

// updateAverageSeverity updates the average severity field of all contributors.
func (c *Contributor) updateAverageSeverity() {
	if c.FixCount == 0 {
		return
	}

	c.MeanSeverity = (2*float64(c.Severities[model.Low]) +
		5.5*float64(c.Severities[model.Medium]) +
		9*float64(c.Severities[model.High]) +
		9.5*float64(c.Severities[model.Critical])) / float64(c.FixCount)

}

func (c *Contributor) updateReward(url string, now, date time.Time, reward float64) {
	// Append reward to reward slice
	c.Rewards = append(c.Rewards, RewardEvent{
		Date:   date,
		Reward: reward,
		URL:    url,
	})

	// Updated reward sum
	c.RewardSum += reward

	// Add reward by month
	if month, ok := isInTheLast12Months(now, date); ok {
		c.RewardsLastYear[month].Reward += reward
	}
}
