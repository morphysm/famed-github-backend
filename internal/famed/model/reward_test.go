package model_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/morphysm/famed-github-backend/internal/famed/model"
	model2 "github.com/morphysm/famed-github-backend/internal/respositories/github/model"
)

//nolint:funlen
func TestUpdateReward(t *testing.T) {
	t.Parallel()

	exampleTime := time.Date(2022, 4, 4, 0, 0, 0, 0, time.UTC)
	testCases := []struct {
		Name         string
		Contributors model.Contributors
		WorkLogs     map[string][]model.WorkLog
		Open         time.Time
		Close        time.Time
		K            int
		Expected     model.Contributors
	}{
		{
			Name:         "ContributorsFromIssues nil",
			Contributors: nil,
			Expected:     nil,
		},
		{
			Name:         "ContributorsFromIssues empty",
			Contributors: model.Contributors{},
			Expected:     model.Contributors{},
		},
		{
			Name:         "contributor empty",
			Contributors: model.Contributors{"TestUser": {}},
			Expected:     model.Contributors{"TestUser": {}},
		},
		{
			Name: "contributor without work log",
			Contributors: model.Contributors{"TestUser": {
				Login:            "TestUser",
				AvatarURL:        "nil",
				HTMLURL:          "nil",
				FixCount:         0,
				RewardsLastYear:  nil,
				Rewards:          nil,
				RewardSum:        0,
				TimeToDisclosure: model.TimeToDisclosure{},
				Severities:       nil,
				MeanSeverity:     0,
			}},
			Expected: model.Contributors{"TestUser": {
				Login:            "TestUser",
				AvatarURL:        "nil",
				HTMLURL:          "nil",
				FixCount:         0,
				RewardsLastYear:  nil,
				Rewards:          nil,
				RewardSum:        0,
				TimeToDisclosure: model.TimeToDisclosure{},
				Severities:       nil,
				MeanSeverity:     0,
			}},
		},
		{
			Name: "contributor with 0 duration work log",
			Contributors: model.Contributors{"TestUser": {
				Login:            "TestUser",
				AvatarURL:        "nil",
				HTMLURL:          "nil",
				FixCount:         0,
				RewardsLastYear:  nil,
				Rewards:          nil,
				RewardSum:        0,
				TimeToDisclosure: model.TimeToDisclosure{},
				Severities:       nil,
				MeanSeverity:     0,
			}},
			WorkLogs: map[string][]model.WorkLog{"TestUser": {{exampleTime, exampleTime}}},
			Expected: model.Contributors{"TestUser": {
				Login:            "TestUser",
				AvatarURL:        "nil",
				HTMLURL:          "nil",
				FixCount:         0,
				RewardsLastYear:  nil,
				Rewards:          []model.RewardEvent{{Date: time.Time{}, Reward: 1, URL: "TestURL"}},
				RewardSum:        1,
				TimeToDisclosure: model.TimeToDisclosure{},
				Severities:       nil,
				MeanSeverity:     0,
			}},
		},
		{
			Name: "Two contributor with 0 duration work log",
			Contributors: model.Contributors{
				"TestUser1": {
					Login:            "TestUser1",
					AvatarURL:        "nil",
					HTMLURL:          "nil",
					FixCount:         0,
					RewardsLastYear:  nil,
					Rewards:          nil,
					RewardSum:        0,
					TimeToDisclosure: model.TimeToDisclosure{},
					Severities:       nil,
					MeanSeverity:     0,
				},
				"TestUser2": {
					Login:            "TestUser2",
					AvatarURL:        "nil",
					HTMLURL:          "nil",
					FixCount:         0,
					RewardsLastYear:  nil,
					Rewards:          nil,
					RewardSum:        0,
					TimeToDisclosure: model.TimeToDisclosure{},
					Severities:       nil,
					MeanSeverity:     0,
				},
			},
			WorkLogs: map[string][]model.WorkLog{"TestUser1": {{exampleTime, exampleTime}}, "TestUser2": {{exampleTime, exampleTime}}},
			Expected: model.Contributors{
				"TestUser1": {
					Login:            "TestUser1",
					AvatarURL:        "nil",
					HTMLURL:          "nil",
					FixCount:         0,
					RewardsLastYear:  nil,
					Rewards:          []model.RewardEvent{{Date: time.Time{}, Reward: 0.5, URL: "TestURL"}},
					RewardSum:        0.5,
					TimeToDisclosure: model.TimeToDisclosure{},
					Severities:       nil,
					MeanSeverity:     0,
				},
				"TestUser2": {
					Login:            "TestUser2",
					AvatarURL:        "nil",
					HTMLURL:          "nil",
					FixCount:         0,
					RewardsLastYear:  nil,
					Rewards:          []model.RewardEvent{{Date: time.Time{}, Reward: 0.5, URL: "TestURL"}},
					RewardSum:        0.5,
					TimeToDisclosure: model.TimeToDisclosure{},
					Severities:       nil,
					MeanSeverity:     0,
				},
			},
		},
		{
			Name: "Two contributor with 0 and 1 day duration work log",
			Contributors: model.Contributors{
				"TestUser1": {
					Login:            "TestUser1",
					AvatarURL:        "nil",
					HTMLURL:          "nil",
					FixCount:         0,
					RewardsLastYear:  nil,
					Rewards:          nil,
					RewardSum:        0,
					TimeToDisclosure: model.TimeToDisclosure{},
					Severities:       nil,
					MeanSeverity:     0,
				},
				"TestUser2": {
					Login:            "TestUser2",
					AvatarURL:        "nil",
					HTMLURL:          "nil",
					FixCount:         0,
					RewardsLastYear:  nil,
					Rewards:          nil,
					RewardSum:        0,
					TimeToDisclosure: model.TimeToDisclosure{},
					Severities:       nil,
					MeanSeverity:     0,
				},
			},
			WorkLogs: map[string][]model.WorkLog{"TestUser1": {{exampleTime, exampleTime}}, "TestUser2": {{exampleTime, exampleTime.Add(24 * time.Hour)}}},
			Expected: model.Contributors{
				"TestUser1": {
					Login:            "TestUser1",
					AvatarURL:        "nil",
					HTMLURL:          "nil",
					FixCount:         0,
					RewardsLastYear:  nil,
					Rewards:          []model.RewardEvent{{Date: time.Time{}, Reward: 0, URL: "TestURL"}},
					RewardSum:        0,
					TimeToDisclosure: model.TimeToDisclosure{},
					Severities:       nil,
					MeanSeverity:     0,
				},
				"TestUser2": {
					Login:            "TestUser2",
					AvatarURL:        "nil",
					HTMLURL:          "nil",
					FixCount:         0,
					RewardsLastYear:  nil,
					Rewards:          []model.RewardEvent{{Date: time.Time{}, Reward: 1, URL: "TestURL"}},
					RewardSum:        1,
					TimeToDisclosure: model.TimeToDisclosure{},
					Severities:       nil,
					MeanSeverity:     0,
					TotalWorkTime:    24 * time.Hour,
				},
			},
		},
		{
			Name: "contributor with 1 day duration work log",
			Contributors: model.Contributors{"TestUser": {
				Login:            "TestUser",
				AvatarURL:        "nil",
				HTMLURL:          "nil",
				FixCount:         0,
				RewardsLastYear:  model.NewRewardsLastYear(exampleTime.Add(time.Hour * 24)),
				Rewards:          nil,
				RewardSum:        0,
				TimeToDisclosure: model.TimeToDisclosure{},
				Severities:       nil,
				MeanSeverity:     0,
			}},
			WorkLogs: map[string][]model.WorkLog{"TestUser": {{exampleTime, exampleTime.Add(time.Hour * 24)}}},
			Open:     exampleTime,
			Close:    exampleTime.Add(time.Hour * 24),
			Expected: model.Contributors{"TestUser": {
				Login:     "TestUser",
				AvatarURL: "nil",
				HTMLURL:   "nil",
				FixCount:  0,
				RewardsLastYear: model.RewardsLastYear{
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
				Rewards:          []model.RewardEvent{{Date: exampleTime.Add(time.Hour * 24), Reward: 0.975, URL: "TestURL"}},
				RewardSum:        0.975,
				TimeToDisclosure: model.TimeToDisclosure{},
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
			rewardStructure := model.NewRewardStructure(map[model2.IssueSeverity]float64{model2.Low: 1}, 40, 2)
			BoardOptions := model.NewBoardOptions("POINTS", rewardStructure, time.Date(2022, 4, 4, 0, 0, 0, 0, time.UTC))
			tC.Contributors.UpdateRewards("TestURL", tC.WorkLogs, tC.Open, tC.Close, tC.K, model2.Low, BoardOptions)

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
			rewardStructure := model.NewRewardStructure(map[model2.IssueSeverity]float64{model2.Low: 1}, 40, 2)
			severityResult := rewardStructure.Reward(testCase.T, testCase.K, model2.Low)

			// THEN
			assert.Equal(t, testCase.Expected, severityResult)
		})
	}
}
