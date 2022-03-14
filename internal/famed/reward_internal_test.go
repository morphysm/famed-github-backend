package famed

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

//nolint:funlen
func TestUpdateReward(t *testing.T) {
	t.Parallel()

	exampleTime := time.Now()
	testCases := []struct {
		Name         string
		Contributors Contributors
		WorkLogs     map[string][]WorkLog
		Open         time.Time
		Closed       time.Time
		K            int
		Expected     Contributors
	}{
		{
			Name:         "ContributorsForIssues nil",
			Contributors: nil,
			Expected:     nil,
		},
		{
			Name:         "ContributorsForIssues empty",
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
				AvatarURL:        nil,
				HTMLURL:          nil,
				GravatarID:       nil,
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
				AvatarURL:        nil,
				HTMLURL:          nil,
				GravatarID:       nil,
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
				AvatarURL:        nil,
				HTMLURL:          nil,
				GravatarID:       nil,
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
				AvatarURL:        nil,
				HTMLURL:          nil,
				GravatarID:       nil,
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
				AvatarURL:        nil,
				HTMLURL:          nil,
				GravatarID:       nil,
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
				AvatarURL:        nil,
				HTMLURL:          nil,
				GravatarID:       nil,
				FixCount:         0,
				RewardsLastYear:  nil,
				Rewards:          nil,
				RewardSum:        0,
				TimeToDisclosure: TimeToDisclosure{},
				Severities:       nil,
				MeanSeverity:     0,
			}},
		},
	}

	for _, tC := range testCases {
		tC := tC
		t.Run(tC.Name, func(t *testing.T) {
			t.Parallel()
			// WHEN
			tC.Contributors.updateRewards(tC.WorkLogs, tC.Open, tC.Closed, tC.K, 1)

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
			severityResult := reward(testCase.T, testCase.K)

			// THEN
			assert.Equal(t, testCase.Expected, severityResult)
		})
	}
}
