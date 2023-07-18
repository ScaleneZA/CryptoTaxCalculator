package readtransformer

import (
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/conversionrate"
)

type ReadTransformer interface {
	ReadAndTransform() ([]conversionrate.MarketPair, error)
}
