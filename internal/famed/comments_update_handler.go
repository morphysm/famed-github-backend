package famed

import (
	"context"
	"github.com/morphysm/famed-github-backend/internal/config"
	"github.com/phuslu/log"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"

	model2 "github.com/morphysm/famed-github-backend/internal/famed/model"
	"github.com/morphysm/famed-github-backend/internal/famed/model/comment"
	"github.com/morphysm/famed-github-backend/internal/repositories/github/model"
)

type action string

const (
	updateAction action = "update"
	orderAction  action = "order"
	deleteAction action = "delete"
)

type safeIssueCommentsUpdates struct {
	m map[int]issueCommentUpdate
	sync.RWMutex
}

type issueCommentUpdate struct {
	EligibleComment commentUpdate `json:"eligibleComment"`
	RewardComment   commentUpdate `json:"rewardComment"`
}

func NewIssueCommentUpdate() issueCommentUpdate {
	return issueCommentUpdate{
		EligibleComment: NewCommentUpdate(),
		RewardComment:   NewCommentUpdate(),
	}
}

type commentUpdate struct {
	Actions []action `json:"actions"`
	Errors  []string `json:"errors"`
}

// NewCommentUpdate returns a new commentUpdate with initialized slice fields.
func NewCommentUpdate() commentUpdate {
	return commentUpdate{
		Actions: []action{},
		Errors:  []string{},
	}
}

func (c *commentUpdate) AddAction(action action) {
	c.Actions = append(c.Actions, action)
}

func (c *commentUpdate) AddError(err error) {
	c.Errors = append(c.Errors, err.Error())
}

func NewSafeIssueCommentsUpdates() *safeIssueCommentsUpdates {
	return &safeIssueCommentsUpdates{
		m: make(map[int]issueCommentUpdate),
	}
}

func (sICU *safeIssueCommentsUpdates) AddAction(issueNumber int, action action, commentType comment.Type) {
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

func (sICU *safeIssueCommentsUpdates) AddError(issueNumber int, err error, commentType comment.Type) {
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
	Updates map[int]issueCommentUpdate `json:"updates"`
}

// GetUpdateComments updates the comments in a GitHub repo.
// TODO improve efficiency. Multiple get comments calls could be avoided. Updates could be reduced to necessary.
func (gH *githubHandler) GetUpdateComments(c echo.Context) error {
	owner := c.Param("owner")
	if owner == "" {
		return echo.NewHTTPError(http.StatusBadRequest, model2.ErrMissingOwnerPathParameter.Error())
	}

	repoName := c.Param("repo_name")
	if repoName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, model2.ErrMissingRepoPathParameter.Error())
	}

	if ok := gH.githubInstallationClient.CheckInstallation(owner); !ok {
		return echo.NewHTTPError(http.StatusBadRequest, model2.ErrAppNotInstalled.Error())
	}

	famedLabel := gH.famedConfig.Labels[config.FamedLabelKey]
	issues, err := gH.githubInstallationClient.GetIssuesByRepo(c.Request().Context(), owner, repoName, []string{famedLabel.Name}, nil)
	if err != nil {
		return err
	}

	updates := NewSafeIssueCommentsUpdates()
	gH.deleteDuplicateComments(c.Request().Context(), owner, repoName, issues, updates)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		gH.updateRewardComments(c.Request().Context(), owner, repoName, issues, updates)
	}()
	go func() {
		defer wg.Done()
		gH.updateEligibleComments(c.Request().Context(), owner, repoName, issues, updates)
	}()

	wg.Wait()
	gH.orderComments(c.Request().Context(), owner, repoName, issues, updates)

	return c.JSON(http.StatusOK, updateCommentsResponse{Updates: updates.m})
}

