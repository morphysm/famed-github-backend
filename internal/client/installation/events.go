package installation

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type EventsResponse []Event

type Event struct {
	ID        *int       `json:"id"`
	NodeID    *string    `json:"node_id"`
	URL       *string    `json:"url"`
	Actor     *User      `json:"actor"`
	Event     *string    `json:"event"`
	CommitID  *string    `json:"commit_id"`
	CommitURL *string    `json:"commit_url"`
	CreatedAt *time.Time `json:"created_at"`
	Issue     *Issue     `json:"issue"`
}

type EventAction string

const (
	ActionOpened     EventAction = "opened"
	ActionEdited     EventAction = "edited"
	ActionClosed     EventAction = "closed"
	ActionReopened   EventAction = "reopened"
	ActionAssigned   EventAction = "assigned"
	ActionUnassigned EventAction = "unassigned"
	ActionLabeled    EventAction = "labeled"
	ActionUnlabeled  EventAction = "unlabeled"
)

// TODO rename to GetRepoEvents
func (c *githubInstallationClient) GetEvents(ctx context.Context, repoID string) (EventsResponse, error) {
	var (
		resp EventsResponse
		path = fmt.Sprintf("/repos/%s/%s/issues/events", c.owner, repoID)
	)

	installationToken, err := c.token(ctx)
	if err != nil {
		return resp, err
	}

	_, err = c.execute(ctx, http.MethodGet, path, installationToken, nil, &resp)
	if err != nil {
		return resp, err
	}

	return resp, err
}
