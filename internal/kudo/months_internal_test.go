package kudo

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
