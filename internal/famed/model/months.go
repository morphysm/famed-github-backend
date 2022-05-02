package model

import (
	"fmt"
	"time"
)

type RewardsLastYear []monthlyReward

type monthlyReward struct {
	Month  string  `json:"month"`
	Reward float64 `json:"reward"`
}

const monthsInAYear = 12

// NewRewardsLastYear returns rewardsLastYear with instantiated months starting at the current month and going back 11 months.
func NewRewardsLastYear(timeStart time.Time) RewardsLastYear {
	rewardsLastYear := make([]monthlyReward, monthsInAYear)
	year, month, _ := timeStart.Date()
	for i := 0; i < 12; i++ {
		rewardsLastYear[i].Month = fmt.Sprintf("%d.%d", month, year)
		month--
		if month < 1 {
			month = 12
			year--
		}
	}

	return rewardsLastYear
}

// lastDayOfMonth returns the last day of a month of a given time.
func lastDayOfMonth(now time.Time) time.Time {
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)

	return firstOfMonth.AddDate(0, 1, -1)
}

// isInTheLast12Months returns how many months ago the then date is and
// true if the month of the passed date is less than the current month and 11 months ago.
func isInTheLast12Months(now time.Time, then time.Time) (int, bool) {
	lastOfMonth := lastDayOfMonth(now)
	aYearAgo := lastOfMonth.AddDate(-1, 0, 0)
	if then.Sub(aYearAgo) > 0 {
		// Same year
		if now.Year()-then.Year() == 0 {
			return int(now.Month() - then.Month()), true
		}
		// Different year
		monthsTillTheEndOfTheYear := 12 - then.Month()
		return int(monthsTillTheEndOfTheYear + now.Month()), true
	}

	return 0, false
}
