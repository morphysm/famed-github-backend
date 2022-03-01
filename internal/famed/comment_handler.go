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

	response, err := gH.checkAndUpdateComments(c.Request().Context(), owner, repoName)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response)
}

// checkAndUpdateComments checks all comments and updates comments where necessary in a concurrent fashion.
func (gH *githubHandler) checkAndUpdateComments(ctx context.Context, owner string, repoName string) ([]*IssueCommentUpdate, error) {
	repo := NewRepo(gH.famedConfig, gH.githubInstallationClient, gH.currencyClient, owner, repoName)

	comments, err := repo.GetComments(ctx)
	if err != nil {
		return nil, err
	}

	issueCommentUpdates := make([]*IssueCommentUpdate, len(comments))
	var wg sync.WaitGroup
	i := 0
	for issueNumber, comment := range comments {
		issueCommentUpdates[i] = &IssueCommentUpdate{}
		wg.Add(1)
		go gH.checkAndUpdateComment(ctx, &wg, owner, repoName, issueNumber, comment, issueCommentUpdates[i])
		i++
	}
	wg.Wait()

	return issueCommentUpdates, nil
}

// checkAndUpdateComment should be run as  a go routine to check a comment and update the comment if necessary.
func (gH *githubHandler) checkAndUpdateComment(ctx context.Context, wg *sync.WaitGroup, owner string, repoName string, issueNumber int, comment string, issueCommentUpdate *IssueCommentUpdate) {
	defer wg.Done()
	issueCommentUpdate.IssueNumber = issueNumber

	issueComments, err := gH.githubInstallationClient.GetComments(ctx, owner, repoName, issueNumber)
	if err != nil {
		log.Printf("[UpdateComments] error while getting comments for issue #%d, error: %v", issueNumber, err)
		issueCommentUpdate.Error = err.Error()
		return
	}

	lastCommentByBot := getLastCommentsByUser(issueComments, gH.famedConfig.BotUserID)
	if isCommentValid(lastCommentByBot) && *lastCommentByBot.Body != comment {
		log.Printf("[UpdateComments] updating comment for issue #%d", issueNumber)
		err := gH.githubInstallationClient.PostComment(ctx, owner, repoName, issueNumber, comment)
		if err != nil {
			log.Printf("[UpdateComments] error while posting comment for issue #%d, error: %v", issueNumber, err)
			issueCommentUpdate.Error = err.Error()
			return
		}

		issueCommentUpdate.Updated = true
	}
}

// getLastCommentsByUser returns the last comment made by a user in a list of comments.
func getLastCommentsByUser(issueComments []*github.IssueComment, userID int64) *github.IssueComment {
	for i := len(issueComments) - 1; i >= 0; i-- {
		comment := issueComments[i]

		if comment.User != nil && *comment.User.ID == userID {
			return comment
		}
	}

	return nil
}
