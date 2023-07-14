package writer

import (
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/conversionrate/sharedtypes"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/db/markets"
)

type SQLLiteWriter struct {
	Pair sharedtypes.Pair
}

func (w SQLLiteWriter) WriteAll(b Backends, mps []sharedtypes.MarketSlice) error {
	for _, mp := range mps {
		_, err := markets.Create(b.DB(), sharedtypes.MarketPair{
			Pair:        w.Pair,
			MarketSlice: mp,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (w SQLLiteWriter) DeleteAll(b Backends) error {
	return markets.Truncate(b.DB())
}
