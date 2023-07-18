package calculator_test

import (
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/taxcalculator"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/taxcalculator/ops/calculator"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCalculate(t *testing.T) {

	testCases := []struct {
		name     string
		seed     []taxcalculator.Transaction
		expected calculator.YearEndTotals
	}{
		{
			name: "Happy Path",
			seed: []taxcalculator.Transaction{
				{
					Currency:          "ETH",
					DetectedType:      taxcalculator.TypeBuy,
					Amount:            0.56,
					Timestamp:         1519812503,
					WholePriceAtPoint: 100,
				},
				{
					Currency:          "BTC",
					DetectedType:      taxcalculator.TypeBuy,
					Amount:            0.5,
					Timestamp:         1535450915,
					WholePriceAtPoint: 900,
				},
				{
					Currency:          "ETH",
					DetectedType:      taxcalculator.TypeBuy,
					Amount:            1.2,
					Timestamp:         1535450903,
					WholePriceAtPoint: 200,
				},
				{
					Currency:          "ETH",
					DetectedType:      taxcalculator.TypeBuy,
					OverridedType:     taxcalculator.TypeSell,
					Amount:            0.25,
					Timestamp:         1656410903,
					WholePriceAtPoint: 300,
				},
				{
					Currency:          "ETH",
					DetectedType:      taxcalculator.TypeSell,
					Amount:            1.25,
					Timestamp:         1687946903,
					WholePriceAtPoint: 400,
				},
				{
					Currency:          "BTC",
					DetectedType:      taxcalculator.TypeSell,
					Amount:            0.4,
					Timestamp:         1687947903,
					WholePriceAtPoint: 2000,
				},
			},
			expected: calculator.YearEndTotals{
				2023: {
					"ETH":   50,
					"TOTAL": 50,
				},
				2024: {
					"BTC":   440,
					"ETH":   281,
					"TOTAL": 721,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := calculator.Calculate("ZAR", tc.seed)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
