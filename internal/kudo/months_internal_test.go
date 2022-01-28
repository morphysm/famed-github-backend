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
		ExpectedMonths time.Month
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
			Name:           "Difference 12 months",
			Now:            time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			Then:           time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			ExpectedMonths: 0,
			ExpectedOk:     false,
		},
		{
			Name:           "Difference 11 months and 30 days",
			Now:            time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			Then:           time.Date(2021, 12, 31, 0, 0, 0, 0, time.UTC),
			ExpectedMonths: 11,
			ExpectedOk:     true,
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
