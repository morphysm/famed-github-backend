package client

import (
	"context"
	"net/http"
	"time"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
)

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
	AccessTokensUrl string `json:"access_tokens_url"`
	RepositoriesUrl string `json:"repositories_url"`
	HtmlUrl         string `json:"html_url"`
	AppId           int    `json:"app_id"`
	TargetId        int    `json:"target_id"`
	TargetType      string `json:"target_type"`
	Permissions     struct {
		Checks   string `json:"checks"`
		Metadata string `json:"metadata"`
		Contents string `json:"contents"`
	} `json:"permissions"`
	Events                 []string    `json:"events"`
	SingleFileName         string      `json:"single_file_name"`
	HasMultipleSingleFiles bool        `json:"has_multiple_single_files"`
	SingleFilePaths        []string    `json:"single_file_paths"`
	RepositorySelection    string      `json:"repository_selection"`
	CreatedAt              time.Time   `json:"created_at"`
	UpdatedAt              time.Time   `json:"updated_at"`
	AppSlug                string      `json:"app_slug"`
	SuspendedAt            interface{} `json:"suspended_at"`
	SuspendedBy            interface{} `json:"suspended_by"`
}

func (c *githubClient) GetInstallations(ctx context.Context) ([]Installation, error) {
	var (
		resp []Installation
	)
	// TODO to function
	t := jwt.New()
	t.Set(jwt.IssuerKey, c.appID)
	t.Set(jwt.IssuedAtKey, time.Now().Add(-time.Minute).Unix())
	t.Set(jwt.ExpirationKey, time.Now().Add(time.Minute * 5).Unix())

	jwkey, err := jwk.ParseKey([]byte(c.apiKey), jwk.WithPEM(true))
	if err != nil {
		return resp, err
	}

	signedToken, err := jwt.Sign(t, jwa.RS256, jwkey)
	if err != nil {
		return resp, err
	}

	_, err = c.execute(ctx, http.MethodGet, "/app/installations", string(signedToken[:]), nil, &resp)
	if err != nil {
		return resp, err
	}

	return resp, err
}

