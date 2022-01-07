package app

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate . Client
type Client interface {
	GetInstallations(ctx context.Context) (InstallationResponse, error)
	GetAccessTokens(ctx context.Context, installationID int, repositoryIDs []int) (AccessTokensResponse, error)
}

type githubAppClient struct {
	baseURL string
	apiKey  jwk.Key
	appID   string
	client  *http.Client
	// TODO replace by redis cache
	accessTokens *AccessTokensResponse
}

// NewClient returns a new instance of the Github client
func NewClient(baseURL string, apiKey string, appID string) (Client, error) {
	jwkey, err := jwk.ParseKey([]byte(apiKey), jwk.WithPEM(true))
	if err != nil {
		return nil, err
	}

	return &githubAppClient{
		baseURL: baseURL,
		apiKey:  jwkey,
		appID:   appID,
		client:  &http.Client{},
	}, nil
}

// token generates a GitHub App token
func (c *githubAppClient) token() (string, error) {
	token := jwt.New()
	err := token.Set(jwt.IssuerKey, c.appID)
	if err != nil {
		return "nil", err
	}

	err = token.Set(jwt.IssuedAtKey, time.Now().Add(-time.Minute).Unix())
	if err != nil {
		return "nil", err
	}

	err = token.Set(jwt.ExpirationKey, time.Now().Add(time.Minute*5).Unix())
	if err != nil {
		return "nil", err
	}

	signedToken, err := jwt.Sign(token, jwa.RS256, c.apiKey)
	return string(signedToken), err
}

// execute prepares and sends http requests to GitHub api.
func (c *githubAppClient) execute(ctx context.Context, method string, path string, token string, body []byte, object interface{}) (*http.Response, error) {
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
		//defer resp.Body.Close()
		//buf, bodyErr := ioutil.ReadAll(resp.Body)
		//fmt.Println(buf)
		//fmt.Println(bodyErr)
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
