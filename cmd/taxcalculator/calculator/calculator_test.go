package calculator_test

import (
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/taxcalculator/calculator"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/taxcalculator/sharedtypes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCalculate(t *testing.T) {
	// TODO(Make this table-driven test)
	ts := []sharedtypes.Transaction{
		{
			Currency:          "ETH",
			Typ:               sharedtypes.TypeBuy,
			Amount:            0.56,
			Timestamp:         1519812503,
			WholePriceAtPoint: 100,
		},
		{
			Currency:          "ETH",
			Typ:               sharedtypes.TypeBuy,
			Amount:            1.2,
			Timestamp:         1535450903,
			WholePriceAtPoint: 200,
		},
		{
			Currency:          "ETH",
			Typ:               sharedtypes.TypeSell,
			Amount:            0.25,
			Timestamp:         1656410903,
			WholePriceAtPoint: 300,
		},
		{
			Currency:          "ETH",
			Typ:               sharedtypes.TypeSell,
			Amount:            1.25,
			Timestamp:         1687946903,
			WholePriceAtPoint: 400,
		},
	}

	ta := calculator.Calculate(ts)

	assert.Equal(t, map[int]map[string]float64{2023: {"ETH": 50}, 2024: {"ETH": 281}}, ta)
}
