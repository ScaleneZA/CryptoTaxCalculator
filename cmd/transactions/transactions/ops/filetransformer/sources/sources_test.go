package sources_test

import (
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/transactions/transactions"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/transactions/transactions/ops/filetransformer/sources"
	"github.com/luno/jettison/jtest"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
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
					Timestamp:    time.Unix(1519812503, 0),
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
					Timestamp:    time.Unix(1519812503, 0),
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
					Timestamp:    time.Unix(1600838472, 0),
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
					Timestamp:    time.Unix(1610812475, 0),
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
					Timestamp:    time.Unix(1610817323, 0),
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
					Timestamp:    time.Unix(1610817323, 0),
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
					Timestamp:    time.Unix(1648778496, 0),
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

func TestKrakenSource_TransformRow(t *testing.T) {
	testCases := []struct {
		name     string
		row      []string
		expected []transactions.Transaction
	}{
		{
			name: "kraken withdraw",
			row:  []string{"LUFQ64-G33PS-2AGKX4", "ASBHH3K-SSTHXU-T7T3JD", "2017-07-12 12:55:51", "withdrawal", "", "currency", "XXDG", "-499998.00000000", "2.00000000", "0.00000000"},
			expected: []transactions.Transaction{
				{
					UID:          "fba5670c958c5429948dc1954ac2f640",
					Transformer:  transactions.TransformTypeKraken,
					Currency:     "DOGE",
					DetectedType: transactions.TypeSendInternal,
					Amount:       499998.00000000,
					Timestamp:    time.Unix(1499864151, 0),
				},
				{
					UID:          "c18e5beb15875f9075867e87962a17e9",
					Transformer:  transactions.TransformTypeKraken,
					Currency:     "DOGE",
					DetectedType: transactions.TypeFee,
					Amount:       2,
					Timestamp:    time.Unix(1499864151, 0),
				},
			},
		},
		{
			name: "kraken deposit",
			row:  []string{"", "QGBBX47-IMV4HY-D4L774", "2017-01-16 04:55:07", "deposit", "", "currency", "XXBT", "1.0000000000", "0.0000000000", ""},
			expected: []transactions.Transaction{
				{
					UID:          "4324ce5d00ca0890eb658a7a8b9bb516",
					Transformer:  transactions.TransformTypeKraken,
					Currency:     "BTC",
					DetectedType: transactions.TypeReceiveInternal,
					Amount:       1,
					Timestamp:    time.Unix(1484542507, 0),
				},
			},
		},
		{
			name: "kraken trade (buy)",
			row:  []string{"LDGXVI-CFEUU-RDF5EO", "TRHBHN-B7MT4-5GFBCK", "2017-01-16 05:08:57", "trade", "", "currency", "XETH", "69.4815000000", "0.0000000000", "69.4815000000"},
			expected: []transactions.Transaction{
				{
					UID:          "76a51894e501c19cac22522653568991",
					Transformer:  transactions.TransformTypeKraken,
					Currency:     "ETH",
					DetectedType: transactions.TypeBuy,
					Amount:       69.4815000000,
					Timestamp:    time.Unix(1484543337, 0),
				},
			},
		},
		{
			name: "kraken trade (sell)",
			row:  []string{"LQA2A6-MGIDN-2G3OJU", "TRHBHN-B7MT4-5GFBCK", "2017-01-16 05:08:57", "trade", "", "currency", "XXBT", "-0.8143930000", "0.0013030000", "0.1843040000"},
			expected: []transactions.Transaction{
				{
					UID:          "e1c7fbac21d5917bf1177d28cfcbefda",
					Transformer:  transactions.TransformTypeKraken,
					Currency:     "BTC",
					DetectedType: transactions.TypeSell,
					Amount:       0.8143930000,
					Timestamp:    time.Unix(1484543337, 0),
				},
				{
					UID:          "8ad9938d00e0c05f9088da603275d9d0",
					Transformer:  transactions.TransformTypeKraken,
					Currency:     "BTC",
					DetectedType: transactions.TypeFee,
					Amount:       0.0013030000,
					Timestamp:    time.Unix(1484543337, 0),
				},
			},
		},
		{
			name: "kraken airdrop",
			row:  []string{"L5M3LT-5ZXMV-6HY4DR", "LA6H43Q-QITN2-Z6JZGC", "2017-08-01 16:26:36", "transfer", "", "currency", "BCH", "0.6529973350", "0.0000000000", "0.6529973350"},
			expected: []transactions.Transaction{
				{
					UID:          "0abc00109247ccf2345f27d905bb1940",
					Transformer:  transactions.TransformTypeKraken,
					Currency:     "BCH",
					DetectedType: transactions.TypeAirdrop,
					Amount:       0.6529973350,
					Timestamp:    time.Unix(1501604796, 0),
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := sources.KrakenSource{}.TransformRow(tc.row)
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
					Timestamp:    time.Unix(1624989420, 0),
				},
				{
					UID:          "44de929a7ebca4380db896d63fb063f0",
					Transformer:  transactions.TransformTypeCoinomi,
					Currency:     "ETH",
					DetectedType: transactions.TypeFee,
					Amount:       0.00058689,
					Timestamp:    time.Unix(1624989420, 0),
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
					Timestamp:    time.Unix(1624989420, 0),
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
					Timestamp:    time.Unix(1612849620, 0),
				},
				{
					UID:          "8aab065f6d121deb579346e4ece74635",
					Transformer:  transactions.TransformTypeCoinomi,
					Currency:     "XRP",
					DetectedType: transactions.TypeFee,
					Amount:       0.000012,
					Timestamp:    time.Unix(1612849620, 0),
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
					Timestamp:    time.Unix(1453347106, 0),
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
					Timestamp:    time.Unix(1610107160, 0),
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
					Timestamp:    time.Unix(1629962522, 0),
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
					Timestamp:    time.Unix(1484714102, 0),
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
					Timestamp:    time.Unix(1619701433, 0),
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
