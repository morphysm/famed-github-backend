package installation

import (
	"context"
	"fmt"
	"net/http"
)

type EventResponse []Event

type Event struct {
	Event    string `json:"event"`
	ID       int    `json:"id"`
	Assignee *User  `json:"assignee,omitempty"`
	Issue    Issue  `json:"issue"`
}

func (c *githubInstallationClient) GetEvents(ctx context.Context, repoID string) (EventResponse, error) {
	var (
		resp EventResponse
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
