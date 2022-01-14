package kudo

import "fmt"

func GenerateCommentFromContributors(contributors map[string]*Contributor) string {
	if len(contributors) > 0 {
		comment := "Kudo suggests:"
		for _, contributor := range contributors {
			comment = fmt.Sprintf("%s\n Contributor: %s, Reward: %f\n", comment, contributor.Login, contributor.RewardSum)
		}
		return comment
	}

	return "Kudo could not find valid contributors."
}
