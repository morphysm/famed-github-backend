package famed

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
)

type IssueCommentUpdate struct {
	IssueNumber int    `json:"issueNumber"`
	Updated     bool   `json:"updated"`
	Error       string `json:"error"`
}

// GetUpdateComments updates the comments in a GitHub repo.
func (gH *githubHandler) GetUpdateComments(c echo.Context) error {
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

	response, err := gH.updateComments(c.Request().Context(), owner, repoName)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response)
}

// updateComments checks all comments and updates comments where necessary in a concurrent fashion.
func (gH *githubHandler) updateComments(ctx context.Context, owner string, repoName string) ([]*IssueCommentUpdate, error) {
	issues, err := gH.loadIssuesAndEvents(ctx, owner, repoName)
	if err != nil {
		return nil, err
	}

	issueCommentUpdates := make([]*IssueCommentUpdate, len(issues))
	if len(issues) == 0 {
		return issueCommentUpdates, nil
	}

	comments := make(map[int]string, len(issues))
	for issueNumber, issue := range issues {
		contributors, err := ContributorsFromIssue(issue, BoardOptions{
			currency: gH.famedConfig.Currency,
			rewards:  gH.famedConfig.Rewards,
		})
		if err != nil {
			comments[issueNumber] = rewardCommentFromError(err)
			continue
		}

		comments[issueNumber] = rewardComment(contributors, gH.famedConfig.Currency, owner, repoName)
	}

	var wg sync.WaitGroup
	i := 0
	for issueNumber, comment := range comments {
		issueCommentUpdates[i] = &IssueCommentUpdate{}
		wg.Add(1)
		go gH.updateComment(ctx, &wg, owner, repoName, issueNumber, comment, issueCommentUpdates[i])
		i++
	}
	wg.Wait()

	return issueCommentUpdates, nil
}

// updateComment should be run as  a go routine to check a handleClosedEvent and update the handleClosedEvent if necessary.
func (gH *githubHandler) updateComment(ctx context.Context, wg *sync.WaitGroup, owner string, repoName string, issueNumber int, comment string, issueCommentUpdate *IssueCommentUpdate) {
	defer wg.Done()
	issueCommentUpdate.IssueNumber = issueNumber

	issueComments, err := gH.githubInstallationClient.GetComments(ctx, owner, repoName, issueNumber)
	if err != nil {
		log.Printf("[CleanState] error while getting comments for issue #%d, error: %v", issueNumber, err)
		issueCommentUpdate.Error = err.Error()
		return
	}

	lastCommentByBot, found := findComment(issueComments, gH.famedConfig.BotLogin, commentReward)
	if !found {
		log.Printf("[CleanState] did not find expected comment for issue #%d", issueNumber)
		log.Printf("[CleanState] posting comment for issue #%d", issueNumber)
		err := gH.githubInstallationClient.PostComment(ctx, owner, repoName, issueNumber, comment)
		if err != nil {
			log.Printf("[CleanState] error while posting rewardComment for issue #%d, error: %v", issueNumber, err)
			issueCommentUpdate.Error = err.Error()
			return
		}
	}

	if found && lastCommentByBot.Body != comment {
		log.Printf("[CleanState] updating comment for issue #%d", issueNumber)
		err := gH.githubInstallationClient.UpdateComment(ctx, owner, repoName, lastCommentByBot.ID, comment)
		if err != nil {
			log.Printf("[CleanState] error while posting rewardComment for issue #%d, error: %v", issueNumber, err)
			issueCommentUpdate.Error = err.Error()
			return
		}
	}

	issueCommentUpdate.Updated = true
}
