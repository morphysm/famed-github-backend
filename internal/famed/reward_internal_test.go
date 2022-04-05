package famed

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

//nolint:funlen
func TestUpdateReward(t *testing.T) {
	t.Parallel()

	exampleTime := time.Date(2022, 4, 4, 0, 0, 0, 0, time.UTC)
	testCases := []struct {
		Name         string
		Contributors Contributors
		WorkLogs     map[string][]WorkLog
		Open         time.Time
		Close        time.Time
		K            int
		Expected     Contributors
	}{
		{
			Name:         "ContributorsFromIssues nil",
			Contributors: nil,
			Expected:     nil,
		},
		{
			Name:         "ContributorsFromIssues empty",
			Contributors: Contributors{},
			Expected:     Contributors{},
		},
		{
			Name:         "Contributor empty",
			Contributors: Contributors{"TestUser": {}},
			Expected:     Contributors{"TestUser": {}},
		},
		{
			Name: "Contributor without work log",
			Contributors: Contributors{"TestUser": {
				Login:            "TestUser",
				AvatarURL:        "nil",
				HTMLURL:          "nil",
				FixCount:         0,
				RewardsLastYear:  nil,
				Rewards:          nil,
				RewardSum:        0,
				TimeToDisclosure: TimeToDisclosure{},
				Severities:       nil,
				MeanSeverity:     0,
			}},
			Expected: Contributors{"TestUser": {
				Login:            "TestUser",
				AvatarURL:        "nil",
				HTMLURL:          "nil",
				FixCount:         0,
				RewardsLastYear:  nil,
				Rewards:          nil,
				RewardSum:        0,
				TimeToDisclosure: TimeToDisclosure{},
				Severities:       nil,
				MeanSeverity:     0,
			}},
		},
		{
			Name: "Contributor with empty work log",
			Contributors: Contributors{"TestUser": {
				Login:            "TestUser",
				AvatarURL:        "nil",
				HTMLURL:          "nil",
				FixCount:         0,
				RewardsLastYear:  nil,
				Rewards:          nil,
				RewardSum:        0,
				TimeToDisclosure: TimeToDisclosure{},
				Severities:       nil,
				MeanSeverity:     0,
			}},
			WorkLogs: map[string][]WorkLog{"TestUser": {}},
			Expected: Contributors{"TestUser": {
				Login:            "TestUser",
				AvatarURL:        "nil",
				HTMLURL:          "nil",
				FixCount:         0,
				RewardsLastYear:  nil,
				Rewards:          nil,
				RewardSum:        0,
				TimeToDisclosure: TimeToDisclosure{},
				Severities:       nil,
				MeanSeverity:     0,
			}},
		},
		{
			Name: "Contributor with 0 duration work log",
			Contributors: Contributors{"TestUser": {
				Login:            "TestUser",
				AvatarURL:        "nil",
				HTMLURL:          "nil",
				FixCount:         0,
				RewardsLastYear:  nil,
				Rewards:          nil,
				RewardSum:        0,
				TimeToDisclosure: TimeToDisclosure{},
				Severities:       nil,
				MeanSeverity:     0,
			}},
			WorkLogs: map[string][]WorkLog{"TestUser": {{exampleTime, exampleTime}}},
			Expected: Contributors{"TestUser": {
				Login:            "TestUser",
				AvatarURL:        "nil",
				HTMLURL:          "nil",
				FixCount:         0,
				RewardsLastYear:  nil,
				Rewards:          nil,
				RewardSum:        0,
				TimeToDisclosure: TimeToDisclosure{},
				Severities:       nil,
				MeanSeverity:     0,
			}},
		},
		{
			Name: "Contributor with 1 day duration work log",
			Contributors: Contributors{"TestUser": {
				Login:            "TestUser",
				AvatarURL:        "nil",
				HTMLURL:          "nil",
				FixCount:         0,
				RewardsLastYear:  newRewardsLastYear(exampleTime.Add(time.Hour * 24)),
				Rewards:          nil,
				RewardSum:        0,
				TimeToDisclosure: TimeToDisclosure{},
				Severities:       nil,
				MeanSeverity:     0,
			}},
			WorkLogs: map[string][]WorkLog{"TestUser": {{exampleTime, exampleTime.Add(time.Hour * 24)}}},
			Open:     exampleTime,
			Close:    exampleTime.Add(time.Hour * 24),
			Expected: Contributors{"TestUser": {
				Login:     "TestUser",
				AvatarURL: "nil",
				HTMLURL:   "nil",
				FixCount:  0,
				RewardsLastYear: RewardsLastYear{
					{Month: "4.2022", Reward: 0.975},
					{Month: "3.2022", Reward: 0},
					{Month: "2.2022", Reward: 0},
					{Month: "1.2022", Reward: 0},
					{Month: "12.2021", Reward: 0},
					{Month: "11.2021", Reward: 0},
					{Month: "10.2021", Reward: 0},
					{Month: "9.2021", Reward: 0},
					{Month: "8.2021", Reward: 0},
					{Month: "7.2021", Reward: 0},
					{Month: "6.2021", Reward: 0},
					{Month: "5.2021", Reward: 0},
				},
				Rewards:          []Reward{{Date: exampleTime.Add(time.Hour * 24), Reward: 0.975}},
				RewardSum:        0.975,
				TimeToDisclosure: TimeToDisclosure{},
				Severities:       nil,
				MeanSeverity:     0,
				TotalWorkTime:    86400000000000,
			}},
		},
	}

	for _, tC := range testCases {
		tC := tC
		t.Run(tC.Name, func(t *testing.T) {
			t.Parallel()
			// WHEN
			tC.Contributors.updateRewards(tC.WorkLogs, tC.Open, tC.Close, tC.K, 40, 1)

			// THEN
			assert.Equal(t, tC.Expected, tC.Contributors)
		})
	}
}

func TestReward(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name     string
		T        time.Duration
		K        int
		Expected float64
	}{
		{
			Name:     "Instant fix",
			T:        time.Duration(0),
			K:        0,
			Expected: 1,
		},
		{
			Name:     "Instant fix with reopens",
			T:        time.Duration(0),
			K:        10,
			Expected: 1,
		},
		{
			Name:     "Fix after 40 days",
			T:        time.Hour * 24 * 40,
			K:        0,
			Expected: 0.0,
		},
		{
			Name:     "Fix after 40 days with reopens",
			T:        time.Hour * 24 * 40,
			K:        10,
			Expected: 0.0,
		},
		{
			Name:     "Fix after 20 days",
			T:        time.Hour * 24 * 20,
			K:        0,
			Expected: 0.5,
		},
		{
			Name:     "Fix after 20 days with reopens",
			T:        time.Hour * 24 * 20,
			K:        2,
			Expected: 0.03125,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			// WHEN
			severityResult := reward(testCase.T, testCase.K, 40)

			// THEN
			assert.Equal(t, testCase.Expected, severityResult)
		})
	}
}
