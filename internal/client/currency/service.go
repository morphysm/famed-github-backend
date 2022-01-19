package currency

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate . Client
type Client interface {
	GetUSDToETHConversion(ctx context.Context) (float64, error)
}

type currencyClient struct {
	baseURL string
	client  *http.Client
}

// NewCurrencyClient returns a new instance of the currency client
func NewCurrencyClient(baseURL string) Client {
	return &currencyClient{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

// execute is responsible for preparing and sending http requests to the currency api.
func (c *currencyClient) execute(ctx context.Context, method string, path string, body []byte, obj interface{}) (*http.Response, error) {
	// Prepare request to send.
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Add("Accept", "application/json;version=1")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return nil, fmt.Errorf("invalid status code %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(obj)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
