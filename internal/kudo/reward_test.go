package kudo

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUpdateReward(t *testing.T) {
	exampleTime := time.Now()

	testCases := []struct {
		Name         string
		Contributors map[string]*Contributor
		WorkLogs     map[string][]WorkLog
		Open         time.Time
		Closed       time.Time
		K            int
		Expected     map[string]*Contributor
	}{
		{
			Name:         "Contributors nil",
			Contributors: nil,
			Expected:     nil,
		},
		{
			Name:         "Contributors empty",
			Contributors: map[string]*Contributor{},
			Expected:     map[string]*Contributor{},
		},
		{
			Name:         "Contributor empty",
			Contributors: map[string]*Contributor{"TestUser": {}},
			Expected:     map[string]*Contributor{"TestUser": {}},
		},
		{
			Name: "Contributor without work log",
			Contributors: map[string]*Contributor{"TestUser": {
				Login:            "TestUser",
				AvatarURL:        nil,
				HTMLURL:          nil,
				GravatarID:       nil,
				FixCount:         0,
				MonthFixCount:    nil,
				Rewards:          nil,
				RewardSum:        0,
				TimeToDisclosure: TimeToDisclosure{},
				IssueSeverities:  nil,
				MeanSeverity:     0,
			}},
			Expected: map[string]*Contributor{"TestUser": {
				Login:            "TestUser",
				AvatarURL:        nil,
				HTMLURL:          nil,
				GravatarID:       nil,
				FixCount:         0,
				MonthFixCount:    nil,
				Rewards:          nil,
				RewardSum:        0,
				TimeToDisclosure: TimeToDisclosure{},
				IssueSeverities:  nil,
				MeanSeverity:     0,
			}},
		},
		{
			Name: "Contributor with empty work log",
			Contributors: map[string]*Contributor{"TestUser": {
				Login:            "TestUser",
				AvatarURL:        nil,
				HTMLURL:          nil,
				GravatarID:       nil,
				FixCount:         0,
				MonthFixCount:    nil,
				Rewards:          nil,
				RewardSum:        0,
				TimeToDisclosure: TimeToDisclosure{},
				IssueSeverities:  nil,
				MeanSeverity:     0,
			}},
			WorkLogs: map[string][]WorkLog{"TestUser": {}},
			Expected: map[string]*Contributor{"TestUser": {
				Login:            "TestUser",
				AvatarURL:        nil,
				HTMLURL:          nil,
				GravatarID:       nil,
				FixCount:         0,
				MonthFixCount:    nil,
				Rewards:          nil,
				RewardSum:        0,
				TimeToDisclosure: TimeToDisclosure{},
				IssueSeverities:  nil,
				MeanSeverity:     0,
			}},
		},
		{
			Name: "Contributor with 0 duration work log",
			Contributors: map[string]*Contributor{"TestUser": {
				Login:            "TestUser",
				AvatarURL:        nil,
				HTMLURL:          nil,
				GravatarID:       nil,
				FixCount:         0,
				MonthFixCount:    nil,
				Rewards:          nil,
				RewardSum:        0,
				TimeToDisclosure: TimeToDisclosure{},
				IssueSeverities:  nil,
				MeanSeverity:     0,
			}},
			WorkLogs: map[string][]WorkLog{"TestUser": {{exampleTime, exampleTime}}},
			Expected: map[string]*Contributor{"TestUser": {
				Login:            "TestUser",
				AvatarURL:        nil,
				HTMLURL:          nil,
				GravatarID:       nil,
				FixCount:         0,
				MonthFixCount:    nil,
				Rewards:          nil,
				RewardSum:        0,
				TimeToDisclosure: TimeToDisclosure{},
				IssueSeverities:  nil,
				MeanSeverity:     0,
			}},
		},
	}

	for _, tC := range testCases {
		tC := tC
		t.Run(tC.Name, func(t *testing.T) {
			t.Parallel()
			// WHEN
			contributors := updateReward(tC.Contributors, tC.WorkLogs, tC.Open, tC.Closed, tC.K)

			// THEN
			assert.Equal(t, tC.Expected, contributors)
		})
	}
}

func TestReward(t *testing.T) {
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
