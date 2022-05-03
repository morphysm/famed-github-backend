package comment

import (
	"errors"
	"fmt"
	"strings"

	model2 "github.com/morphysm/famed-github-backend/internal/famed/model"
)

var ErrNoContributors = errors.New("GitHub data incomplete")

const RewardCommentTableHeader = "| Contributor | Time | Reward |\n| ----------- | ----------- | ----------- |"

type RewardComment struct {
	identifier Identifier
	headline   string
	table      string
}

// NewRewardComment return a RewardComment.
func NewRewardComment(contributors []*model2.Contributor, currency string, owner string, repoName string) RewardComment {
	rewardComment := RewardComment{}
	// TODO pass version information
	rewardComment.identifier = NewIdentifier(RewardCommentType, "TODO")

	for _, contributor := range contributors {
		rewardComment.headline = fmt.Sprintf("%s@%s ", rewardComment.headline, contributor.Login)
	}

	rewardComment.headline = fmt.Sprintf("%s- you Got Famed! 💎 Check out your new score here: https://www.famed.morphysm.com/teams/%s/%s", rewardComment.headline, owner, repoName)
	rewardComment.table = RewardCommentTableHeader

	for _, contributor := range contributors {
		rewardComment.table = fmt.Sprintf("%s\n|%s|%s|%d %s|", rewardComment.table, contributor.Login, contributor.TotalWorkTime, int(contributor.RewardSum), currency)
	}

	return rewardComment
}

func (c RewardComment) String() (string, error) {
	var sb strings.Builder

	identifier, err := c.identifier.String()
	if err != nil {
		return "", err
	}

	sb.WriteString(identifier)
	sb.WriteString("\n")
	sb.WriteString(c.headline)
	sb.WriteString("\n")
	sb.WriteString(c.table)

	return sb.String(), nil
}

func (c RewardComment) Type() Type {
	return c.identifier.Type
}
