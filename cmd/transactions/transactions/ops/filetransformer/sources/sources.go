package sources

import "github.com/ScaleneZA/CryptoTaxCalculator/cmd/transactions/transactions"

type Source interface {
	TransformRow(row []string) (transactions.Transaction, error)
}
