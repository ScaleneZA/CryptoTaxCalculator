package readtransformer

import (
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/rates/conversionrate"
)

type ReadTransformer interface {
	ReadAndTransform() ([]conversionrate.MarketPair, error)
}
