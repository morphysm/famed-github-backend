package installation

import (
	"context"
	"log"

	"github.com/morphysm/famed-github-backend/internal/client/apps"
	"golang.org/x/oauth2"
)

type gitHubTokenSource struct {
	client         apps.Client
	installationID int64
	repoIDs        []int64
}

func NewGithubTokenSource(client apps.Client, installationID int64, repoIDs []int64) oauth2.TokenSource {
	return &gitHubTokenSource{
		client:         client,
		installationID: installationID,
		repoIDs:        repoIDs,
	}
}

// Token returns an oauth2 token.
func (tS *gitHubTokenSource) Token() (*oauth2.Token, error) {
	tokenResp, err := tS.client.GetAccessTokens(
		context.Background(),
		tS.installationID,
		tS.repoIDs)
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
