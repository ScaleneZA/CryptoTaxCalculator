package sources

import (
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/sharedtypes"
)

type Source interface {
	TransformRow(row []string) (sharedtypes.Transaction, error)
}
