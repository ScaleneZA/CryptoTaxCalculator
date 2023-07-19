package di

import (
	"database/sql"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/rates/db"
)

type Backends interface {
	DB() *sql.DB
}

func SetupDI() Backends {
	di := new(DI)
	di.db = db.Connect()

	return di
}

func SetupDIForTesting() Backends {
	di := new(DI)
	di.db = db.ConnectForTesting()

	return di
}

type DI struct {
	db *sql.DB
}

func (di DI) DB() *sql.DB {
	return di.db
}
