package webhandlers

import (
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/rates/conversionrate"
)

type Backends interface {
	RatesClient() conversionrate.Client
}
