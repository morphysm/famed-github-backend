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
	commentEligible commentType = iota
	commentReward
)

var ErrEventNotHandled = errors.New("the event is not handled")

// handleIssuesEvent handles issue events and posts a suggested payout rewardComment to the GitHub API,
// if the famed label is set and the issue is closed.
func (gH *githubHandler) handleIssuesEvent(c echo.Context, event *github.IssuesEvent) error {
	ctx := c.Request().Context()

	// Check event base requirements
	if !isWebhookEventValid(event) {
		log.Printf("[handleIssuesEvent] err: %v", ErrEventMissingData)
		return ErrEventMissingData
	}

	var (
		commentType commentType
		comment     string
		err         error
	)
	switch *event.Action {
	case string(installation.Closed):
		comment, err = gH.rewardComment(ctx, event)
		if err != nil {
			log.Printf("[handleIssuesEvent] error while generating reward comment for closed event: %v", err)
			return err
		}
		commentType = commentReward

	case string(installation.Assigned):
		fallthrough

	case string(installation.Unassigned):
		fallthrough

	case string(installation.Unlabeled):
		fallthrough

	case string(installation.Labeled):
		comment, err = gH.eligibleComment(event)
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
	err = gH.postOrUpdateComment(ctx, *event.Repo.Owner.Login, *event.Repo.Name, *event.Issue.Number, comment, commentType)
	if err != nil {
		log.Printf("[handleIssueEvent] error while posting RewardComment: %v", err)
		return err
	}

	return c.NoContent(http.StatusOK)
}

// rewardComment returns a reward rewardComment if event and issue qualifies
func (gH *githubHandler) rewardComment(ctx context.Context, event *github.IssuesEvent) (string, error) {
	famedLabel := gH.famedConfig.Labels[config.FamedLabel]

	_, err := isCloseEventValid(event, famedLabel.Name)
	if err != nil {
		if errors.Is(err, ErrIssueMissingAssignee) {
			return RewardCommentFromError(err), nil
		}

		return "", err
	}

	repo := NewRepo(gH.famedConfig, gH.githubInstallationClient, *event.Repo.Owner.Login, *event.Repo.Name)
	return repo.RewardComment(ctx, event.Issue)
}

// rewardComment returns an eligible rewardComment if event and issue qualifies
func (gH *githubHandler) eligibleComment(event *github.IssuesEvent) (string, error) {
	famedLabel := gH.famedConfig.Labels[config.FamedLabel]
	if !isIssueFamedLabeled(event.Issue, famedLabel.Name) {
		return "", ErrEventMissingFamedLabel
	}

	if !isEligibleIssueValid(event.Issue) {
		return "", ErrIssueMissingData
	}

	return IssueEligibleComment(event.Issue)
}

// postOrUpdateComment checks if a rewardComment of a type is present,
// if so, the rewardComment body is checked for equality against the new rewardComment,
// if the comments are not equal the rewardComment is updated,
// if no rewardComment was found a new rewardComment is posted.
func (gH *githubHandler) postOrUpdateComment(ctx context.Context, owner string, repoName string, issueNumber int, comment string, commentType commentType) error {
	comments, err := gH.githubInstallationClient.GetComments(ctx, owner, repoName, issueNumber)
	if err != nil {
		return err
	}

	foundComment, found := findComment(comments, gH.famedConfig.BotLogin, commentType)
	if !found {
		return gH.githubInstallationClient.PostComment(ctx, owner, repoName, issueNumber, comment)

	}

	if *foundComment.Body != comment {
		return gH.githubInstallationClient.UpdateComment(ctx, owner, repoName, *foundComment.ID, comment)
	}

	return nil
}

// findComment finds the last of with the commentType and posted by the user with a login equal to botLogin
func findComment(comments []*github.IssueComment, botLogin string, commentType commentType) (*github.IssueComment, bool) {
	for _, comment := range comments {
		if isCommentValid(comment) &&
			isUserValid(comment.User) &&
			*comment.User.Login == botLogin &&
			verifyCommentType(*comment.Body, commentType) {
			return comment, true
		}
	}

	return nil, false
}

// verifyCommentType checks if a RewardComment is of the commentType
func verifyCommentType(str string, commentType commentType) bool {
	var substr string
	switch commentType {
	case commentEligible:
		substr = "are now eligible to Get Famed."
	case commentReward:
		substr = "Famed could not generate a reward suggestion."
		if strings.Contains(str, substr) {
			return true
		}
		substr = "Famed suggests:"
	}

	return strings.Contains(str, substr)
}
