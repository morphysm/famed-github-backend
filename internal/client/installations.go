package client

import (
	"context"
	"net/http"
	"time"
)

type InstallationResponse []Installation

type Installation struct {
	Id      int `json:"id"`
	Account struct {
		Login             string `json:"login"`
		Id                int    `json:"id"`
		NodeId            string `json:"node_id"`
		AvatarUrl         string `json:"avatar_url"`
		GravatarId        string `json:"gravatar_id"`
		Url               string `json:"url"`
		HtmlUrl           string `json:"html_url"`
		FollowersUrl      string `json:"followers_url"`
		FollowingUrl      string `json:"following_url"`
		GistsUrl          string `json:"gists_url"`
		StarredUrl        string `json:"starred_url"`
		SubscriptionsUrl  string `json:"subscriptions_url"`
		OrganizationsUrl  string `json:"organizations_url"`
		ReposUrl          string `json:"repos_url"`
		EventsUrl         string `json:"events_url"`
		ReceivedEventsUrl string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"account"`
	RepositorySelection    string 			 `json:"repository_selection"`
	AccessTokensUrl        string            `json:"access_tokens_url"`
	RepositoriesUrl        string            `json:"repositories_url"`
	HtmlUrl                string            `json:"html_url"`
	AppId                  int               `json:"app_id"`
	TargetId               int               `json:"target_id"`
	TargetType             string            `json:"target_type"`
	Permissions            map[string]string `json:"permissions"`
	Events                 []string          `json:"events"`
	SingleFileName         string            `json:"single_file_name"`
	HasMultipleSingleFiles bool              `json:"has_multiple_single_files"`
	SingleFilePaths        []string          `json:"single_file_paths"`
	CreatedAt              time.Time         `json:"created_at"`
	UpdatedAt              time.Time         `json:"updated_at"`
	AppSlug                string            `json:"app_slug"`
	SuspendedAt            interface{}       `json:"suspended_at"`
	SuspendedBy            interface{}       `json:"suspended_by"`
}

func (c *githubClient) GetInstallations(ctx context.Context) (InstallationResponse, error) {
	var (
		resp InstallationResponse
	)
	appToken, err := c.token()
	if err != nil {
		return resp, err
	}

	_, err = c.execute(ctx, http.MethodGet, "/app/installations", appToken, nil, &resp)
	if err != nil {
		return resp, err
	}

	return resp, err
}
