package model

import "time"

type BoardOptions struct {
	Currency        string
	RewardStructure RewardStructure
	Now             time.Time
}

func NewBoardOptions(currency string, rewardStructure RewardStructure, now time.Time) BoardOptions {
	return BoardOptions{
		Currency:        currency,
		RewardStructure: rewardStructure,
		Now:             now,
	}
}
