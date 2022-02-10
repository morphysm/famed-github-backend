package famed

import (
	"context"
	"log"
	"sync"
)

// getEvents requests all events of an issue from the GitHub API in a concurrent fashion.
func (r *repo) getEvents(ctx context.Context, repoName string) error {
	errChannel := make(chan error, len(r.issues))

	// Create context with cancel to cancel all request if one fails
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Create wait group to wait for all requests to finish
	wg := sync.WaitGroup{}

	for _, issue := range r.issues {
		wg.Add(1)

		// Start go routine to get the issue's events
		go func(ctx context.Context, repoName string, issueNumber int, issueID int64) {
			defer wg.Done()

			// Check if one of the requests returned an error otherwise run the request
			select {
			case <-ctx.Done():
				return
			default:
				eventsResp, err := r.installationClient.GetIssueEvents(ctx, repoName, issueNumber)
				if err != nil {
					log.Printf("[getEvents] error while getting events for issue with Number: %d, error: %v\n", issueNumber, err)
					cancel()
					errChannel <- err
					return
				}

				issue := r.issues[issueID]
				issue.Events = eventsResp
				r.issues[issueID] = issue
			}
		}(ctx, repoName, *issue.Issue.Number, *issue.Issue.ID)
	}

	wg.Wait()

	// Checking for error
	// TODO should we return multiple errors
	select {
	case err := <-errChannel:
		return err
	default:
	}

	return nil
}
