package sources

import (
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/taxcalculator/filetransformer/filevalidator"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/taxcalculator/sharedtypes"
)

type Source interface {
	TransformRow(vr filevalidator.ValidatedRow) (sharedtypes.Transaction, error)
}
