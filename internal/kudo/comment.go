package kudo

import (
	"errors"
	"fmt"

	"github.com/google/go-github/v41/github"
)

func GenerateComment(issue *github.Issue, events []*github.IssueEvent, currency string, rewards map[IssueSeverity]float64, usdToEthRate float64) string {
	contributors := Contributors{}
	err := contributors.MapIssue(issue, events, currency, rewards, usdToEthRate)
	if err != nil {
		return GenerateCommentFromError(err)
	}

	return contributors.generateCommentFromContributors()
}

func (contributors Contributors) generateCommentFromContributors() string {
	if len(contributors) > 0 {
		comment := "Kudo suggests:"
		for _, contributor := range contributors {
			comment = fmt.Sprintf("%s\n Contributor: %s, Reward: %f\n", comment, contributor.Login, contributor.RewardSum)
		}
		return comment
	}

	return "Kudo could not find valid contributors."
}

func GenerateCommentFromError(err error) string {
	comment := "Kudo could not generate a rewards suggestion. \n" +
		"Reason: "

	if errors.Is(err, ErrIssueMissingAssignee) {
		return fmt.Sprintf("%s The issue is missing an assignee.", comment)
	}

	return fmt.Sprintf("%s Unknown.", comment)
}
