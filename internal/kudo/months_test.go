package kudo_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/morphysm/kudos-github-backend/internal/kudo"
)

func TestNewRewardsLastYear(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name     string
		Date     time.Time
		Expected kudo.RewardsLastYear
	}{
		{
			Name: "Start on first of 2021",
			Date: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			Expected: kudo.RewardsLastYear{
				{Month: "1.2021", Reward: 0},
				{Month: "12.2020", Reward: 0},
				{Month: "11.2020", Reward: 0},
				{Month: "10.2020", Reward: 0},
				{Month: "9.2020", Reward: 0},
				{Month: "8.2020", Reward: 0},
				{Month: "7.2020", Reward: 0},
				{Month: "6.2020", Reward: 0},
				{Month: "5.2020", Reward: 0},
				{Month: "4.2020", Reward: 0},
				{Month: "3.2020", Reward: 0},
				{Month: "2.2020", Reward: 0},
			},
		},
		{
			Name: "Start mid first of 2020",
			Date: time.Date(2020, 6, 1, 0, 0, 0, 0, time.UTC),
			Expected: kudo.RewardsLastYear{
				{Month: "6.2020", Reward: 0},
				{Month: "5.2020", Reward: 0},
				{Month: "4.2020", Reward: 0},
				{Month: "3.2020", Reward: 0},
				{Month: "2.2020", Reward: 0},
				{Month: "1.2020", Reward: 0},
				{Month: "12.2019", Reward: 0},
				{Month: "11.2019", Reward: 0},
				{Month: "10.2019", Reward: 0},
				{Month: "9.2019", Reward: 0},
				{Month: "8.2019", Reward: 0},
				{Month: "7.2019", Reward: 0},
			},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			// WHEN
			rewardsLastYear := kudo.NewRewardsLastYear(testCase.Date)

			// THEN
			assert.Equal(t, testCase.Expected, rewardsLastYear)
		})
	}
}
