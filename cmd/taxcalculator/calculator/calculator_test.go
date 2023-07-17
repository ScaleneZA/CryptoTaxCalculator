package calculator_test

import (
	calculator2 "github.com/ScaleneZA/CryptoTaxCalculator/cmd/taxcalculator/calculator"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/taxcalculator/sharedtypes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCalculate(t *testing.T) {

	testCases := []struct {
		name     string
		seed     []sharedtypes.Transaction
		expected calculator2.YearEndTotals
	}{
		{
			name: "Happy Path",
			seed: []sharedtypes.Transaction{
				{
					Currency:          "ETH",
					DetectedType:      sharedtypes.TypeBuy,
					Amount:            0.56,
					Timestamp:         1519812503,
					WholePriceAtPoint: 100,
				},
				{
					Currency:          "BTC",
					DetectedType:      sharedtypes.TypeBuy,
					Amount:            0.5,
					Timestamp:         1535450915,
					WholePriceAtPoint: 900,
				},
				{
					Currency:          "ETH",
					DetectedType:      sharedtypes.TypeBuy,
					Amount:            1.2,
					Timestamp:         1535450903,
					WholePriceAtPoint: 200,
				},
				{
					Currency:          "ETH",
					DetectedType:      sharedtypes.TypeBuy,
					OverridedType:     sharedtypes.TypeSell,
					Amount:            0.25,
					Timestamp:         1656410903,
					WholePriceAtPoint: 300,
				},
				{
					Currency:          "ETH",
					DetectedType:      sharedtypes.TypeSell,
					Amount:            1.25,
					Timestamp:         1687946903,
					WholePriceAtPoint: 400,
				},
				{
					Currency:          "BTC",
					DetectedType:      sharedtypes.TypeSell,
					Amount:            0.4,
					Timestamp:         1687947903,
					WholePriceAtPoint: 2000,
				},
			},
			expected: calculator2.YearEndTotals{
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
			actual, err := calculator2.Calculate("ZAR", tc.seed)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
