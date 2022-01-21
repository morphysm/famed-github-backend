package kudo

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

func NewRewardsLastYear(timeStart time.Time) RewardsLastYear {
	rewardsLastYear := make([]MonthlyReward, monthsInAYear)
	for i := 0; i < 12; i++ {
		timeInMonth := timeStart.AddDate(0, -i, 0)
		year, month, _ := timeInMonth.Date()
		rewardsLastYear[i].Month = fmt.Sprintf("%d.%d", month, year)
	}

	return rewardsLastYear
}

func lastCurrentOfMonth() time.Time {
	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	firstOfMonth := time.Date(currentYear, currentMonth+1, -1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	return lastOfMonth
}

// isLessThenAYearAndThisMonthAgo returns how many month ago the passed date is and true
// if the month of the passed date is less than 12 months ago.
func isLessThenAYearAndThisMonthAgo(date time.Time) (time.Month, bool) {
	lastOfMonth := lastCurrentOfMonth()
	aYearAgo := lastOfMonth.AddDate(-1, 0, 0)
	if date.Sub(aYearAgo) > 0 {
		return date.Month() - time.Now().Month(), true
	}

	return 0, false
}
