package famed

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"

	"github.com/morphysm/famed-github-backend/internal/client/github"
)

type commentType int

const (
	commentEligible commentType = iota
	commentReward
)

var (
	ErrEventNotHandled         = errors.New("the event is not handled")
	ErrIssueMissingPullRequest = errors.New("the issue is missing a pull request")
)

// handleIssuesEvent handles issue events and posts a suggested payout handleClosedEvent to the GitHub API,
// if the famed label is set and the issue is closed.
func (gH *githubHandler) handleIssuesEvent(c echo.Context, event github.IssuesEvent) error {
	var (
		commentType commentType
		comment     string
		err         error
		ctx         = c.Request().Context()
	)

	switch event.Action {
	case string(github.Closed):
		comment, err = gH.handleClosedEvent(ctx, event)
		if err != nil {
			log.Printf("[handleIssuesEvent] error while generating reward comment for closed event: %v", err)
			return err
		}
		commentType = commentReward

	case string(github.Assigned):
		fallthrough

	case string(github.Unassigned):
		fallthrough

	case string(github.Labeled):
		fallthrough

	case string(github.Unlabeled):
		comment, err = gH.handleUpdatedEvent(ctx, event)
		if err != nil {
			log.Printf("[handleIssuesEvent] error while generating eligible comment for labeled event: %v", err)
			return err
		}
		commentType = commentEligible

	default:
		log.Printf("[handleIssueEvent] error: %v", ErrEventNotHandled)
		return ErrEventNotHandled
	}

	// Post comment to GitHub
	_, err = gH.postOrUpdateComment(ctx, event.Repo.Owner.Login, event.Repo.Name, event.Issue.Number, comment, commentType)
	if err != nil {
		log.Printf("[handleIssueEvent] error while posting rewardComment: %v", err)
		return err
	}

	return c.NoContent(http.StatusOK)
}

// handleClosedEvent returns a reward comment if event and issue qualifies and reopens the issue if close conditions are not met.
func (gH *githubHandler) handleClosedEvent(ctx context.Context, event github.IssuesEvent) (string, error) {
	if len(event.Issue.Assignees) == 0 {
		return rewardCommentFromError(ErrIssueMissingAssignee), nil
	}

	// TODO removed for DevConnect add rule
	//pullRequest, err := gH.githubInstallationClient.GetIssuePullRequest(ctx, event.Repo.Owner.Login, event.Repo.Name, event.Issue.Number)
	//if err != nil {
	//	return "", err
	//}
	//
	//if pullRequest == nil {
	//	return rewardCommentFromError(ErrIssueMissingPullRequest), nil
	//}

	issue, err := gH.loadIssueEvents(ctx, event.Repo.Owner.Login, event.Repo.Name, event.Issue)
	if err != nil {
		return "", err
	}

	contributors, err := ContributorsFromIssue(issue, BoardOptions{
		currency:  gH.famedConfig.Currency,
		rewards:   gH.famedConfig.Rewards,
		daysToFix: gH.famedConfig.DaysToFix,
	})
	if err != nil {
		return rewardCommentFromError(err), nil
	}

	return rewardComment(contributors, gH.famedConfig.Currency, event.Repo.Owner.Login, event.Repo.Name), nil
}

// handleUpdatedEvent returns an eligible comment if event and issue qualifies
func (gH *githubHandler) handleUpdatedEvent(ctx context.Context, event github.IssuesEvent) (string, error) {
	pullRequest, err := gH.githubInstallationClient.GetIssuePullRequest(ctx, event.Repo.Owner.Login, event.Repo.Name, event.Issue.Number)
	if err != nil {
		return "", err
	}

	return issueEligibleComment(event.Issue, pullRequest), nil
}

// postOrUpdateComment checks if a handleClosedEvent of a type is present,
// if so, the handleClosedEvent body is checked for equality against the new handleClosedEvent,
// if the comments are not equal the handleClosedEvent is updated,
// if no handleClosedEvent was found a new handleClosedEvent is posted.
func (gH *githubHandler) postOrUpdateComment(ctx context.Context, owner string, repoName string, issueNumber int, comment string, commentType commentType) (bool, error) {
	comments, err := gH.githubInstallationClient.GetComments(ctx, owner, repoName, issueNumber)
	if err != nil {
		return false, err
	}

	foundComment, found := findComment(comments, gH.famedConfig.BotLogin, commentType)
	if !found {
		err := gH.githubInstallationClient.PostComment(ctx, owner, repoName, issueNumber, comment)
		if err != nil {
			return false, err
		}
		return true, nil
	}

	// Case-insensitive compare because the board url gets transformed to upper case.
	if !strings.EqualFold(foundComment.Body, comment) {
		err := gH.githubInstallationClient.UpdateComment(ctx, owner, repoName, foundComment.ID, comment)
		if err != nil {
			return false, err
		}
		return true, nil
	}

	return false, nil
}
