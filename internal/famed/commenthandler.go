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
	repoName := c.Param("repo_name")
	if repoName == "" {
		return echo.ErrBadRequest.SetInternal(ErrMissingRepoPathParameter)
	}

	response, err := gH.checkAndUpdateComments(c.Request().Context(), repoName)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response)
}

func (gH *githubHandler) checkAndUpdateComments(ctx context.Context, repoName string) ([]*IssueCommentUpdate, error) {
	repo := NewRepo(gH.famedConfig, gH.githubInstallationClient, gH.currencyClient, repoName)

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
		go gH.checkAndUpdateComment(ctx, &wg, repoName, issueNumber, comment, issueCommentUpdates[i])
		i++
	}
	wg.Wait()

	return issueCommentUpdates, nil
}

func (gH *githubHandler) checkAndUpdateComment(ctx context.Context, wg *sync.WaitGroup, repoName string, issueNumber int, comment string, issueCommentUpdate *IssueCommentUpdate) {
	defer wg.Done()
	issueCommentUpdate.IssueNumber = issueNumber

	issueComments, err := gH.githubInstallationClient.GetComments(ctx, repoName, issueNumber)
	if err != nil {
		log.Printf("[UpdateComments] error while getting comments for issue #%d, error: %v", issueNumber, err)
		issueCommentUpdate.Error = err.Error()
		return
	}
	lastCommentByBot := getLastCommentsByUser(issueComments, 96487857)

	if isCommentValid(lastCommentByBot) && *lastCommentByBot.Body != comment {
		log.Printf("[UpdateComments] updating comment for issue #%d", issueNumber)
		_, err := gH.githubInstallationClient.PostComment(ctx, repoName, issueNumber, comment)
		if err != nil {
			log.Printf("[UpdateComments] error while posting comment for issue #%d, error: %v", issueNumber, err)
			issueCommentUpdate.Error = err.Error()
			return
		}

		issueCommentUpdate.Updated = true
	}
}

func getLastCommentsByUser(issueComments []*github.IssueComment, userID int64) *github.IssueComment {
	for i := len(issueComments) - 1; i >= 0; i-- {
		comment := issueComments[i]

		if comment.User != nil && *comment.User.ID == userID {
			return comment
		}
	}

	return nil
}
