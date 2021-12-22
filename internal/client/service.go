package client

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate . Client
type Client interface {
	GetInstallations(ctx context.Context) ([]Installation, error)
}

type githubClient struct {
	baseURL   string
	apiKey    string
	appID 	  string
	client    *http.Client
}

// NewClient returns a new instance of the Github client
func NewClient(baseURL string, apiKey string, appID string) Client {
	return &githubClient{
		baseURL:	baseURL,
		apiKey:    	apiKey,
		appID: 		appID,
		client:    	&http.Client{},
	}
}

// execute is responsible of preparing, sending http requests to TOPAS api.
func (c *githubClient) execute(ctx context.Context, method string, path string, token string, body []byte, object interface{}) (*http.Response, error) {
	// Prepare request to send.
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Add(http.CanonicalHeaderKey("api-key"), c.apiKey)
	req.Header.Add(http.CanonicalHeaderKey("Accept"), "application/vnd.github.v3+json")
	req.Header.Add(http.CanonicalHeaderKey("Authorization"), "Bearer " + token)
	if method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch {
		req.Header.Add(http.CanonicalHeaderKey("Content-Type"), "application/json;charset=UTF-8")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		//defer resp.Body.Close()
		//buf, bodyErr := ioutil.ReadAll(resp.Body)
		//fmt.Println(buf)
		//fmt.Println(bodyErr)
		return nil, err
	}

	defer resp.Body.Close()
	//TODO Handle non 2xx codes
	//buf, bodyErr := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(buf))
	//fmt.Println(bodyErr)
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(object)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

