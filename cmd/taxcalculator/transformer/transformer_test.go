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
		typ      transformer.TransformType
		expected []sharedtypes.Transaction
	}{
		{
			name:     "Basic Example",
			typ:      transformer.TransformTypeBasic,
			seedFile: "./testData/basic.csv",
			expected: []sharedtypes.Transaction{
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
			},
		},
		{
			name:     "Basic Example - Unordered, no header happy path",
			typ:      transformer.TransformTypeBasic,
			seedFile: "./testData/basic_unordered_no_header.csv",
			expected: []sharedtypes.Transaction{
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
			},
		},
		{
			name:     "Luno",
			typ:      transformer.TransformTypeLuno,
			seedFile: "./testData/LUNO_XBT.csv",
			expected: []sharedtypes.Transaction{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ts, err := transformer.Transform(tc.seedFile, tc.typ)
			assert.NoError(t, err)

			assert.Equal(t, tc.expected, ts)
		})
	}

}
