package sources

import "github.com/ScaleneZA/CryptoTaxCalculator/cmd/transactions/transactions"

type Source interface {
	// TransformRow returns a slice of transactions because some sources contain the fees in the same row
	TransformRow(row []string) ([]transactions.Transaction, error)
}
