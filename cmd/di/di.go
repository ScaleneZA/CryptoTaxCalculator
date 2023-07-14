package di

import (
	"database/sql"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/db"
)

type Backends interface {
	DB() *sql.DB
}

func SetupDI() Backends {
	di := new(DI)
	di.db = db.Connect()

	return di
}

type DI struct {
	db *sql.DB
}

func (di DI) DB() *sql.DB {
	return di.db
}
