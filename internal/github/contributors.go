package github

import (
	"context"
	"errors"
	"log"
	"net/http"
	"sort"
	"sync"

	"github.com/google/go-github/v41/github"
	"github.com/labstack/echo/v4"

	"github.com/morphysm/kudos-github-backend/internal/client/installation"
	"github.com/morphysm/kudos-github-backend/internal/kudo"
)

func (gH *githubHandler) GetContributors(c echo.Context) error {
	repoName := c.Param("repo_name")
	if repoName == "" {
		return echo.ErrBadRequest.SetInternal(errors.New("missing repo name path parameter"))
	}

	// Get all issues in repo
	issuesResponse, err := gH.githubInstallationClient.GetIssuesByRepo(c.Request().Context(), repoName, []string{gH.kudoLabel}, installation.Closed)
	if err != nil {
		return echo.ErrBadGateway.SetInternal(err)
	}

	// Use issues to generate contributor list
	contributors, err := gH.issuesToContributors(c.Request().Context(), issuesResponse, repoName)
	if err != nil {
		return echo.ErrBadGateway.SetInternal(err)
	}

	return c.JSON(http.StatusOK, contributors)
}

type eventContainer struct {
	mu     sync.Mutex
	events map[int64][]*github.IssueEvent
}

func (eC *eventContainer) add(key int64, events []*github.IssueEvent) {
	eC.mu.Lock()
	defer eC.mu.Unlock()
	eC.events[key] = events
}

// issuesToContributors generates a contributor list based on a list of issues
func (gH *githubHandler) issuesToContributors(ctx context.Context, issues []*github.Issue, repoName string) ([]*kudo.Contributor, error) {
	var (
		contributorsArray = []*kudo.Contributor{}
		filteredIssues    []*github.Issue
	)

	if len(issues) == 0 {
		return contributorsArray, nil
	}

	usdToEthRate, err := gH.currencyClient.GetUSDToETHConversion(ctx)
	if err != nil {
		return nil, err
	}

	for _, issue := range issues {
		if !kudo.IsIssueValid(issue) {
			log.Printf("[issuesToContributors] issue invalid with ID: %d \n", issue.ID)
			continue
		}
		filteredIssues = append(filteredIssues, issue)
	}

	events, err := gH.getEvents(ctx, filteredIssues, repoName)
	if err != nil {
		return nil, err
	}

	contributors := kudo.GenerateContributors(filteredIssues, events, gH.kudoRewardUnit, gH.kudoRewards, usdToEthRate)

	// Transformation of contributors map to contributors array
	for _, contributor := range contributors {
		contributorsArray = append(contributorsArray, contributor)
	}

	// Sort contributors array by total rewards
	sort.SliceStable(contributorsArray, func(i, j int) bool {
		return contributorsArray[i].RewardSum > contributorsArray[j].RewardSum
	})

	return contributorsArray, nil
}

// getEvents gets all events for an issue in a concurrent fashion.
// TODO look through this with a fresh mind
func (gH *githubHandler) getEvents(ctx context.Context, issues []*github.Issue, repoName string) (map[int64][]*github.IssueEvent, error) {
	var (
		events = eventContainer{
			events: map[int64][]*github.IssueEvent{},
		}
		errChannel = make(chan error, len(issues))
	)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	wg := sync.WaitGroup{}

	for _, issue := range issues {
		wg.Add(1)
		issue := issue
		go func(ctx context.Context, repoName string, issueNumber int, issueID int64) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				return
			default:
				eventsResp, err := gH.githubInstallationClient.GetIssueEvents(ctx, repoName, issueNumber)
				if err != nil {
					log.Printf("[getEvents] error while getting events for issue with Number: %d, error: %v\n", issueNumber, err)
					cancel()
					errChannel <- err
					return
				}

				events.add(issueID, eventsResp)
			}
		}(ctx, repoName, *issue.Number, *issue.ID)
	}

	wg.Wait()

	// Checking for error
	// TODO should we return multiple errors
	select {
	case err := <-errChannel:
		return nil, err
	default:
	}

	return events.events, nil
}
