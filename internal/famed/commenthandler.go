package famed

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/google/go-github/v41/github"
	"github.com/labstack/echo/v4"
)

type UpdateCommentsResponse struct {
	mu            sync.Mutex
	UpdatedIssues []int `json:"updatedIssues"`
}

func (UCR *UpdateCommentsResponse) append(i int) {
	UCR.mu.Lock()
	defer UCR.mu.Unlock()
	UCR.UpdatedIssues = append(UCR.UpdatedIssues, i)
}

// UpdateComments updates the comments in a GitHub repo.
func (gH *githubHandler) UpdateComments(c echo.Context) error {
	repoName := c.Param("repo_name")
	if repoName == "" {
		return echo.ErrBadRequest.SetInternal(ErrMissingRepoPathParameter)
	}

	response := UpdateCommentsResponse{UpdatedIssues: []int{}}
	repo := NewRepo(gH.famedConfig, gH.githubInstallationClient, gH.currencyClient, repoName)

	comments, err := repo.GetComments(c.Request().Context())
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for issueNumber, comment := range comments {
		wg.Add(1)
		go gH.checkAndUpdateComment(c.Request().Context(), &wg, repoName, issueNumber, comment, &response)
	}

	wg.Wait()

	return c.JSON(http.StatusOK, response)
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

func (gH *githubHandler) checkAndUpdateComment(ctx context.Context, wg *sync.WaitGroup, repoName string, issueNumber int, comment string, response *UpdateCommentsResponse) {
	defer wg.Done()
	issueComments, err := gH.githubInstallationClient.GetComments(ctx, repoName, issueNumber)
	if err != nil {
		log.Printf("[UpdateComments] error while getting comments for issue #%d, error: %v", issueNumber, err)
		return
	}
	lastCommentByBot := getLastCommentsByUser(issueComments, 96487857)

	// TODO Add checks
	if lastCommentByBot != nil && lastCommentByBot.Body != nil && *lastCommentByBot.Body != comment {
		log.Printf("[UpdateComments] updating comment for issue #%d", issueNumber)
		_, err := gH.githubInstallationClient.PostComment(ctx, repoName, issueNumber, comment)
		if err != nil {
			log.Printf("[UpdateComments] error while posting comment for issue #%d, error: %v", issueNumber, err)
			return
		}

		response.append(issueNumber)
	}
}
