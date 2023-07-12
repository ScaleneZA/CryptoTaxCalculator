package writer

import (
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/conversionrate/sharedtypes"
)

type Writer interface {
	Write([]sharedtypes.MarketSlice) error
}
