package di

import (
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/rates/conversionrate"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/rates/conversionrate/client/grpc"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/rates/conversionrate/client/testing"
)

type Backends interface {
	RatesClient() conversionrate.Client
}

func SetupDI() Backends {
	di := new(DI)

	ratesClient, err := grpc.New()
	if err != nil {
		panic("failed to create new grpc client")
	}

	di.ratesClient = ratesClient

	return di
}

func SetupDIForTesting() Backends {
	di := new(DI)

	ratesClient, err := testing.New()
	if err != nil {
		panic("failed to create new testing client")
	}

	di.ratesClient = ratesClient

	return di
}

type DI struct {
	ratesClient conversionrate.Client
}

func (di DI) RatesClient() conversionrate.Client {
	return di.ratesClient
}
