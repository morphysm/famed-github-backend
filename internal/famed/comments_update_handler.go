package famed

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/phuslu/log"

	"github.com/labstack/echo/v4"

	famedModel "github.com/morphysm/famed-github-backend/internal/famed/model"
	"github.com/morphysm/famed-github-backend/internal/famed/model/comment"
	"github.com/morphysm/famed-github-backend/internal/repositories/github/model"
)

type action string

const (
	updateAction action = "update"
	orderAction  action = "order"
	deleteAction action = "delete"
)

type SafeIssueCommentsUpdates struct {
	m map[int]IssueCommentUpdate
	sync.RWMutex
}

type IssueCommentUpdate struct {
	EligibleComment CommentUpdate `json:"eligibleComment"`
	RewardComment   CommentUpdate `json:"rewardComment"`
}

func NewIssueCommentUpdate() IssueCommentUpdate {
	return IssueCommentUpdate{
		EligibleComment: NewCommentUpdate(),
		RewardComment:   NewCommentUpdate(),
	}
}

type CommentUpdate struct {
	Actions []action `json:"actions"`
	Errors  []string `json:"errors"`
}

// NewCommentUpdate returns a new commentUpdate with initialized slice fields.
func NewCommentUpdate() CommentUpdate {
	return CommentUpdate{
		Actions: []action{},
		Errors:  []string{},
	}
}

func (c *CommentUpdate) AddAction(action action) {
	c.Actions = append(c.Actions, action)
}

func (c *CommentUpdate) AddError(err error) {
	c.Errors = append(c.Errors, err.Error())
}

func NewSafeIssueCommentsUpdates() *SafeIssueCommentsUpdates {
	return &SafeIssueCommentsUpdates{
		m: make(map[int]IssueCommentUpdate),
	}
}

func (sICU *SafeIssueCommentsUpdates) AddAction(issueNumber int, action action, commentType comment.Type) {
	sICU.Lock()
	defer sICU.Unlock()

	update, ok := sICU.m[issueNumber]
	if !ok {
		update = NewIssueCommentUpdate()
	}

	if commentType == comment.EligibleCommentType {
		update.EligibleComment.AddAction(action)
	}
	if commentType == comment.RewardCommentType {
		update.RewardComment.AddAction(action)
	}

	sICU.m[issueNumber] = update
}

func (sICU *SafeIssueCommentsUpdates) AddError(issueNumber int, err error, commentType comment.Type) {
	sICU.Lock()
	defer sICU.Unlock()

	update, ok := sICU.m[issueNumber]
	if !ok {
		update = NewIssueCommentUpdate()
	}

	if commentType != comment.EligibleCommentType {
		update.EligibleComment.AddError(err)
	}

	if commentType != comment.RewardCommentType {
		update.RewardComment.AddError(err)
	}

	sICU.m[issueNumber] = update
}

type updateCommentsResponse struct {
	Updates map[int]IssueCommentUpdate `json:"updates"`
}

// GetUpdateComments updates the comments in a GitHub repo.
// TODO improve efficiency. Multiple get comments calls could be avoided. Updates could be reduced to necessary.
func (gH *githubHandler) GetUpdateComments(ctx echo.Context) error {
	owner := ctx.Param("owner")
	if owner == "" {
		return echo.NewHTTPError(http.StatusBadRequest, famedModel.ErrMissingOwnerPathParameter.Error())
	}

	repoName := ctx.Param("repo_name")
	if repoName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, famedModel.ErrMissingRepoPathParameter.Error())
	}

	if ok := gH.githubInstallationClient.CheckInstallation(owner); !ok {
		return echo.NewHTTPError(http.StatusBadRequest, famedModel.ErrAppNotInstalled.Error())
	}

	issues, err := gH.githubInstallationClient.GetEnrichedIssues(ctx.Request().Context(), owner, repoName, model.Opened)
	if err != nil {
		return fmt.Errorf("failed to get issues for repository: %w", err)
	}

	commentsIssues := make(map[*model.EnrichedIssue][]model.IssueComment, len(issues))
	for _, issue := range issues {
		comments, err := gH.githubInstallationClient.GetComments(ctx.Request().Context(), owner, repoName, issue.Number)
		if err != nil {
			return err
		}
		commentsIssues[&issue] = comments
	}

	updates := NewSafeIssueCommentsUpdates()
	gH.deleteDuplicateComments(ctx.Request().Context(), owner, repoName, commentsIssues, updates)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		gH.updateRewardComments(ctx.Request().Context(), owner, repoName, commentsIssues, updates)
	}()
	go func() {
		defer wg.Done()
		gH.updateEligibleComments(ctx.Request().Context(), owner, repoName, commentsIssues, updates)
	}()

	wg.Wait()
	gH.orderComments(ctx.Request().Context(), owner, repoName, commentsIssues, updates)

	return ctx.JSON(http.StatusOK, updateCommentsResponse{Updates: updates.m})
}

