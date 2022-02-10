package famed

import (
	"fmt"
	"time"
)

type RewardsLastYear []MonthlyReward

type MonthlyReward struct {
	Month  string  `json:"month"`
	Reward float64 `json:"reward"`
}

const monthsInAYear = 12

// newRewardsLastYear returns RewardsLastYear with instantiated months starting at the current month and going back 11 months.
func newRewardsLastYear(timeStart time.Time) RewardsLastYear {
	rewardsLastYear := make([]MonthlyReward, monthsInAYear)
	for i := 0; i < 12; i++ {
		timeInMonth := timeStart.AddDate(0, -i, 0)
		year, month, _ := timeInMonth.Date()
		rewardsLastYear[i].Month = fmt.Sprintf("%d.%d", month, year)
	}

	return rewardsLastYear
}

// lastCurrentOfMonth returns the last day of the current month.
func lastCurrentOfMonth(now time.Time) time.Time {
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	return lastOfMonth
}

// isInTheLast12Months returns how many months ago the then date is and
// true if the month of the passed date is less than the current month and 11 months ago.
func isInTheLast12Months(now time.Time, then time.Time) (int, bool) {
	lastOfMonth := lastCurrentOfMonth(now)
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