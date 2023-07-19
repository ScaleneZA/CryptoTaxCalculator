package testing

import (
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/rates/conversionrate/ops/marketvalue"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/rates/di"
)

type Client struct {
	b Backends
}

func New() (*Client, error) {
	return &Client{
		di.SetupDIForTesting(),
	}, nil
}

func (c Client) ValueAtTime(from, to string, timestamp int64) (float64, error) {
	return marketvalue.ValueAtTime(c.b, from, to, timestamp)
}
