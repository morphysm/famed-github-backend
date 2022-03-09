package famed

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/google/go-github/v41/github"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/morphysm/famed-github-backend/internal/client/installation"
	"github.com/morphysm/famed-github-backend/internal/config"
)

type commentType int

const (
	closed commentType = iota
	labeled
)

// handleIssuesEvent handles issue events and posts a suggested payout comment to the GitHub API,
// if the famed label is set and the issue is closed.
func (gH *githubHandler) handleIssuesEvent(c echo.Context, event *github.IssuesEvent) error {
	ctx := c.Request().Context()

	var commentType commentType
	var comment string
	var err error
	//TODO add action check
	switch *event.Action {
	case string(installation.Closed):
		comment, err = gH.closeEventToComment(ctx, event)
		if err != nil {
			log.Printf("[handleIssuesEvent] error while generating comment for closed event: %v", err)
			// TODO return c.NoContent(http.StatusOK)
			return err
		}
		commentType = closed
	case string(installation.Labeled):
		comment, err = gH.labeledEventToComment(ctx, event)
		if err != nil {
			log.Printf("[handleIssuesEvent] error while generating comment for labeled event: %v", err)
			// TODO return c.NoContent(http.StatusOK)
			return err
		}
		commentType = labeled
	default:
		log.Print("received unhandled issues event")
		return c.NoContent(http.StatusOK)
	}

	// Post comment to GitHub
	err = gH.postOrUpdateComment(ctx, *event.Repo.Owner.Login, *event.Repo.Name, *event.Issue.Number, comment, commentType)
	if err != nil {
		log.Printf("[handleIssueEvent] error while posting comment: %v", err)
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (gH *githubHandler) closeEventToComment(ctx context.Context, event *github.IssuesEvent) (string, error) {
	famedLabel := gH.famedConfig.Labels[config.FamedLabel]

	_, err := isValidCloseEvent(event, famedLabel.Name)
	if err != nil {
		if errors.Is(err, ErrIssueMissingAssignee) {
			return commentFromError(err), nil
		}
		return "", err
	}

	repo := NewRepo(gH.famedConfig, gH.githubInstallationClient, *event.Repo.Owner.Login, *event.Repo.Name)
	return repo.ContributorComment(ctx, event.Issue)
}

func (gH *githubHandler) labeledEventToComment(ctx context.Context, event *github.IssuesEvent) (string, error) {
	famedLabel := gH.famedConfig.Labels[config.FamedLabel]

	// TODO add severity labels
	// TODO add validity check
	if event.Label == nil || event.Label.Name == nil || *event.Label.Name != famedLabel.Name {
		return "", errors.New("label is not \"famed\" label")
	}

	// TODO move this
	issue, _ := gH.githubInstallationClient.GetIssue(ctx, *event.Repo.Owner.Login, *event.Repo.Name, *event.Issue.Number)
	repo := NewRepo(gH.famedConfig, gH.githubInstallationClient, *event.Repo.Owner.Login, *event.Repo.Name)
	return repo.IssueStateComment(ctx, issue)
}

func (gH *githubHandler) postOrUpdateComment(ctx context.Context, owner string, repoName string, issueNumber int, comment string, commentType commentType) error {
	comments, err := gH.githubInstallationClient.GetComments(ctx, owner, repoName, issueNumber)
	if err != nil {
		return err
	}

	commentID, found := findComment(comments, gH.famedConfig.BotLogin, commentType)
	if found {
		return gH.githubInstallationClient.UpdateComment(ctx, owner, repoName, commentID, comment)
	}

	return gH.githubInstallationClient.PostComment(ctx, owner, repoName, issueNumber, comment)
}

func findComment(comments []*github.IssueComment, botLogin string, commentType commentType) (int64, bool) {
	for _, comment := range comments {
		// TODO validate comment
		if *comment.User.Login == botLogin && checkCommentType(*comment.Body, commentType) {
			return *comment.ID, true
		}
	}

	return -1, false
}

func checkCommentType(str string, commentType commentType) bool {
	var substr string
	switch commentType {
	case labeled:
		substr = "are now eligible to Get Famed."
	case closed:
		substr = "Famed suggests:"
	}

	return strings.Contains(str, substr)
}
