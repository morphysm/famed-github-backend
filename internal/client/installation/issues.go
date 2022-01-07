package installation

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type IssueResponse []Issue

type Issue struct {
	ID                *int         `json:"id,omitempty"`
	NodeID            *string      `json:"node_id,omitempty"`
	URL               *string      `json:"url,omitempty"`
	RepositoryURL     *string      `json:"repository_url,omitempty"`
	LabelsURL         *string      `json:"labels_url,omitempty"`
	CommentsURL       *string      `json:"comments_url,omitempty"`
	EventsURL         *string      `json:"events_url,omitempty"`
	HTMLURL           *string      `json:"html_url,omitempty"`
	Number            *int         `json:"number,omitempty"`
	State             *string      `json:"state,omitempty"`
	Title             *string      `json:"title,omitempty"`
	Body              *string      `json:"body,omitempty"`
	User              *User        `json:"user,omitempty"`
	Labels            []Label      `json:"labels,omitempty"`
	Assignee          *User        `json:"assignee,omitempty"`
	Assignees         []User       `json:"assignees,omitempty"`
	Milestone         *Milestone   `json:"milestone,omitempty"`
	Locked            *bool        `json:"locked,omitempty"`
	ActiveLockReason  *string      `json:"active_lock_reason,omitempty"`
	Comments          *int         `json:"comments,omitempty"`
	PullRequest       *PullRequest `json:"pull_request,omitempty"`
	ClosedAt          *time.Time   `json:"closed_at,omitempty"`
	CreatedAt         *time.Time   `json:"created_at,omitempty"`
	UpdatedAt         *time.Time   `json:"updated_at,omitempty"`
	ClosedBy          *User        `json:"closed_by,omitempty"`
	AuthorAssociation *string      `json:"author_association,omitempty"`
}

type Milestone struct {
	URL         *string `json:"url,omitempty"`
	HTMLURL     *string `json:"html_url,omitempty"`
	LabelsURL   *string `json:"labels_url,omitempty"`
	ID          *int    `json:"id,omitempty"`
	NodeID      *string `json:"node_id,omitempty"`
	Number      *int    `json:"number,omitempty"`
	State       *string `json:"state,omitempty"`
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Creator     *User   `json:"creator,omitempty"`
}

type PullRequest struct {
	URL      *string `json:"url,omitempty"`
	HTMLURL  *string `json:"html_url,omitempty"`
	DiffURL  *string `json:"diff_url,omitempty"`
	PatchURL *string `json:"patch_url,omitempty"`
}

type User struct {
	Login             *string `json:"login,omitempty"`
	ID                *int    `json:"id,omitempty"`
	NodeID            *string `json:"node_id,omitempty"`
	AvatarURL         *string `json:"avatar_url,omitempty"`
	GravatarID        *string `json:"gravatar_id,omitempty"`
	URL               *string `json:"url,omitempty"`
	HTMLURL           *string `json:"html_url,omitempty"`
	FollowersURL      *string `json:"followers_url,omitempty"`
	FollowingURL      *string `json:"following_url,omitempty"`
	GistsURL          *string `json:"gists_url,omitempty"`
	StarredURL        *string `json:"starred_url,omitempty"`
	SubscriptionsURL  *string `json:"subscriptions_url,omitempty"`
	OrganizationsURL  *string `json:"organizations_url,omitempty"`
	ReposURL          *string `json:"repos_url,omitempty"`
	EventsURL         *string `json:"events_url,omitempty"`
	ReceivedEventsURL *string `json:"received_events_url,omitempty"`
	Type              *string `json:"type,omitempty"`
	SiteAdmin         *bool   `json:"site_admin,omitempty"`
}

type Comment struct {
	ID                *int       `json:"id"`
	NodeID            *string    `json:"node_id"`
	URL               *string    `json:"url"`
	HTMLURL           *string    `json:"html_url"`
	Body              *string    `json:"body"`
	User              *User      `json:"user"`
	CreatedAt         *time.Time `json:"created_at"`
	UpdatedAt         *time.Time `json:"updated_at"`
	IssueURL          *string    `json:"issue_url"`
	AuthorAssociation *string    `json:"author_association"`
}

type IssueState string

const (
	Open   IssueState = "open"
	Closed IssueState = "closed"
	All    IssueState = "all"
)

func (c *githubInstallationClient) GetIssues(ctx context.Context, repoName string, labels string, state IssueState) (IssueResponse, error) {
	var (
		resp IssueResponse
		// TODO verify query params
		path = fmt.Sprintf("/repos/%s/%s/issues", c.owner, repoName)
	)

	issuePathAndQuery := url.URL{Path: path}

	query := url.Values{}
	query.Add("per_page", "100")
	query.Add("labels", labels)
	query.Add("state", string(state))

	issuePathAndQuery.RawQuery = query.Encode()

	installationToken, err := c.token(ctx)
	if err != nil {
		return resp, err
	}

	_, err = c.execute(ctx, http.MethodGet, issuePathAndQuery.String(), installationToken, nil, &resp)
	if err != nil {
		return resp, err
	}

	return resp, err
}

type CommentRequest struct {
	Body string `json:"body"`
}

func (c *githubInstallationClient) PostComment(ctx context.Context, repoName string, issueNumber int, comment string) (Comment, error) {
	var (
		commentResponse Comment
		commentRequest  = CommentRequest{Body: comment}
		path            = fmt.Sprintf("/repos/%s/%s/issues/%d/comments", c.owner, repoName, issueNumber)
	)

	installationToken, err := c.token(ctx)
	if err != nil {
		return commentResponse, err
	}

	body, err := json.Marshal(commentRequest)
	if err != nil {
		return commentResponse, err
	}

	_, err = c.execute(ctx, http.MethodPost, path, installationToken, body, &commentResponse)
	if err != nil {
		return commentResponse, err
	}

	return commentResponse, err
}