// updateRewardComments checks all comments and updates comments where necessary in a concurrent fashion.
func (gH *githubHandler) updateRewardComments(ctx context.Context, owner, repoName string, commentsIssues map[*model.EnrichedIssue][]model.IssueComment, updates *SafeIssueCommentsUpdates) {
	var wg sync.WaitGroup
	defer wg.Wait()
	i := 0
	for issue, comments := range commentsIssues {
		wg.Add(1)
		go func(ctx context.Context, wg *sync.WaitGroup, owner, repoName string, issue model.EnrichedIssue) {
			defer wg.Done()
			update, err := gH.updateRewardComment(ctx, owner, repoName, issue, comments)
			if updates != nil && err != nil {
				updates.AddError(issue.Number, err, comment.RewardCommentType)
			}
			if updates != nil && update {
				updates.AddAction(issue.Number, updateAction, comment.RewardCommentType)
			}
		}(ctx, &wg, owner, repoName, *issue)
		i++
	}
}

// updateRewardComment should be run as  a go routine to check a handleClosedEvent and update the handleClosedEvent if necessary.
func (gH *githubHandler) updateRewardComment(ctx context.Context, owner, repoName string, issue model.EnrichedIssue, comments []model.IssueComment) (bool, error) {
	boardOptions := famedModel.NewBoardOptions(gH.famedConfig.Currency, famedModel.NewRewardStructure(gH.famedConfig.Rewards, gH.famedConfig.DaysToFix, 2), gH.now())
	contributors, err := famedModel.NewBlueTeamFromIssue(issue, boardOptions)
	var newComment comment.Comment
	if err != nil {
		newComment = comment.NewErrorRewardComment(err)
	}
	if len(contributors) == 0 {
		newComment = comment.NewErrorRewardComment(comment.ErrNoContributors)
	}
	if err == nil && len(contributors) > 0 {
		newComment = comment.NewRewardComment(contributors, gH.famedConfig.Currency, owner, repoName)
	}

	updated, err := gH.postOrUpdateComment(ctx, owner, repoName, issue.Number, newComment, comments)
	if err != nil {
		log.Error().Err(err).Msg("[updateRewardComment] error while posting reward comment")
		return false, err
	}

	return updated, nil
}

func (gH *githubHandler) updateEligibleComments(ctx context.Context, owner, repoName string, commentsIssues map[*model.EnrichedIssue][]model.IssueComment, updates *SafeIssueCommentsUpdates) {
	var wg sync.WaitGroup
	defer wg.Wait()
	for issue, comments := range commentsIssues {
		wg.Add(1)
		go func(issue model.Issue) {
			defer wg.Done()
			update, err := gH.updateEligibleComment(ctx, owner, repoName, issue, comments)
			if updates != nil && err != nil {
				updates.AddError(issue.Number, err, comment.EligibleCommentType)
			}
			if updates != nil && update {
				updates.AddAction(issue.Number, updateAction, comment.EligibleCommentType)
			}
		}(issue.Issue)
	}

	return
}

