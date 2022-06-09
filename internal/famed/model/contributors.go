package model

import (
	"github.com/phuslu/log"
	"sort"
	"time"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"

	"github.com/morphysm/famed-github-backend/internal/repositories/github/model"
)

type Contributors map[string]*Contributor

// mapAssigneeIfMissing adds a contributor to the contributors' map if the contributor is missing.
func (cs Contributors) mapAssigneeIfMissing(assignee model.User, currency string, now time.Time) {
	_, ok := cs[assignee.Login]
	if !ok {
		cs[assignee.Login] = newContributor(assignee, currency, now)
	}
}

// updateFixCounters updates the fix counters of the contributor who is assigned to the issue in the contributors' map.
func (cs Contributors) incrementFixCounters(login string, timeToDisclosure float64, severity model.IssueSeverity) {
	contributor := cs[login]
	contributor.incrementFixCounters(timeToDisclosure, severity)
}

// updateMeanAndDeviationOfDisclosure updates the mean and deviation of the time to disclosure of all contributors.
func (cs Contributors) updateMeanAndDeviationOfDisclosure() {
	for _, contributor := range cs {
		contributor.updateMeanAndDeviationOfDisclosure()
	}
}

// updateAverageSeverity updates the average severity field of all contributors.
func (cs Contributors) updateAverageSeverity() {
	for _, contributor := range cs {
		contributor.updateAverageSeverity()
	}
}

// UpdateRewards updates the reward for each contributor based on
// open (time when issue was opened)
// close (time issue was closed)
// k (number of times the issue was reopened)
// workLogs (time each contributor worked on the issue)
func (cs Contributors) UpdateRewards(url string, workLogs WorkLogs, open time.Time, close time.Time, k int, severity model.IssueSeverity, boardOptions BoardOptions) {
	points := boardOptions.RewardStructure.Reward(close.Sub(open), k, severity)
	// Get the sum of work per contributor and the total sum of work
	contributorsWork, workSum := workLogs.Sum()

	// Divide base reward based on percentage of each contributor
	for login, contributorTotalWork := range contributorsWork {
		if contributorTotalWork < 0 {
			// < is a safety measure, should not happen
			log.Info().Msgf("contributor total work < 0: %d\n", contributorTotalWork)
			continue
		}
		contributor := cs[login]

		// Assign total work to contributor for issue rewardComment generation
		contributor.TotalWorkTime = contributorTotalWork

		// Calculated share of reward
		// workSum can be 0
		var reward float64
		if workSum == 0 {
			reward = points / float64(len(contributorsWork))
		} else {
			reward = points * float64(contributorTotalWork) / float64(workSum)
		}

		// Updated reward sum
		contributor.updateReward(url, boardOptions.Now, close, reward)
	}
}

func (cs Contributors) toSortedSlice() []*Contributor {
	contributorsSlice := cs.toSlice()
	sortContributors(contributorsSlice)
	return contributorsSlice
}

// mapToSlice transforms the contributors map to a contributors slice.
func (cs Contributors) toSlice() []*Contributor {
	contributorsSlice := make([]*Contributor, 0)
	for _, contributor := range cs {
		contributorsSlice = append(contributorsSlice, contributor)
	}

	return contributorsSlice
}

// sortContributors sorts the contributors by descending reward sum.
func sortContributors(contributors []*Contributor) {
	c := collate.New(language.Und, collate.IgnoreCase)
	sort.SliceStable(contributors, func(i, j int) bool {
		if contributors[i].RewardSum == contributors[j].RewardSum {
			return c.CompareString(contributors[i].Login, contributors[j].Login) == -1
		}
		return contributors[i].RewardSum > contributors[j].RewardSum
	})
}
