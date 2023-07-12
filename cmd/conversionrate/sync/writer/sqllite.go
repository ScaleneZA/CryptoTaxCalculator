package writer

import (
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/conversionrate/sync/sharedtypes"
)

type SQLLiteWriter struct {
	Filename string
}

func (w SQLLiteWriter) Write(mps []sharedtypes.MarketSlice) error {
	// TODO: work
	return nil
}
