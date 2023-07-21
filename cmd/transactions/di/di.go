package di

import (
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/rates/conversionrate"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/rates/conversionrate/client/grpc"
)

type Backends interface {
	RatesClient() conversionrate.Client
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

	return di
}

func SetupDIForTesting(b BackendsTest) Backends {
	di := new(DI)
	di.ratesClient = b.RatesClient

	return di
}

type DI struct {
	ratesClient conversionrate.Client
}

func (di DI) RatesClient() conversionrate.Client {
	return di.ratesClient
}
