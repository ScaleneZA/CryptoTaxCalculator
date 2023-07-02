package filetransformer_test

import (
	"os"
	"testing"

	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/filetransformer"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/sharedtypes"
	"github.com/stretchr/testify/assert"
)

func TestTransform(t *testing.T) {
	testCases := []struct {
		name     string
		seedFile string
		typ      filetransformer.TransformType
		expected []sharedtypes.Transaction
	}{
		{
			name:     "Basic Example",
			typ:      filetransformer.TransformTypeBasic,
			seedFile: "./testData/basic.csv",
			expected: []sharedtypes.Transaction{
				{
					Currency:          "ETH",
					DetectedType:      sharedtypes.TypeBuy,
					Amount:            0.56,
					Timestamp:         1519812503,
					WholePriceAtPoint: 100,
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
					DetectedType:      sharedtypes.TypeSell,
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
			},
		},
		{
			name:     "Basic Example - Unordered, no header happy path",
			typ:      filetransformer.TransformTypeBasic,
			seedFile: "./testData/basic_unordered_no_header.csv",
			expected: []sharedtypes.Transaction{
				{
					Currency:          "ETH",
					DetectedType:      sharedtypes.TypeBuy,
					Amount:            0.56,
					Timestamp:         1519812503,
					WholePriceAtPoint: 100,
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
					DetectedType:      sharedtypes.TypeSell,
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
			},
		},
		{
			name:     "Basic Example - Multi-currency",
			typ:      filetransformer.TransformTypeBasic,
			seedFile: "./testData/basic_multi_currency.csv",
			expected: []sharedtypes.Transaction{
				{
					Currency:          "ETH",
					DetectedType:      sharedtypes.TypeBuy,
					Amount:            0.56,
					Timestamp:         1519812503,
					WholePriceAtPoint: 100,
				},
				{
					Currency:          "ETH",
					DetectedType:      sharedtypes.TypeBuy,
					Amount:            1.2,
					Timestamp:         1535450903,
					WholePriceAtPoint: 200,
				},
				{
					Currency:          "BTC",
					DetectedType:      sharedtypes.TypeBuy,
					Amount:            1,
					Timestamp:         1535450913,
					WholePriceAtPoint: 1000,
				},
				{
					Currency:          "ETH",
					DetectedType:      sharedtypes.TypeSell,
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
					Amount:            0.5,
					Timestamp:         1687946913,
					WholePriceAtPoint: 2000,
				},
			},
		},
		{
			name:     "Luno",
			typ:      filetransformer.TransformTypeLuno,
			seedFile: "./testData/LUNO_XBT.csv",
			expected: []sharedtypes.Transaction{
				{
					Currency:          "BTC",
					DetectedType:      sharedtypes.TypeBuy,
					Amount:            0.99579,
					Timestamp:         1453347106,
					WholePriceAtPoint: 7531.708492754497,
				},
				{
					Currency:          "BTC",
					DetectedType:      sharedtypes.TypeBuy,
					Amount:            1.37329,
					Timestamp:         1453624978,
					WholePriceAtPoint: 7281.783163060971,
				},
				{
					Currency:          "BTC",
					DetectedType:      sharedtypes.TypeBuy,
					Amount:            0.073002,
					Timestamp:         1466570728,
					WholePriceAtPoint: 10699.980822443222,
				},
				{
					Currency:          "BTC",
					DetectedType:      sharedtypes.TypeBuy,
					Amount:            0.38,
					Timestamp:         1466571536,
					WholePriceAtPoint: 10705,
				},
				{
					Currency:          "BTC",
					DetectedType:      sharedtypes.TypeFee,
					Amount:            0.0038,
					Timestamp:         1466571536,
					WholePriceAtPoint: 10734.21052631579,
				},
				{
					Currency:          "BTC",
					DetectedType:      sharedtypes.TypeBuy,
					Amount:            1.74346,
					Timestamp:         1480915521,
					WholePriceAtPoint: 11471.441845525564,
				},
				{
					Currency:          "BTC",
					DetectedType:      sharedtypes.TypeSendInternal,
					Amount:            1,
					Timestamp:         1484542490,
					WholePriceAtPoint: 12051,
				},
				{
					Currency:          "BTC",
					DetectedType:      sharedtypes.TypeBuy,
					Amount:            0.244261,
					Timestamp:         1484714102,
					WholePriceAtPoint: 12898.989195982984,
				},
				{
					Currency:          "BTC",
					DetectedType:      sharedtypes.TypeFee,
					Amount:            0.002442,
					Timestamp:         1484714102,
					WholePriceAtPoint: 12895.167895167893,
				},
				{
					Currency:          "BTC",
					DetectedType:      sharedtypes.TypeSendInternal,
					Amount:            0.0285,
					Timestamp:         1484840347,
					WholePriceAtPoint: 12710.877192982456,
				},
				{
					Currency:          "BTC",
					DetectedType:      sharedtypes.TypeSendInternal,
					Amount:            0.0283,
					Timestamp:         1485066626,
					WholePriceAtPoint: 12991.872791519436,
				},
				{
					Currency:          "BTC",
					DetectedType:      sharedtypes.TypeSendInternal,
					Amount:            1,
					Timestamp:         1488432689,
					WholePriceAtPoint: 16616,
				},
				{
					Currency:          "BTC",
					DetectedType:      sharedtypes.TypeSell,
					Amount:            0.08,
					Timestamp:         1610107160,
					WholePriceAtPoint: 635178.375,
				},
				{
					Currency:          "BTC",
					DetectedType:      sharedtypes.TypeReceiveInternal,
					Amount:            0.09034,
					Timestamp:         1619086743,
					WholePriceAtPoint: 800433.9163161389,
				},
				{
					Currency:          "BTC",
					DetectedType:      sharedtypes.TypeReceiveInternal,
					Amount:            0.09093,
					Timestamp:         1619701433,
					WholePriceAtPoint: 803903.9920818213,
				},
				{
					Currency:          "BTC",
					DetectedType:      sharedtypes.TypeSell,
					Amount:            0.76043057,
					Timestamp:         1619796553,
					WholePriceAtPoint: 825713.6742937623,
				},
				{
					Currency:          "BTC",
					DetectedType:      sharedtypes.TypeSell,
					Amount:            0.52295337,
					Timestamp:         1621189915,
					WholePriceAtPoint: 669275.732939631,
				},
				{
					Currency:          "BTC",
					DetectedType:      sharedtypes.TypeSendInternal,
					Amount:            0.00072384,
					Timestamp:         1629962522,
					WholePriceAtPoint: 718901.9672855879,
				},
				{
					Currency:          "BTC",
					DetectedType:      sharedtypes.TypeSendInternal,
					Amount:            0.00101206,
					Timestamp:         1652554670,
					WholePriceAtPoint: 488775.36904926586,
				},
				{
					Currency:          "BTC",
					DetectedType:      sharedtypes.TypeFee,
					Amount:            0.00002327,
					Timestamp:         1652554670,
					WholePriceAtPoint: 488611.94671250536,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			file, err := os.Open(tc.seedFile)
			assert.NoError(t, err)
			defer file.Close()

			ts, err := filetransformer.Transform(file, tc.typ)
			assert.NoError(t, err)

			assert.Equal(t, tc.expected, ts)
		})
	}

}
