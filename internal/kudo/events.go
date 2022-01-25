package kudo

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

// getEvents gets all events for an issue in a concurrent fashion.
// TODO look through this with a fresh mind
func (bG *boardGenerator) getEvents(ctx context.Context, issues []*github.Issue, repoName string) (map[int64][]*github.IssueEvent, error) {
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
