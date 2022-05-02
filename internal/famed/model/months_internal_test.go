package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsLessThenAYearAndThisMonthAgo(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name           string
		Now            time.Time
		Then           time.Time
		ExpectedMonths int
		ExpectedOk     bool
	}{
		{
			Name:           "Difference 0",
			Now:            time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			Then:           time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			ExpectedMonths: 0,
			ExpectedOk:     true,
		},
		{
			Name:           "Difference 1 month same year",
			Now:            time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC),
			Then:           time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			ExpectedMonths: 1,
			ExpectedOk:     true,
		},
		{
			Name:           "Difference 1 month different year",
			Now:            time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			Then:           time.Date(2020, 12, 1, 0, 0, 0, 0, time.UTC),
			ExpectedMonths: 1,
			ExpectedOk:     true,
		},
		{
			Name:           "Difference 11 months and 29 days different year",
			Now:            time.Date(2021, 11, 30, 0, 0, 0, 0, time.UTC),
			Then:           time.Date(2020, 12, 1, 0, 0, 0, 0, time.UTC),
			ExpectedMonths: 11,
			ExpectedOk:     true,
		},
		{
			Name:           "Difference 12 months different year",
			Now:            time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			Then:           time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			ExpectedMonths: 0,
			ExpectedOk:     false,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			// WHEN
			months, ok := isInTheLast12Months(testCase.Now, testCase.Then)

			// THEN
			assert.Equal(t, testCase.ExpectedMonths, months)
			assert.Equal(t, testCase.ExpectedOk, ok)
		})
	}
}

func TestNewRewardsLastYear(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name     string
		Date     time.Time
		Expected RewardsLastYear
	}{
		{
			Name: "Start on first of 2021",
			Date: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			Expected: RewardsLastYear{
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
			Expected: RewardsLastYear{
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
		{
			Name: "Start 31.03.2020",
			Date: time.Date(2020, 3, 31, 0, 0, 0, 0, time.UTC),
			Expected: RewardsLastYear{
				{Month: "3.2020", Reward: 0},
				{Month: "2.2020", Reward: 0},
				{Month: "1.2020", Reward: 0},
				{Month: "12.2019", Reward: 0},
				{Month: "11.2019", Reward: 0},
				{Month: "10.2019", Reward: 0},
				{Month: "9.2019", Reward: 0},
				{Month: "8.2019", Reward: 0},
				{Month: "7.2019", Reward: 0},
				{Month: "6.2019", Reward: 0},
				{Month: "5.2019", Reward: 0},
				{Month: "4.2019", Reward: 0},
			},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			// WHEN
			rewardsLastYear := NewRewardsLastYear(testCase.Date)

			// THEN
			assert.Equal(t, testCase.Expected, rewardsLastYear)
		})
	}
}
