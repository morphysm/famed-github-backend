package model

import (
	"math"
	"time"

	"github.com/morphysm/famed-github-backend/internal/respositories/github/model"
)

type RewardStructure struct {
	severityReward map[model.IssueSeverity]float64
	maxDaysToFix   int
	kMultiplier    int
}

func NewRewardStructure(severityReward map[model.IssueSeverity]float64, maxDaysToFix, kMultiplier int) RewardStructure {
	return RewardStructure{
		severityReward: severityReward,
		maxDaysToFix:   maxDaysToFix,
		kMultiplier:    kMultiplier,
	}
}

// Reward returns the base reward multiplied by the severity reward.
func (RW RewardStructure) Reward(t time.Duration, k int, severity model.IssueSeverity) float64 {
	return RW.baseReward(t, k) * RW.severityReward[severity]
}

// reward returns the base reward for t (time the issue was open) and k (number of times the issue was reopened).
func (RW RewardStructure) baseReward(t time.Duration, k int) float64 {
	// 1 - t (in days) / 40 ^ 2*k+1
	reward := math.Pow(1.0-t.Hours()/float64(RW.maxDaysToFix*24), float64(RW.kMultiplier)*float64(k)+1)
	return math.Max(0, reward)
}
