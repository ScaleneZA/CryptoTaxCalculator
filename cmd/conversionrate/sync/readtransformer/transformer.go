package readtransformer

import "github.com/ScaleneZA/CryptoTaxCalculator/cmd/conversionrate/sync/sharedtypes"

type ReadTransformer interface {
	ReadAndTransform() ([]sharedtypes.MarketSlice, error)
}
