package famed

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/morphysm/famed-github-backend/internal/client/github"
	"github.com/morphysm/famed-github-backend/internal/config"
	"github.com/morphysm/famed-github-backend/pkg/pointer"
)

type safeIssueCommentsUpdates struct {
	m    map[int]issueCommentUpdate
	lock sync.RWMutex
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

func (sICU *safeIssueCommentsUpdates) Add(issueNumber int, commentUpdate commentUpdate, commentType commentType) {
	sICU.lock.Lock()
	defer sICU.lock.Unlock()

	elmt := sICU.m[issueNumber]
	if commentType != commentEligible {
		elmt.EligibleComment = commentUpdate
	}
	if commentType != commentReward {
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
		return echo.ErrBadRequest.SetInternal(ErrMissingOwnerPathParameter)
	}

	repoName := c.Param("repo_name")
	if repoName == "" {
		return echo.ErrBadRequest.SetInternal(ErrMissingRepoPathParameter)
	}

	if installed := gH.githubInstallationClient.CheckInstallation(owner); !installed {
		log.Printf("[GetUpdateComments] error on request for contributors: %v", ErrAppNotInstalled)
		return ErrAppNotInstalled
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
	wrappedIssues, err := gH.loadIssuesAndEvents(ctx, owner, repoName)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	i := 0
	for issueNumber, issue := range wrappedIssues {
		wg.Add(1)
		go func(ctx context.Context, wg *sync.WaitGroup, owner string, repoName string, issue WrappedIssue) {
			update := gH.updateRewardComment(ctx, wg, owner, repoName, issue)
			if updates != nil {
				updates.Add(issueNumber, update, commentReward)
			}
		}(ctx, &wg, owner, repoName, issue)
		i++
	}
	wg.Wait()

	return nil
}

// updateRewardComment should be run as  a go routine to check a handleClosedEvent and update the handleClosedEvent if necessary.
func (gH *githubHandler) updateRewardComment(ctx context.Context, wg *sync.WaitGroup, owner string, repoName string, issue WrappedIssue) commentUpdate {
	defer wg.Done()

	update := commentUpdate{}
	comment := ""
	contributors, err := ContributorsFromIssue(issue, BoardOptions{
		currency: gH.famedConfig.Currency,
		rewards:  gH.famedConfig.Rewards,
	})
	if err != nil {
		comment = rewardCommentFromError(err)
	}
	if err == nil {
		comment = rewardComment(contributors, gH.famedConfig.Currency, owner, repoName)
	}

	updated, err := gH.postOrUpdateComment(ctx, owner, repoName, issue.Issue.Number, comment, commentReward)
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
	famedLabel := gH.famedConfig.Labels[config.FamedLabel]
	issues, err := gH.githubInstallationClient.GetIssuesByRepo(ctx, owner, repoName, []string{famedLabel.Name}, nil)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, issue := range issues {
		wg.Add(1)
		go func(issue github.Issue) {
			update := gH.updateEligibleComment(ctx, &wg, owner, repoName, issue)
			if updates != nil {
				updates.Add(issue.Number, update, commentEligible)
			}
		}(issue)
	}

	return nil
}

func (gH *githubHandler) updateEligibleComment(ctx context.Context, wg *sync.WaitGroup, owner string, repoName string, issue github.Issue) commentUpdate {
	defer wg.Done()

	update := commentUpdate{}
	pullRequest, err := gH.githubInstallationClient.GetIssuePullRequest(ctx, owner, repoName, issue.Number)
	if err != nil {
		log.Printf("[updateEligibleComment] error while fetching pull request: %v", err)
		update.Error = err.Error()
		return update
	}

	comment := issueEligibleComment(issue, pullRequest)

	updated, err := gH.postOrUpdateComment(ctx, owner, repoName, issue.Number, comment, commentEligible)
	if err != nil {
		log.Printf("[updateEligibleComment] error while posting eligable comment: %v", err)
		update.Error = err.Error()
		return update
	}

	update.Updated = updated
	return update
}