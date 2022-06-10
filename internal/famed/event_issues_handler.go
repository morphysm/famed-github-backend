package famed

import (
	"context"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/phuslu/log"

	famedModel "github.com/morphysm/famed-github-backend/internal/famed/model"
	"github.com/morphysm/famed-github-backend/internal/famed/model/comment"
	"github.com/morphysm/famed-github-backend/internal/repositories/github/model"
)

// handleIssuesEvent handles issue events and posts a suggested payout handleClosedEvent to the GitHub API,
// if the famed label is set and the issue is closed.
func (gH *githubHandler) handleIssuesEvent(c echo.Context, event model.IssuesEvent) error {
	var (
		comment comment.Comment
		err     error
		ctx     = c.Request().Context()
	)

	switch event.Action {
	case string(model.Closed):
		comment = gH.handleClosedEvent(ctx, event)
		if err != nil {
			log.Error().Err(err).Msg("[handleIssuesEvent] error while generating reward comment for closed event")
			return err
		}

	case string(model.Assigned):
		fallthrough

	case string(model.Unassigned):
		fallthrough

	case string(model.Labeled):
		fallthrough

	case string(model.Unlabeled):
		comment, err = gH.handleUpdatedEvent(ctx, event)
		if err != nil {
			log.Error().Err(err).Msg("[handleIssuesEvent] error while generating eligible comment for labeled event")
			return err
		}

	default:
		log.Error().Err(famedModel.ErrEventNotHandled).Msg("[handleIssueEvent] error")
		return famedModel.ErrEventNotHandled
	}

	comments, err := gH.githubInstallationClient.GetComments(ctx, event.Repo.Owner.Login, event.Repo.Name, event.Issue.Number)

	// Post comment to GitHub
	_, err = gH.postOrUpdateComment(ctx, event.Repo.Owner.Login, event.Repo.Name, event.Issue.Number, comment, comments)
	if err != nil {
		log.Error().Err(err).Msg("[handleIssueEvent] error while posting rewardComment")
		return err
	}

	return c.NoContent(http.StatusOK)
}

// handleClosedEvent returns a reward comment if event and issue qualifies and reopens the issue if close conditions are not met.
func (gH *githubHandler) handleClosedEvent(ctx context.Context, event model.IssuesEvent) comment.Comment {
	if len(event.Issue.Assignees) == 0 {
		return comment.NewErrorRewardComment(famedModel.ErrIssueMissingAssignee)
	}

	issue := gH.githubInstallationClient.EnrichIssue(ctx, event.Repo.Owner.Login, event.Repo.Name, event.Issue)
	// TODO: Commented out for dev connect
	//if issue.PullRequest == nil {
	//	return comment.NewErrorRewardComment(famedModel.ErrIssueMissingPullRequest)
	//}

	rewardStructure := famedModel.NewRewardStructure(gH.famedConfig.Rewards, gH.famedConfig.DaysToFix, 2)
	boardOptions := famedModel.NewBoardOptions(gH.famedConfig.Currency, rewardStructure, gH.now())
	contributors, err := famedModel.NewBlueTeamFromIssue(issue, boardOptions)
	if err != nil {
		return comment.NewErrorRewardComment(err)
	}
	if len(contributors) == 0 {
		return comment.NewErrorRewardComment(comment.ErrNoContributors)
	}

	return comment.NewRewardComment(contributors, gH.famedConfig.Currency, event.Repo.Owner.Login, event.Repo.Name)
}

// handleUpdatedEvent returns an eligible comment if event and issue qualifies
func (gH *githubHandler) handleUpdatedEvent(ctx context.Context, event model.IssuesEvent) (comment.Comment, error) {
	pullRequest, err := gH.githubInstallationClient.GetIssuePullRequest(ctx, event.Repo.Owner.Login, event.Repo.Name, event.Issue.Number)
	if err != nil {
		return nil, err
	}

	return comment.NewEligibleComment(event.Issue, pullRequest), nil
}

// postOrUpdateComment checks if a handleClosedEvent of a type is present,
// if so, the handleClosedEvent body is checked for equality against the new handleClosedEvent,
// if the comments are not equal the handleClosedEvent is updated,
// if no handleClosedEvent was found a new handleClosedEvent is posted.
func (gH *githubHandler) postOrUpdateComment(ctx context.Context, owner string, repoName string, issueNumber int, updatedComment comment.Comment, comments []model.IssueComment) (bool, error) {
	body, err := updatedComment.String()
	if err != nil {
		return false, err
	}

	foundComment, found := comment.Comments(comments).FindComment(gH.famedConfig.BotLogin, updatedComment.Type())
	if !found {
		err := gH.githubInstallationClient.PostComment(ctx, owner, repoName, issueNumber, body)
		if err != nil {
			return false, err
		}
		return true, nil
	}

	// Case-insensitive compare because the board url gets transformed to upper case.
	if !strings.EqualFold(foundComment.Body, body) {
		err := gH.githubInstallationClient.UpdateComment(ctx, owner, repoName, foundComment.ID, body)
		if err != nil {
			return false, err
		}
		return true, nil
	}

	return false, nil
}
