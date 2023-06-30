package transformer_test

import (
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/taxcalculator/sharedtypes"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/taxcalculator/transformer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTransform(t *testing.T) {
	testCases := []struct {
		name     string
		seedFile string
		expected []sharedtypes.Transaction
	}{
		{
			name:     "Standard happy path",
			seedFile: "./testData/example.xlsx",
			expected: []sharedtypes.Transaction{
				{
					Typ:               sharedtypes.TypeBuy,
					Amount:            0.56,
					Timestamp:         1519812503,
					WholePriceAtPoint: 100,
				},
				{
					Typ:               sharedtypes.TypeBuy,
					Amount:            1.2,
					Timestamp:         1535450903,
					WholePriceAtPoint: 200,
				},
				{
					Typ:               sharedtypes.TypeSell,
					Amount:            0.25,
					Timestamp:         1656410903,
					WholePriceAtPoint: 300,
				},
				{
					Typ:               sharedtypes.TypeSell,
					Amount:            1.25,
					Timestamp:         1687946903,
					WholePriceAtPoint: 400,
				},
			},
		},
		{
			name:     "Unordered, no header happy path",
			seedFile: "./testData/example_unordered_no_header.xlsx",
			expected: []sharedtypes.Transaction{
				{
					Typ:               sharedtypes.TypeBuy,
					Amount:            0.56,
					Timestamp:         1519812503,
					WholePriceAtPoint: 100,
				},
				{
					Typ:               sharedtypes.TypeBuy,
					Amount:            1.2,
					Timestamp:         1535450903,
					WholePriceAtPoint: 200,
				},
				{
					Typ:               sharedtypes.TypeSell,
					Amount:            0.25,
					Timestamp:         1656410903,
					WholePriceAtPoint: 300,
				},
				{
					Typ:               sharedtypes.TypeSell,
					Amount:            1.25,
					Timestamp:         1687946903,
					WholePriceAtPoint: 400,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ts, err := transformer.Transform(tc.seedFile, transformer.TransformTypeTest)
			assert.NoError(t, err)

			assert.Equal(t, tc.expected, ts)
		})
	}

}