// updateRewardComments checks all comments and updates comments where necessary in a concurrent fashion.
func (gH *githubHandler) updateRewardComments(ctx context.Context, owner string, repoName string, issues []model.Issue, updates *safeIssueCommentsUpdates) {
	enrichedIssues := gH.githubInstallationClient.EnrichIssues(ctx, owner, repoName, issues)

	var wg sync.WaitGroup
	defer wg.Wait()
	i := 0
	for _, issue := range enrichedIssues {
		wg.Add(1)
		go func(ctx context.Context, wg *sync.WaitGroup, owner string, repoName string, issue model.EnrichedIssue) {
			defer wg.Done()
			update, err := gH.updateRewardComment(ctx, owner, repoName, issue)
			if updates != nil && err != nil {
				updates.AddError(issue.Number, err, comment.RewardCommentType)
			}
			if updates != nil && update {
				updates.AddAction(issue.Number, updateAction, comment.RewardCommentType)
			}
		}(ctx, &wg, owner, repoName, issue)
		i++
	}
}

// updateRewardComment should be run as  a go routine to check a handleClosedEvent and update the handleClosedEvent if necessary.
func (gH *githubHandler) updateRewardComment(ctx context.Context, owner string, repoName string, issue model.EnrichedIssue) (bool, error) {
	boardOptions := model2.NewBoardOptions(gH.famedConfig.Currency, model2.NewRewardStructure(gH.famedConfig.Rewards, gH.famedConfig.DaysToFix, 2), gH.now())
	contributors, err := model2.NewBlueTeamFromIssue(issue, boardOptions)
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

	updated, err := gH.postOrUpdateComment(ctx, owner, repoName, issue.Number, newComment)
	if err != nil {
		log.Error().Err(err).Msg("[updateRewardComment] error while posting reward comment")
		return false, err
	}

	return updated, nil
}

func (gH *githubHandler) updateEligibleComments(ctx context.Context, owner string, repoName string, issues []model.Issue, updates *safeIssueCommentsUpdates) {
	var wg sync.WaitGroup
	defer wg.Wait()
	for _, issue := range issues {
		wg.Add(1)
		go func(issue model.Issue) {
			defer wg.Done()
			update, err := gH.updateEligibleComment(ctx, owner, repoName, issue)
			if updates != nil && err != nil {
				updates.AddError(issue.Number, err, comment.EligibleCommentType)
			}
			if updates != nil && update {
				updates.AddAction(issue.Number, updateAction, comment.EligibleCommentType)
			}
		}(issue)
	}

	return
}

func (gH *githubHandler) updateEligibleComment(ctx context.Context, owner string, repoName string, issue model.Issue) (bool, error) {
	pullRequest, err := gH.githubInstallationClient.GetIssuePullRequest(ctx, owner, repoName, issue.Number)
	if err != nil {
		log.Error().Err(err).Msg("[updateEligibleComment] error while fetching pull request")
		return false, err
	}

	comment := comment.NewEligibleComment(issue, pullRequest)
	updated, err := gH.postOrUpdateComment(ctx, owner, repoName, issue.Number, comment)
	if err != nil {
		log.Error().Err(err).Msg("[updateEligibleComment] error while posting eligable comment")
		return false, err
	}

	return updated, nil
}

// updateRewardComments checks all comments and updates comments where necessary in a concurrent fashion.
func (gH *githubHandler) deleteDuplicateComments(ctx context.Context, owner string, repoName string, issues []model.Issue, updates *safeIssueCommentsUpdates) error {
	for _, issue := range issues {
		comments, err := gH.githubInstallationClient.GetComments(ctx, owner, repoName, issue.Number)
		if err != nil {
			return err
		}

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

// deleteComment deletes a given comment and returns a commentUpdate.
func (gH *githubHandler) deleteComment(ctx context.Context, owner string, repoName string, comment model.IssueComment) (bool, error) {
	err := gH.githubInstallationClient.DeleteComment(ctx, owner, repoName, comment.ID)
	if err != nil {
		log.Error().Err(err).Msgf("[deleteComment] error while deleting comment with id: %d", comment.ID)
		return false, err
	}

	return true, nil
}

// updateRewardComments checks all comments and updates comments where necessary in a concurrent fashion.
func (gH *githubHandler) orderComments(ctx context.Context, owner string, repoName string, issues []model.Issue, updates *safeIssueCommentsUpdates) error {
	for _, issue := range issues {
		comments, err := gH.githubInstallationClient.GetComments(ctx, owner, repoName, issue.Number)
		if err != nil {
			return err
		}

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
