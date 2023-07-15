package readtransformer

import (
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/conversionrate/sharedtypes"
)

type ReadTransformer interface {
	ReadAndTransform() ([]sharedtypes.MarketPair, error)
}
