package installation

import (
	"context"
	"fmt"
	"net/http"
)

type LabelResponse []Label

type Label struct {
	ID          int    `json:"id"`
	NodeID      string `json:"node_id"`
	URL         string `json:"url"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
	Default     bool   `json:"default"`
}

func (c *githubInstallationClient) GetLabels(ctx context.Context, repoID string) (LabelResponse, error) {
	var (
		resp LabelResponse
		path = fmt.Sprintf("/repos/%s/%s/labels", c.owner, repoID)
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
