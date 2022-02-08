package famed

import (
	"context"
	"log"
	"sync"

	"github.com/google/go-github/v41/github"
)

type eventContainer struct {
	mu     sync.Mutex
	events map[int64][]*github.IssueEvent
}

func (eC *eventContainer) add(key int64, events []*github.IssueEvent) {
	eC.mu.Lock()
	defer eC.mu.Unlock()
	eC.events[key] = events
}

// getEvents requests all events of an issue from the GitHub API in a concurrent fashion.
func (bG *boardGenerator) getEvents(ctx context.Context, issues []*github.Issue, repoName string) (map[int64][]*github.IssueEvent, error) {
	var (
		events = eventContainer{
			events: map[int64][]*github.IssueEvent{},
		}
		errChannel = make(chan error, len(issues))
	)

	// Create context with cancel to cancel all request if one fails
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Create wait group to wait for all requests to finish
	wg := sync.WaitGroup{}

	for _, issue := range issues {
		wg.Add(1)

		// Start go routine to get the issue's events
		go func(ctx context.Context, repoName string, issueNumber int, issueID int64) {
			defer wg.Done()

			// Check if one of the requests returned an error otherwise run the request
			select {
			case <-ctx.Done():
				return
			default:
				eventsResp, err := bG.installationClient.GetIssueEvents(ctx, repoName, issueNumber)
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
