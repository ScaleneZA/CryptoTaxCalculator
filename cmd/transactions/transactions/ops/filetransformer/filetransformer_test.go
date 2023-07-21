package filetransformer_test

import (
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/transactions/transactions"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/transactions/transactions/ops/filetransformer"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransform(t *testing.T) {
	testCases := []struct {
		name      string
		seedFiles []string
		typ       transactions.TransformType
		expected  []transactions.Transaction
	}{
		{
			name:      "Basic Example",
			typ:       transactions.TransformTypeBasic,
			seedFiles: []string{"./test_data/basic.csv"},
			expected: []transactions.Transaction{
				{
					Transformer:  transactions.TransformTypeBasic,
					Currency:     "ETH",
					DetectedType: transactions.TypeBuy,
					Amount:       0.56,
					Timestamp:    1519812503,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 100,
					},
				},
				{
					Transformer:  transactions.TransformTypeBasic,
					Currency:     "ETH",
					DetectedType: transactions.TypeBuy,
					Amount:       1.2,
					Timestamp:    1535450903,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 200,
					},
				},
				{
					Transformer:  transactions.TransformTypeBasic,
					Currency:     "ETH",
					DetectedType: transactions.TypeSell,
					Amount:       0.25,
					Timestamp:    1656410903,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 300,
					},
				},
				{
					Transformer:  transactions.TransformTypeBasic,
					Currency:     "ETH",
					DetectedType: transactions.TypeSell,
					Amount:       1.25,
					Timestamp:    1687946903,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 400,
					},
				},
			},
		},
		{
			name:      "Basic Example - Unordered, no header happy path",
			typ:       transactions.TransformTypeBasic,
			seedFiles: []string{"./test_data/basic_unordered_no_header.csv"},
			expected: []transactions.Transaction{
				{
					Transformer:  transactions.TransformTypeBasic,
					Currency:     "ETH",
					DetectedType: transactions.TypeBuy,
					Amount:       0.56,
					Timestamp:    1519812503,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 100,
					},
				},
				{
					Transformer:  transactions.TransformTypeBasic,
					Currency:     "ETH",
					DetectedType: transactions.TypeBuy,
					Amount:       1.2,
					Timestamp:    1535450903,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 200,
					},
				},
				{
					Transformer:  transactions.TransformTypeBasic,
					Currency:     "ETH",
					DetectedType: transactions.TypeSell,
					Amount:       0.25,
					Timestamp:    1656410903,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 300,
					},
				},
				{
					Transformer:  transactions.TransformTypeBasic,
					Currency:     "ETH",
					DetectedType: transactions.TypeSell,
					Amount:       1.25,
					Timestamp:    1687946903,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 400,
					},
				},
			},
		},

		{
			name: "Basic Example - Two files",
			typ:  transactions.TransformTypeBasic,
			seedFiles: []string{
				"./test_data/basic.csv",
				"./test_data/basic_unordered_no_header.csv",
			},
			expected: []transactions.Transaction{
				{
					Transformer:  transactions.TransformTypeBasic,
					Currency:     "ETH",
					DetectedType: transactions.TypeBuy,
					Amount:       0.56,
					Timestamp:    1519812503,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 100,
					},
				},
				{
					Transformer:  transactions.TransformTypeBasic,
					Currency:     "ETH",
					DetectedType: transactions.TypeBuy,
					Amount:       0.56,
					Timestamp:    1519812503,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 100,
					},
				},
				{
					Transformer:  transactions.TransformTypeBasic,
					Currency:     "ETH",
					DetectedType: transactions.TypeBuy,
					Amount:       1.2,
					Timestamp:    1535450903,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 200,
					},
				},
				{
					Transformer:  transactions.TransformTypeBasic,
					Currency:     "ETH",
					DetectedType: transactions.TypeBuy,
					Amount:       1.2,
					Timestamp:    1535450903,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 200,
					},
				},
				{
					Transformer:  transactions.TransformTypeBasic,
					Currency:     "ETH",
					DetectedType: transactions.TypeSell,
					Amount:       0.25,
					Timestamp:    1656410903,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 300,
					},
				},
				{
					Transformer:  transactions.TransformTypeBasic,
					Currency:     "ETH",
					DetectedType: transactions.TypeSell,
					Amount:       0.25,
					Timestamp:    1656410903,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 300,
					},
				},
				{
					Transformer:  transactions.TransformTypeBasic,
					Currency:     "ETH",
					DetectedType: transactions.TypeSell,
					Amount:       1.25,
					Timestamp:    1687946903,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 400,
					},
				},
				{
					Transformer:  transactions.TransformTypeBasic,
					Currency:     "ETH",
					DetectedType: transactions.TypeSell,
					Amount:       1.25,
					Timestamp:    1687946903,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 400,
					},
				},
			},
		},
		{
			name:      "Basic Example - Multi-currency",
			typ:       transactions.TransformTypeBasic,
			seedFiles: []string{"./test_data/basic_multi_currency.csv"},
			expected: []transactions.Transaction{
				{
					Transformer:  transactions.TransformTypeBasic,
					Currency:     "ETH",
					DetectedType: transactions.TypeBuy,
					Amount:       0.56,
					Timestamp:    1519812503,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 100,
					},
				},
				{
					Transformer:  transactions.TransformTypeBasic,
					Currency:     "ETH",
					DetectedType: transactions.TypeBuy,
					Amount:       1.2,
					Timestamp:    1535450903,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 200,
					},
				},
				{
					Transformer:  transactions.TransformTypeBasic,
					Currency:     "BTC",
					DetectedType: transactions.TypeBuy,
					Amount:       1,
					Timestamp:    1535450913,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 1000,
					},
				},
				{
					Transformer:  transactions.TransformTypeBasic,
					Currency:     "ETH",
					DetectedType: transactions.TypeSell,
					Amount:       0.25,
					Timestamp:    1656410903,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 300,
					},
				},
				{
					Transformer:  transactions.TransformTypeBasic,
					Currency:     "ETH",
					DetectedType: transactions.TypeSell,
					Amount:       1.25,
					Timestamp:    1687946903,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 400,
					},
				},
				{
					Transformer:  transactions.TransformTypeBasic,
					Currency:     "BTC",
					DetectedType: transactions.TypeSell,
					Amount:       0.5,
					Timestamp:    1687946913,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 2000,
					},
				},
			},
		},
		{
			name:      "Luno",
			typ:       transactions.TransformTypeLuno,
			seedFiles: []string{"./test_data/LUNO_XBT.csv"},
			expected: []transactions.Transaction{
				{
					Transformer:  transactions.TransformTypeLuno,
					Currency:     "BTC",
					DetectedType: transactions.TypeBuy,
					Amount:       0.99579,
					Timestamp:    1453347106,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 7531.708492754497,
					},
				},
				{
					Transformer:  transactions.TransformTypeLuno,
					Currency:     "BTC",
					DetectedType: transactions.TypeBuy,
					Amount:       1.37329,
					Timestamp:    1453624978,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 7281.783163060971,
					},
				},
				{
					Transformer:  transactions.TransformTypeLuno,
					Currency:     "BTC",
					DetectedType: transactions.TypeBuy,
					Amount:       0.073002,
					Timestamp:    1466570728,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 10699.980822443222,
					},
				},
				{
					Transformer:  transactions.TransformTypeLuno,
					Currency:     "BTC",
					DetectedType: transactions.TypeBuy,
					Amount:       0.38,
					Timestamp:    1466571536,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 10705,
					},
				},
				{
					Transformer:  transactions.TransformTypeLuno,
					Currency:     "BTC",
					DetectedType: transactions.TypeFee,
					Amount:       0.0038,
					Timestamp:    1466571536,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 10734.21052631579,
					},
				},
				{
					Transformer:  transactions.TransformTypeLuno,
					Currency:     "BTC",
					DetectedType: transactions.TypeBuy,
					Amount:       1.74346,
					Timestamp:    1480915521,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 11471.441845525564,
					},
				},
				{
					Transformer:  transactions.TransformTypeLuno,
					Currency:     "BTC",
					DetectedType: transactions.TypeSendInternal,
					Amount:       1,
					Timestamp:    1484542490,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 12051,
					},
				},
				{
					Transformer:  transactions.TransformTypeLuno,
					Currency:     "BTC",
					DetectedType: transactions.TypeBuy,
					Amount:       0.244261,
					Timestamp:    1484714102,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 12898.989195982984,
					},
				},
				{
					Transformer:  transactions.TransformTypeLuno,
					Currency:     "BTC",
					DetectedType: transactions.TypeFee,
					Amount:       0.002442,
					Timestamp:    1484714102,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 12895.167895167893,
					},
				},
				{
					Transformer:  transactions.TransformTypeLuno,
					Currency:     "BTC",
					DetectedType: transactions.TypeSendInternal,
					Amount:       0.0285,
					Timestamp:    1484840347,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 12710.877192982456,
					},
				},
				{
					Transformer:  transactions.TransformTypeLuno,
					Currency:     "BTC",
					DetectedType: transactions.TypeSendInternal,
					Amount:       0.0283,
					Timestamp:    1485066626,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 12991.872791519436,
					},
				},
				{
					Transformer:  transactions.TransformTypeLuno,
					Currency:     "BTC",
					DetectedType: transactions.TypeSendInternal,
					Amount:       1,
					Timestamp:    1488432689,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 16616,
					},
				},
				{
					Transformer:  transactions.TransformTypeLuno,
					Currency:     "BTC",
					DetectedType: transactions.TypeSell,
					Amount:       0.08,
					Timestamp:    1610107160,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 635178.375,
					},
				},
				{
					Transformer:  transactions.TransformTypeLuno,
					Currency:     "BTC",
					DetectedType: transactions.TypeReceiveInternal,
					Amount:       0.09034,
					Timestamp:    1619086743,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 800433.9163161389,
					},
				},
				{
					Transformer:  transactions.TransformTypeLuno,
					Currency:     "BTC",
					DetectedType: transactions.TypeReceiveInternal,
					Amount:       0.09093,
					Timestamp:    1619701433,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 803903.9920818213,
					},
				},
				{
					Transformer:  transactions.TransformTypeLuno,
					Currency:     "BTC",
					DetectedType: transactions.TypeSell,
					Amount:       0.76043057,
					Timestamp:    1619796553,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 825713.6742937623,
					},
				},
				{
					Transformer:  transactions.TransformTypeLuno,
					Currency:     "BTC",
					DetectedType: transactions.TypeSell,
					Amount:       0.52295337,
					Timestamp:    1621189915,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 669275.732939631,
					},
				},
				{
					Transformer:  transactions.TransformTypeLuno,
					Currency:     "BTC",
					DetectedType: transactions.TypeSendInternal,
					Amount:       0.00072384,
					Timestamp:    1629962522,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 718901.9672855879,
					},
				},
				{
					Transformer:  transactions.TransformTypeLuno,
					Currency:     "BTC",
					DetectedType: transactions.TypeSendInternal,
					Amount:       0.00101206,
					Timestamp:    1652554670,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 488775.36904926586,
					},
				},
				{
					Transformer:  transactions.TransformTypeLuno,
					Currency:     "BTC",
					DetectedType: transactions.TypeFee,
					Amount:       0.00002327,
					Timestamp:    1652554670,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 488611.94671250536,
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var files []io.Reader
			for _, sf := range tc.seedFiles {
				file, err := os.Open(sf)
				assert.NoError(t, err)
				defer file.Close()

				files = append(files, file)
			}

			ts, err := filetransformer.Transform(files, tc.typ)
			assert.NoError(t, err)

			for i, exp := range tc.expected {
				exp.UID = ts[i].UID
				assert.Equal(t, exp, ts[i])
			}
		})
	}

}

