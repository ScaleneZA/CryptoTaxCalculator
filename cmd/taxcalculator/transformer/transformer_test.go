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
			name:     "Basic Example - Multi-currency",
			typ:      transformer.TransformTypeBasic,
			seedFile: "./testData/basic_multi_currency.csv",
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
					Currency:          "BTC",
					Typ:               sharedtypes.TypeBuy,
					Amount:            1,
					Timestamp:         1535450913,
					WholePriceAtPoint: 1000,
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
				{
					Currency:          "BTC",
					Typ:               sharedtypes.TypeSell,
					Amount:            0.5,
					Timestamp:         1687946913,
					WholePriceAtPoint: 2000,
				},
			},
		},
		{
			name:     "Luno",
			typ:      transformer.TransformTypeLuno,
			seedFile: "./testData/LUNO_XBT.csv",
			expected: []sharedtypes.Transaction{
				{
					Currency:          "BTC",
					Typ:               sharedtypes.TypeBuy,
					Amount:            0.99579,
					Timestamp:         1453347106,
					WholePriceAtPoint: 7531.708492754497,
				},
				{
					Currency:          "BTC",
					Typ:               sharedtypes.TypeBuy,
					Amount:            1.37329,
					Timestamp:         1453624978,
					WholePriceAtPoint: 7281.783163060971,
				},
				{
					Currency:          "BTC",
					Typ:               sharedtypes.TypeBuy,
					Amount:            0.073002,
					Timestamp:         1466570728,
					WholePriceAtPoint: 10699.980822443222,
				},
				{
					Currency:          "BTC",
					Typ:               sharedtypes.TypeBuy,
					Amount:            0.38,
					Timestamp:         1466571536,
					WholePriceAtPoint: 10705,
				},
				{
					Currency:          "BTC",
					Typ:               sharedtypes.TypeFee,
					Amount:            0.0038,
					Timestamp:         1466571536,
					WholePriceAtPoint: 10734.21052631579,
				},
				{
					Currency:          "BTC",
					Typ:               sharedtypes.TypeBuy,
					Amount:            1.74346,
					Timestamp:         1480915521,
					WholePriceAtPoint: 11471.441845525564,
				},
				{
					Currency:          "BTC",
					Typ:               sharedtypes.TypeSell,
					Amount:            1,
					Timestamp:         1484542490,
					WholePriceAtPoint: 12051,
				},
			},
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