func (gH *githubHandler) updateEligibleComment(ctx context.Context, owner, repoName string, issue model.Issue, comments []model.IssueComment) (bool, error) {
	pullRequest, err := gH.githubInstallationClient.GetIssuePullRequest(ctx, owner, repoName, issue.Number)
	if err != nil {
		log.Error().Err(err).Msg("[updateEligibleComment] error while fetching pull request")
		return false, err
	}

	eligibleComment := comment.NewEligibleComment(issue, pullRequest)
	updated, err := gH.postOrUpdateComment(ctx, owner, repoName, issue.Number, eligibleComment, comments)
	if err != nil {
		log.Error().Err(err).Msg("[updateEligibleComment] error while posting eligible comment")
		return false, err
	}

	return updated, nil
}

// updateRewardComments checks all comments and updates comments where necessary in a concurrent fashion.
func (gH *githubHandler) deleteDuplicateComments(ctx context.Context, owner, repoName string, commentsIssues map[*model.EnrichedIssue][]model.IssueComment, updates *SafeIssueCommentsUpdates) error {
	for issue, comments := range commentsIssues {
		eligibleCommentFound := false
		rewardCommentFound := false
		for _, com := range comments {
			if comment.VerifyComment(com, gH.famedConfig.BotLogin, comment.EligibleCommentType) {
				if !eligibleCommentFound {
					eligibleCommentFound = true
				} else {
					deleted, err := gH.deleteComment(ctx, owner, repoName, com)
					if updates != nil && err != nil {
						updates.AddError(issue.Number, err, comment.EligibleCommentType)
					}
					if updates != nil && deleted {
						updates.AddAction(issue.Number, deleteAction, comment.EligibleCommentType)
					}
				}
			}
			if comment.VerifyComment(com, gH.famedConfig.BotLogin, comment.RewardCommentType) {
				if !rewardCommentFound {
					rewardCommentFound = true
				} else {
					deleted, err := gH.deleteComment(ctx, owner, repoName, com)
					if updates != nil && err != nil {
						updates.AddError(issue.Number, err, comment.RewardCommentType)
					}
					if updates != nil && deleted {
						updates.AddAction(issue.Number, deleteAction, comment.RewardCommentType)
					}
				}
			}
		}
	}

	return nil
}

// deleteComment deletes a given comment and returns a CommentUpdate.
func (gH *githubHandler) deleteComment(ctx context.Context, owner, repoName string, comment model.IssueComment) (bool, error) {
	err := gH.githubInstallationClient.DeleteComment(ctx, owner, repoName, comment.ID)
	if err != nil {
		log.Error().Err(err).Msgf("[deleteComment] error while deleting comment with id: %d", comment.ID)
		return false, err
	}

	return true, nil
}

// updateRewardComments checks all comments and updates comments where necessary in a concurrent fashion.
func (gH *githubHandler) orderComments(ctx context.Context, owner, repoName string, commentsIssues map[*model.EnrichedIssue][]model.IssueComment, updates *SafeIssueCommentsUpdates) error {
	for issue, comments := range commentsIssues {

		rewardCommentFound := false
		var rewardComment model.IssueComment
		for _, com := range comments {
			if comment.VerifyComment(com, gH.famedConfig.BotLogin, comment.EligibleCommentType) {
				if !rewardCommentFound {
					continue
				}

				// Switch comments
				gH.githubInstallationClient.UpdateComment(ctx, owner, repoName, com.ID, rewardComment.Body)
				gH.githubInstallationClient.UpdateComment(ctx, owner, repoName, rewardComment.ID, com.Body)

				updates.AddAction(issue.Number, orderAction, comment.EligibleCommentType)
				updates.AddAction(issue.Number, orderAction, comment.RewardCommentType)
			}
			if comment.VerifyComment(com, gH.famedConfig.BotLogin, comment.RewardCommentType) {
				rewardCommentFound = true
				rewardComment = com
			}
		}
	}

	return nil
}
