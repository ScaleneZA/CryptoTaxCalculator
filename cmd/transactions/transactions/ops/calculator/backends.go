package calculator

import (
	"database/sql"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/rates/conversionrate"
)

type Backends interface {
	RatesClient() conversionrate.Client
	DB() *sql.DB
}
