package github

import (
	"context"
	"log"

	"golang.org/x/oauth2"
)

type gitHubTokenSource struct {
	client         AppClient
	installationID int64
}

func NewGithubTokenSource(client AppClient, installationID int64) oauth2.TokenSource {
	return &gitHubTokenSource{
		client:         client,
		installationID: installationID,
	}
}

// Token returns an oauth2 token.
func (tS *gitHubTokenSource) Token() (*oauth2.Token, error) {
	tokenResp, err := tS.client.GetAccessToken(
		context.Background(),
		tS.installationID)
	if err != nil {
		log.Printf("error getting access token: %v", err)
		return nil, err
	}

	token := &oauth2.Token{
		AccessToken: tokenResp.GetToken(),
		Expiry:      tokenResp.GetExpiresAt(),
	}

	return token, nil
}
