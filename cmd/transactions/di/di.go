package di

import (
	"database/sql"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/rates/conversionrate"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/rates/conversionrate/client/grpc"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/transactions/db"
)

type Backends interface {
	RatesClient() conversionrate.Client
	DB() *sql.DB
}

type BackendsTest struct {
	RatesClient conversionrate.Client
}

func SetupDI() Backends {
	di := new(DI)

	ratesClient, err := grpc.New()
	if err != nil {
		panic("failed to create new grpc client")
	}

	di.ratesClient = ratesClient
	di.db = db.Connect()

	return di
}

func SetupDIForTesting(b BackendsTest) Backends {
	di := new(DI)

	di.ratesClient = b.RatesClient
	di.db = db.ConnectForTesting()

	return di
}

type DI struct {
	ratesClient conversionrate.Client
	db          *sql.DB
}

func (di DI) RatesClient() conversionrate.Client {
	return di.ratesClient
}

func (di DI) DB() *sql.DB {
	return di.db
}
