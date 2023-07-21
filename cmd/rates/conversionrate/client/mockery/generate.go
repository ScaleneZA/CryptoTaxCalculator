package mockery

import "github.com/ScaleneZA/CryptoTaxCalculator/cmd/rates/conversionrate"

//go:generate mockery --dir=../.. --outpkg=mockery --output=. --case=snake --name=Client

var _ conversionrate.Client = (*Client)(nil)
