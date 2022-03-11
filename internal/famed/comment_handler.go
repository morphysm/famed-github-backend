package famed

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/google/go-github/v41/github"
	"github.com/labstack/echo/v4"
)

type IssueCommentUpdate struct {
	IssueNumber int    `json:"issueNumber"`
	Updated     bool   `json:"updated"`
	Error       string `json:"error"`
}

// UpdateComments updates the comments in a GitHub repo.
func (gH *githubHandler) UpdateComments(c echo.Context) error {
	owner := c.Param("owner")
	if owner == "" {
		return echo.ErrBadRequest.SetInternal(ErrMissingOwnerPathParameter)
	}

	repoName := c.Param("repo_name")
	if repoName == "" {
		return echo.ErrBadRequest.SetInternal(ErrMissingRepoPathParameter)
	}

	if installed := gH.githubInstallationClient.CheckInstallation(owner); !installed {
		log.Printf("[Contributors] error on request for contributors: %v", ErrAppNotInstalled)
		return ErrAppNotInstalled
	}

	response, err := gH.compareAndUpdateComments(c.Request().Context(), owner, repoName)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response)
}

// compareAndUpdateComments checks all comments and updates comments where necessary in a concurrent fashion.
func (gH *githubHandler) compareAndUpdateComments(ctx context.Context, owner string, repoName string) ([]*IssueCommentUpdate, error) {
	repo := NewRepo(gH.famedConfig, gH.githubInstallationClient, owner, repoName)

	comments, err := repo.RewardComments(ctx)
	if err != nil {
		return nil, err
	}

	issueCommentUpdates := make([]*IssueCommentUpdate, len(comments))
	var wg sync.WaitGroup
	i := 0
	for issueNumber, comment := range comments {
		issueCommentUpdates[i] = &IssueCommentUpdate{}
		wg.Add(1)
		go gH.compareAndUpdateComment(ctx, &wg, owner, repoName, issueNumber, comment, issueCommentUpdates[i])
		i++
	}
	wg.Wait()

	return issueCommentUpdates, nil
}

// compareAndUpdateComment should be run as  a go routine to check a rewardComment and update the rewardComment if necessary.
func (gH *githubHandler) compareAndUpdateComment(ctx context.Context, wg *sync.WaitGroup, owner string, repoName string, issueNumber int, comment string, issueCommentUpdate *IssueCommentUpdate) {
	defer wg.Done()
	issueCommentUpdate.IssueNumber = issueNumber

	issueComments, err := gH.githubInstallationClient.GetComments(ctx, owner, repoName, issueNumber)
	if err != nil {
		log.Printf("[UpdateComments] error while getting comments for issue #%d, error: %v", issueNumber, err)
		issueCommentUpdate.Error = err.Error()
		return
	}

	lastCommentByBot := findRewardComment(issueComments, gH.famedConfig.BotLogin)
	if lastCommentByBot == nil || (isCommentValid(lastCommentByBot) && *lastCommentByBot.Body != comment) {
		log.Printf("[UpdateComments] updating RewardComment for issue #%d", issueNumber)
		err := gH.githubInstallationClient.PostComment(ctx, owner, repoName, issueNumber, comment)
		if err != nil {
			log.Printf("[UpdateComments] error while posting RewardComment for issue #%d, error: %v", issueNumber, err)
			issueCommentUpdate.Error = err.Error()
			return
		}

		issueCommentUpdate.Updated = true
	}
}

// findRewardComment returns the last RewardComment made by a user in a list of comments.
func findRewardComment(issueComments []*github.IssueComment, botLogin string) *github.IssueComment {
	for i := len(issueComments) - 1; i >= 0; i-- {
		comment := issueComments[i]

		if isUserValid(comment.User) && *comment.User.Login == botLogin && verifyCommentType(*comment.Body, commentReward) {
			return comment
		}
	}

	return nil
}
