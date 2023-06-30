package sources

import "github.com/ScaleneZA/CryptoTaxCalculator/cmd/taxcalculator/sharedtypes"

type Source interface {
	TransformRow(row []string) (sharedtypes.Transaction, error)
}
