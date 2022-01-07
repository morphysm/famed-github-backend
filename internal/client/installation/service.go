package installation

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/morphysm/kudos-github-backend/internal/client/app"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate . Client
type Client interface {
	GetRepos(ctx context.Context) (ReposResponse, error)
	GetLabels(ctx context.Context, repoID string) (LabelResponse, error)
	GetEvents(ctx context.Context, repoID string) (EventResponse, error)
	GetIssues(ctx context.Context, repoID string, labels string, state IssueState) (IssueResponse, error)

	PostComment(ctx context.Context, repoName string, issueNumber int, comment string) (Comment, error)
}

type githubInstallationClient struct {
	baseURL        string
	appClient      app.Client
	installationID int
	owner          string
	client         *http.Client
	accessToken    app.AccessTokensResponse
}

// NewClient returns a new instance of the Github client
func NewClient(baseURL string, appClient app.Client, owner string, installationID int) (Client, error) {
	return &githubInstallationClient{
		baseURL:        baseURL,
		appClient:      appClient,
		installationID: installationID,
		owner:          owner,
		client:         &http.Client{},
	}, nil
}

// execute prepares and sends http requests to GitHub api.
func (c *githubInstallationClient) execute(ctx context.Context, method string, path string, token string, body []byte, object interface{}) (*http.Response, error) {
	// Set method, url and body
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Add("Accept", "application/vnd.github.v3+json")
	req.Header.Add("Authorization", "Bearer "+token)
	if method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch {
		req.Header.Add("Content-Type", "application/json;charset=UTF-8")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	// TODO extend by all valid status codes
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return nil, errors.New(fmt.Sprintf("invalid status code %d", resp.StatusCode))
	}

	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(object)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *githubInstallationClient) token(ctx context.Context) (string, error) {
	if c.accessToken.Token != "" && c.accessToken.ExpiresAt.Before(time.Now().Add(-time.Minute)) {
		return c.accessToken.Token, nil
	}

	// TODO repoID
	accessTokenResp, err := c.appClient.GetAccessTokens(ctx, c.installationID, []int{434540357, 440546811})
	if err != nil {
		return "", err
	}
	c.accessToken = accessTokenResp

	return c.accessToken.Token, nil
}
