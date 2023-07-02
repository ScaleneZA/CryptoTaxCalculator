package filevalidator

import (
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/taxcalculator/sharedtypes"
)

type ValidatedRow struct {
	Raw                     []string
	OverrideTransactionType sharedtypes.TransactionType
}