func TestTransformAll(t *testing.T) {
	testCases := []struct {
		name      string
		seedFiles map[transactions.TransformType][]string
		typ       transactions.TransformType
		expected  []transactions.Transaction
	}{
		{
			name: "Basic Example - Two different provider files",
			typ:  transactions.TransformTypeBasic,
			seedFiles: map[transactions.TransformType][]string{
				transactions.TransformTypeBasic: {
					"./test_data/basic.csv",
				},
				transactions.TransformTypeLuno: {
					"./test_data/LUNO_XBT.csv",
				},
			},
			expected: []transactions.Transaction{
				{
					Transformer:  transactions.TransformTypeLuno,
					Currency:     "BTC",
					DetectedType: transactions.TypeBuy,
					Amount:       0.99579,
					Timestamp:    1453347106,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 7531.708492754497,
					},
				},
				{
					Transformer:  transactions.TransformTypeLuno,
					Currency:     "BTC",
					DetectedType: transactions.TypeBuy,
					Amount:       1.37329,
					Timestamp:    1453624978,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 7281.783163060971,
					},
				},
				{
					Transformer:  transactions.TransformTypeLuno,
					Currency:     "BTC",
					DetectedType: transactions.TypeBuy,
					Amount:       0.073002,
					Timestamp:    1466570728,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 10699.980822443222,
					},
				},
				{
					Transformer:  transactions.TransformTypeLuno,
					Currency:     "BTC",
					DetectedType: transactions.TypeBuy,
					Amount:       0.38,
					Timestamp:    1466571536,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 10705,
					},
				},
				{
					Transformer:  transactions.TransformTypeLuno,
					Currency:     "BTC",
					DetectedType: transactions.TypeFee,
					Amount:       0.0038,
					Timestamp:    1466571536,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 10734.21052631579,
					},
				},
				{
					Transformer:  transactions.TransformTypeLuno,
					Currency:     "BTC",
					DetectedType: transactions.TypeBuy,
					Amount:       1.74346,
					Timestamp:    1480915521,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 11471.441845525564,
					},
				},
				{
					Transformer:  transactions.TransformTypeLuno,
					Currency:     "BTC",
					DetectedType: transactions.TypeSendInternal,
					Amount:       1,
					Timestamp:    1484542490,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 12051,
					},
				},
				{
					Transformer:  transactions.TransformTypeLuno,
					Currency:     "BTC",
					DetectedType: transactions.TypeBuy,
					Amount:       0.244261,
					Timestamp:    1484714102,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 12898.989195982984,
					},
				},
				{
					Transformer:  transactions.TransformTypeLuno,
					Currency:     "BTC",
					DetectedType: transactions.TypeFee,
					Amount:       0.002442,
					Timestamp:    1484714102,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 12895.167895167893,
					},
				},
				{
					Transformer:  transactions.TransformTypeLuno,
					Currency:     "BTC",
					DetectedType: transactions.TypeSendInternal,
					Amount:       0.0285,
					Timestamp:    1484840347,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 12710.877192982456,
					},
				},
				{
					Transformer:  transactions.TransformTypeLuno,
					Currency:     "BTC",
					DetectedType: transactions.TypeSendInternal,
					Amount:       0.0283,
					Timestamp:    1485066626,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 12991.872791519436,
					},
				},
				{
					Transformer:  transactions.TransformTypeLuno,
					Currency:     "BTC",
					DetectedType: transactions.TypeSendInternal,
					Amount:       1,
					Timestamp:    1488432689,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 16616,
					},
				},
				{
					Transformer:  transactions.TransformTypeBasic,
					Currency:     "ETH",
					DetectedType: transactions.TypeBuy,
					Amount:       0.56,
					Timestamp:    1519812503,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 100,
					},
				},
				{
					Transformer:  transactions.TransformTypeBasic,
					Currency:     "ETH",
					DetectedType: transactions.TypeBuy,
					Amount:       1.2,
					Timestamp:    1535450903,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 200,
					},
				},
				{
					Transformer:  transactions.TransformTypeLuno,
					Currency:     "BTC",
					DetectedType: transactions.TypeSell,
					Amount:       0.08,
					Timestamp:    1610107160,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 635178.375,
					},
				},
				{
					Transformer:  transactions.TransformTypeLuno,
					Currency:     "BTC",
					DetectedType: transactions.TypeReceiveInternal,
					Amount:       0.09034,
					Timestamp:    1619086743,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 800433.9163161389,
					},
				},
				{
					Transformer:  transactions.TransformTypeLuno,
					Currency:     "BTC",
					DetectedType: transactions.TypeReceiveInternal,
					Amount:       0.09093,
					Timestamp:    1619701433,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 803903.9920818213,
					},
				},
				{
					Transformer:  transactions.TransformTypeLuno,
					Currency:     "BTC",
					DetectedType: transactions.TypeSell,
					Amount:       0.76043057,
					Timestamp:    1619796553,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 825713.6742937623,
					},
				},
				{
					Transformer:  transactions.TransformTypeLuno,
					Currency:     "BTC",
					DetectedType: transactions.TypeSell,
					Amount:       0.52295337,
					Timestamp:    1621189915,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 669275.732939631,
					},
				},
				{
					Transformer:  transactions.TransformTypeLuno,
					Currency:     "BTC",
					DetectedType: transactions.TypeSendInternal,
					Amount:       0.00072384,
					Timestamp:    1629962522,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 718901.9672855879,
					},
				},
				{
					Transformer:  transactions.TransformTypeLuno,
					Currency:     "BTC",
					DetectedType: transactions.TypeSendInternal,
					Amount:       0.00101206,
					Timestamp:    1652554670,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 488775.36904926586,
					},
				},
				{
					Transformer:  transactions.TransformTypeLuno,
					Currency:     "BTC",
					DetectedType: transactions.TypeFee,
					Amount:       0.00002327,
					Timestamp:    1652554670,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 488611.94671250536,
					},
				},
				{
					Transformer:  transactions.TransformTypeBasic,
					Currency:     "ETH",
					DetectedType: transactions.TypeSell,
					Amount:       0.25,
					Timestamp:    1656410903,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 300,
					},
				},
				{
					Transformer:  transactions.TransformTypeBasic,
					Currency:     "ETH",
					DetectedType: transactions.TypeSell,
					Amount:       1.25,
					Timestamp:    1687946903,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 400,
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			files := make(map[transactions.TransformType][]io.Reader)
			for typ, fs := range tc.seedFiles {
				for _, sf := range fs {
					file, err := os.Open(sf)
					assert.NoError(t, err)
					defer file.Close()

					files[typ] = append(files[typ], file)
				}
			}

			ts, err := filetransformer.TransformAll(files)
			assert.NoError(t, err)

			for i, exp := range tc.expected {
				exp.UID = ts[i].UID
				assert.Equal(t, exp, ts[i])
			}
		})
	}

}
