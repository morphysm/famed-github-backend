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

func NewRewardsLastYear(timeStart time.Time) RewardsLastYear {
	rewardsLastYear := make([]MonthlyReward, 12)
	for i := 0; i < 12; i++ {
		timeInMonth := timeStart.AddDate(0, i, 0)
		year, month, _ := timeInMonth.Date()
		rewardsLastYear[i].Month = fmt.Sprintf("%d.%d", month, year)
	}

	return rewardsLastYear
}

func lastCurrentOfMonth() time.Time {
	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	return lastOfMonth
}

func isLessThenAYearAndThisMonthAgo(time time.Time) bool {
	return lastCurrentOfMonth().AddDate(-1, 0, 0).Sub(time) > 0
}
