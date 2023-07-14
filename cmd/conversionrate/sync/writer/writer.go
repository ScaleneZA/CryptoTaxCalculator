package writer

import (
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/conversionrate/sharedtypes"
)

type Writer interface {
	WriteAll(Backends, []sharedtypes.MarketSlice) error
	DeleteAll(Backends) error
}
