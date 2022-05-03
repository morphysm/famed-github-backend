package famed

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"

	"github.com/morphysm/famed-github-backend/internal/config"
	model2 "github.com/morphysm/famed-github-backend/internal/famed/model"
	"github.com/morphysm/famed-github-backend/internal/famed/model/comment"
	"github.com/morphysm/famed-github-backend/internal/respositories/github/model"
	"github.com/morphysm/famed-github-backend/pkg/pointer"
)

type safeIssueCommentsUpdates struct {
	m map[int]issueCommentUpdate
	sync.RWMutex
}

type issueCommentUpdate struct {
	EligibleComment commentUpdate `json:"eligibleComment"`
	RewardComment   commentUpdate `json:"rewardComment"`
}

type commentUpdate struct {
	Updated bool   `json:"updated"`
	Error   string `json:"error"`
}

func NewSafeIssueCommentsUpdates() *safeIssueCommentsUpdates {
	return &safeIssueCommentsUpdates{
		m: make(map[int]issueCommentUpdate),
	}
}

func (sICU *safeIssueCommentsUpdates) Add(issueNumber int, commentUpdate commentUpdate, commentType comment.Type) {
	sICU.Lock()
	defer sICU.Unlock()

	elmt := sICU.m[issueNumber]
	if commentType != comment.EligibleCommentType {
		elmt.EligibleComment = commentUpdate
	}
	if commentType != comment.RewardCommentType {
		elmt.RewardComment = commentUpdate
	}
	sICU.m[issueNumber] = elmt
}

type updateCommentsResponse struct {
	RewardCommentsError   *string                    `json:"rewardCommentError,omitempty"`
	EligibleCommentsError *string                    `json:"eligibleCommentError,omitempty"`
	Updates               map[int]issueCommentUpdate `json:"updates"`
}

// GetUpdateComments updates the comments in a GitHub repo.
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

	var wg sync.WaitGroup
	updates := NewSafeIssueCommentsUpdates()
	response := updateCommentsResponse{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		err := gH.updateRewardComments(c.Request().Context(), owner, repoName, updates)
		if err != nil {
			response.RewardCommentsError = pointer.String(err.Error())
		}
	}()
	go func() {
		defer wg.Done()
		err := gH.updateEligibleComments(c.Request().Context(), owner, repoName, updates)
		if err != nil {
			response.EligibleCommentsError = pointer.String(err.Error())
		}
	}()

	wg.Wait()
	response.Updates = updates.m
	return c.JSON(http.StatusOK, response)
}

// updateRewardComments checks all comments and updates comments where necessary in a concurrent fashion.
func (gH *githubHandler) updateRewardComments(ctx context.Context, owner string, repoName string, updates *safeIssueCommentsUpdates) error {
	wrappedIssues, err := gH.loadIssues(ctx, owner, repoName)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	i := 0
	for issueNumber, issue := range wrappedIssues {
		wg.Add(1)
		go func(ctx context.Context, wg *sync.WaitGroup, owner string, repoName string, issue model2.EnrichedIssue) {
			update := gH.updateRewardComment(ctx, wg, owner, repoName, issue)
			if updates != nil {
				updates.Add(issueNumber, update, comment.RewardCommentType)
			}
		}(ctx, &wg, owner, repoName, issue)
		i++
	}
	wg.Wait()

	return nil
}

// updateRewardComment should be run as  a go routine to check a handleClosedEvent and update the handleClosedEvent if necessary.
func (gH *githubHandler) updateRewardComment(ctx context.Context, wg *sync.WaitGroup, owner string, repoName string, issue model2.EnrichedIssue) commentUpdate {
	defer wg.Done()

	update := commentUpdate{}
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
		log.Printf("[updateRewardComment] error while posting reward comment: %v", err)
		update.Error = err.Error()
		return update
	}

	update.Updated = updated
	return update
}

func (gH *githubHandler) updateEligibleComments(ctx context.Context, owner string, repoName string, updates *safeIssueCommentsUpdates) error {
	// TODO duplicate GetIssues call
	famedLabel := gH.famedConfig.Labels[config.FamedLabelKey]
	issues, err := gH.githubInstallationClient.GetIssuesByRepo(ctx, owner, repoName, []string{famedLabel.Name}, nil)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, issue := range issues {
		wg.Add(1)
		go func(issue model.Issue) {
			update := gH.updateEligibleComment(ctx, &wg, owner, repoName, issue)
			if updates != nil {
				updates.Add(issue.Number, update, comment.EligibleCommentType)
			}
		}(issue)
	}

	return nil
}

func (gH *githubHandler) updateEligibleComment(ctx context.Context, wg *sync.WaitGroup, owner string, repoName string, issue model.Issue) commentUpdate {
	defer wg.Done()

	update := commentUpdate{}
	pullRequest, err := gH.githubInstallationClient.GetIssuePullRequest(ctx, owner, repoName, issue.Number)
	if err != nil {
		log.Printf("[updateEligibleComment] error while fetching pull request: %v", err)
		update.Error = err.Error()
		return update
	}

	comment := comment.NewEligibleComment(issue, pullRequest)

	updated, err := gH.postOrUpdateComment(ctx, owner, repoName, issue.Number, comment)
	if err != nil {
		log.Printf("[updateEligibleComment] error while posting eligable comment: %v", err)
		update.Error = err.Error()
		return update
	}

	update.Updated = updated
	return update
}
