package writer

import "github.com/ScaleneZA/CryptoTaxCalculator/cmd/conversionrate/sync/sharedtypes"

type Writer interface {
	Write([]sharedtypes.MarketSlice) error
}
