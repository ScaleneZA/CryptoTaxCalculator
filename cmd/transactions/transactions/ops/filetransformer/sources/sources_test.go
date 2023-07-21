package sources_test

import (
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/transactions/transactions"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/transactions/transactions/ops/filetransformer/sources"
	"github.com/luno/jettison/jtest"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBasicSource_TransformRow(t *testing.T) {
	testCases := []struct {
		name     string
		row      []string
		expected []transactions.Transaction
	}{
		{
			name: "basic buy",
			row:  []string{"1", "ETH", "0.56", "1519812503", "100", "ZAR"},
			expected: []transactions.Transaction{
				{
					UID:          "e62c505e84c2f072e9f247eeb9cb2a23",
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
			},
		},
		{
			name: "basic sell",
			row:  []string{"1", "ETH", "-0.56", "1519812503", "100", "ZAR"},
			expected: []transactions.Transaction{
				{
					UID:          "ff71ca95c3f6f2be107ccbaa2faf03b2",
					Transformer:  transactions.TransformTypeBasic,
					Currency:     "ETH",
					DetectedType: transactions.TypeSell,
					Amount:       0.56,
					Timestamp:    1519812503,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 100,
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := sources.BasicSource{}.TransformRow(tc.row)
			jtest.RequireNil(t, err)
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestBinanceSource_TransformRow(t *testing.T) {
	testCases := []struct {
		name     string
		row      []string
		expected []transactions.Transaction
	}{
		{
			name: "binance withdraw",
			row:  []string{"51633489", "2020-09-23 05:21:12", "Spot", "Withdraw", "ETH", "-1.00000000", "Withdraw fee is included"},
			expected: []transactions.Transaction{
				{
					UID:          "bbc44d01fd191c32457d3d2bfe11d996",
					Transformer:  transactions.TransformTypeBinance,
					Currency:     "ETH",
					DetectedType: transactions.TypeSendInternal,
					Amount:       1,
					Timestamp:    1600838472,
				},
			},
		},
		{
			name: "binance deposit",
			row:  []string{"51633489", "2021-01-16 15:54:35", "Spot", "Deposit", "ETC", "1.00000000", ""},
			expected: []transactions.Transaction{
				{
					UID:          "d9cc77f2bc8d3efbf0b9fd7142ca9221",
					Transformer:  transactions.TransformTypeBinance,
					Currency:     "ETC",
					DetectedType: transactions.TypeReceiveInternal,
					Amount:       1,
					Timestamp:    1610812475,
				},
			},
		},
		{
			name: "binance trade (buy)",
			row:  []string{"51633489", "2021-01-16 17:15:23", "Spot", "Large OTC trading", "ATOM", "87.80402142", ""},
			expected: []transactions.Transaction{
				{
					UID:          "9a6794b7945f6d574bfcabe0e258a24a",
					Transformer:  transactions.TransformTypeBinance,
					Currency:     "ATOM",
					DetectedType: transactions.TypeBuy,
					Amount:       87.80402142,
					Timestamp:    1610817323,
				},
			},
		},
		{
			name: "binance trade (sell)",
			row:  []string{"51633489", "2021-01-16 17:15:23", "Spot", "Large OTC trading", "ETC", "-100.00000000", ""},
			expected: []transactions.Transaction{
				{
					UID:          "2b7b5d15161bbfaa2c501b7ebdf1ea53",
					Transformer:  transactions.TransformTypeBinance,
					Currency:     "ETC",
					DetectedType: transactions.TypeSell,
					Amount:       100,
					Timestamp:    1610817323,
				},
			},
		},
		{
			name: "binance interest earned",
			row:  []string{"51633489", "2022-04-01 02:01:36", "Spot", "POS savings interest", "ATOM", "0.00417670", ""},
			expected: []transactions.Transaction{
				{
					UID:          "16169597266dfa3691c59a4e1b9a04bc",
					Transformer:  transactions.TransformTypeBinance,
					Currency:     "ATOM",
					DetectedType: transactions.TypeInterest,
					Amount:       0.00417670,
					Timestamp:    1648778496,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := sources.BinanceSource{}.TransformRow(tc.row)
			jtest.RequireNil(t, err)
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestCoinomiSource_TransformRow(t *testing.T) {
	testCases := []struct {
		name     string
		row      []string
		expected []transactions.Transaction
	}{
		{
			name: "coinomi withdraw",
			row:  []string{"Ethereum", "Ethereum", "0x0db8b9656f2404647b02979df3667ef39e903863", "", "-0.356137334", "ETH", "0.00058689", "false", "fd21481fbb7a588e41fc331da7973721ac52eeaf6ca6bac7422f32a1d7486ebb", "Tue 29 6 2021 17:57:29", "2021-06-29T17:57Z", "https://etherscan.io/tx/1234"},
			expected: []transactions.Transaction{
				{
					UID:          "fd21481fbb7a588e41fc331da7973721ac52eeaf6ca6bac7422f32a1d7486ebb",
					Transformer:  transactions.TransformTypeCoinomi,
					Currency:     "ETH",
					DetectedType: transactions.TypeSendInternal,
					Amount:       0.356137334,
					Timestamp:    1624989420,
				},
				{
					UID:          "44de929a7ebca4380db896d63fb063f0",
					Transformer:  transactions.TransformTypeCoinomi,
					Currency:     "ETH",
					DetectedType: transactions.TypeFee,
					Amount:       0.00058689,
					Timestamp:    1624989420,
				},
			},
		},
		{
			name: "coinomi withdraw - no fee",
			row:  []string{"Ethereum", "Ethereum", "0x0db8b9656f2404647b02979df3667ef39e903863", "", "-0.356137334", "ETH", "", "false", "fd21481fbb7a588e41fc331da7973721ac52eeaf6ca6bac7422f32a1d7486ebb", "Tue 29 6 2021 17:57:29", "2021-06-29T17:57Z", "https://etherscan.io/tx/1234"},
			expected: []transactions.Transaction{
				{
					UID:          "fd21481fbb7a588e41fc331da7973721ac52eeaf6ca6bac7422f32a1d7486ebb",
					Transformer:  transactions.TransformTypeCoinomi,
					Currency:     "ETH",
					DetectedType: transactions.TypeSendInternal,
					Amount:       0.356137334,
					Timestamp:    1624989420,
				},
			},
		},
		{
			name: "coinomi deposit",
			row:  []string{"XRP", "XRP", "rLHzPsX6oXkzU2qL12kHCH8G8cnZv1rBJh", "", "449.20965", "XRP", "0.000012", "false", "689d3ff18175b81ef776f7c0580d4d81fd2b2cf39cdb227903653dcae15d2142", "Tue 09 2 2021 05:47:02", "2021-02-09T05:47Z", "https://xrpcharts.ripple.com/#/transactions/1234"},
			expected: []transactions.Transaction{
				{
					UID:          "689d3ff18175b81ef776f7c0580d4d81fd2b2cf39cdb227903653dcae15d2142",
					Transformer:  transactions.TransformTypeCoinomi,
					Currency:     "XRP",
					DetectedType: transactions.TypeReceiveInternal,
					Amount:       449.20965,
					Timestamp:    1612849620,
				},
				{
					UID:          "8aab065f6d121deb579346e4ece74635",
					Transformer:  transactions.TransformTypeCoinomi,
					Currency:     "XRP",
					DetectedType: transactions.TypeFee,
					Amount:       0.000012,
					Timestamp:    1612849620,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := sources.CoinomiSource{}.TransformRow(tc.row)
			jtest.RequireNil(t, err)
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestLunoSource_TransformRow(t *testing.T) {
	testCases := []struct {
		name     string
		row      []string
		expected []transactions.Transaction
	}{
		{
			name: "luno buy",
			row:  []string{"123456789", "1", "2016-01-21 03:31:46", "Bought BTC 0.99579 for ZAR 7,500.00", "XBT", "0.99579", "0.99579", "0.99579", "0.99579", "", "", "ZAR", "7500"},
			expected: []transactions.Transaction{
				{
					UID:          "90a0a727f8b3998da38c43dc83c68a4f",
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
			},
		},
		{
			name: "luno sell",
			row:  []string{"123456789", "118", "2021-01-08 11:59:20", "Sold BTC 0.08 for ZAR 50,814.27", "XBT", "-0.08", "-0.08", "0.32008057", "0.32008057", "", "", "ZAR", "50814.27"},
			expected: []transactions.Transaction{
				{
					UID:          "008c5eaeaa06fa34503ffe3af53426f0",
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
			},
		},
		{
			name: "luno send",
			row:  []string{"123456789", "135", "2021-08-26 07:22:02", "Sent Ethereum to KNOWN234", "ETH", "-0.00072384", "0", "0.18967279", "0.18966393", "", "KNOWN234", "ZAR", "520.37"},
			expected: []transactions.Transaction{
				{
					UID:          "5ac4db3377fa6767643af8e42c04c161",
					Transformer:  transactions.TransformTypeLuno,
					Currency:     "ETH",
					DetectedType: transactions.TypeSendInternal,
					Amount:       0.00072384,
					Timestamp:    1629962522,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 718901.9672855879,
					},
				},
			},
		},
		{
			name: "luno fee",
			row:  []string{"123456789", "10", "2017-01-18 04:35:02", "Trading fee", "XBT", "-0.002442", "0", "3.803561", "3.803561", "", "", "ZAR", "31.49"},
			expected: []transactions.Transaction{
				{
					UID:          "1c3aa98c086a24f397128a240303e3bd",
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
			},
		},
		{
			name: "luno receive",
			row:  []string{"123456789", "123", "2021-04-29 13:03:53", "Received Bitcoin", "XBT", "0.09093", "0.09093", "0.76043057", "0.76043057", "", "", "ZAR", "73098.99"},
			expected: []transactions.Transaction{
				{
					UID:          "934d9b1fe0f50f6abaf01d55d0088e7a",
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
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := sources.LunoSource{}.TransformRow(tc.row)
			jtest.RequireNil(t, err)
			require.Equal(t, tc.expected, actual)
		})
	}
}
