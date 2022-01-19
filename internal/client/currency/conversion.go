package currency

import (
	"context"
	"net/http"
)

type Conversion struct {
	Date string  `json:"date"`
	Eth  float64 `json:"eth"`
}

// GetUSDToETHConversion returns the conversion rate from usd to eth.
// TODO think about making this generic
func (c *currencyClient) GetUSDToETHConversion(ctx context.Context) (float64, error) {
	var data Conversion

	path := "/latest/currencies/usd/eth.json"

	// Execute request
	_, err := c.execute(ctx, http.MethodGet, path, nil, &data)

	return data.Eth, err
}
