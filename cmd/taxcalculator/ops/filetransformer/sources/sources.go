package sources

import (
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/taxcalculator"
)

type Source interface {
	TransformRow(row []string) (taxcalculator.Transaction, error)
}
