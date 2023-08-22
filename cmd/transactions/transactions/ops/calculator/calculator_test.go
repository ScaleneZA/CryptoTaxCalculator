package calculator_test

import (
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/rates/conversionrate"
	rates_mock "github.com/ScaleneZA/CryptoTaxCalculator/cmd/rates/conversionrate/client/mockery"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/transactions/db"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/transactions/di"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/transactions/transactions"
	calculatordb "github.com/ScaleneZA/CryptoTaxCalculator/cmd/transactions/transactions/db/calculator"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/transactions/transactions/ops/calculator"
	"github.com/luno/jettison/errors"
	"github.com/luno/jettison/jtest"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCalculate(t *testing.T) {
	testCases := []struct {
		name          string
		seed          []transactions.Transaction
		rateMockCalls []*mock.Call
		expected      []calculator.YearEndTotal
		expectedErr   error
	}{
		{
			name: "Happy Path",
			seed: []transactions.Transaction{
				{
					Currency:     "ETH",
					DetectedType: transactions.TypeBuy,
					Amount:       0.56,
					Timestamp:    time.Unix(1519812503, 0),
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 100,
					},
				},
				{
					Currency:     "BTC",
					DetectedType: transactions.TypeBuy,
					Amount:       0.5,
					Timestamp:    time.Unix(1535450915, 0),
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 900,
					},
				},
				{
					Currency:     "ETH",
					DetectedType: transactions.TypeBuy,
					Amount:       1.2,
					Timestamp:    time.Unix(1535450903, 0),
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 200,
					},
				},
				{
					Currency:      "ETH",
					DetectedType:  transactions.TypeBuy,
					OverridedType: transactions.TypeSell,
					Amount:        0.25,
					Timestamp:     time.Unix(1656410903, 0),
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 300,
					},
				},
				{
					Currency:     "ETH",
					DetectedType: transactions.TypeSell,
					Amount:       1.25,
					Timestamp:    time.Unix(1687946903, 0),
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 400,
					},
				},
				{
					Currency:     "BTC",
					DetectedType: transactions.TypeSell,
					Amount:       0.4,
					Timestamp:    time.Unix(1705835912, 0),
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 2000,
					},
				},
			},
			expected: []calculator.YearEndTotal{
				{
					Year: 2018,
					Gains: []calculator.Gain{
						{
							Asset: "TOTAL",
						},
					},
					Balances: []calculator.Balance{
						{
							Asset:  "ETH",
							Amount: 0.56,
						},
					},
				},
				{
					Year: 2019,
					Gains: []calculator.Gain{
						{
							Asset: "TOTAL",
						},
					},
					Balances: []calculator.Balance{
						{
							Asset:  "ETH",
							Amount: 1.76,
						},
						{
							Asset:  "BTC",
							Amount: 0.5,
						},
					},
				},
				{
					Year: 2020,
					Gains: []calculator.Gain{
						{
							Asset: "TOTAL",
						},
					},
					Balances: []calculator.Balance{
						{
							Asset:  "ETH",
							Amount: 1.76,
						},
						{
							Asset:  "BTC",
							Amount: 0.5,
						},
					},
				},
				{
					Year: 2021,
					Gains: []calculator.Gain{
						{
							Asset: "TOTAL",
						},
					},
					Balances: []calculator.Balance{
						{
							Asset:  "ETH",
							Amount: 1.76,
						},
						{
							Asset:  "BTC",
							Amount: 0.5,
						},
					},
				},
				{
					Year: 2022,
					Gains: []calculator.Gain{
						{
							Asset: "TOTAL",
						},
					},
					Balances: []calculator.Balance{
						{
							Asset:  "ETH",
							Amount: 1.76,
						},
						{
							Asset:  "BTC",
							Amount: 0.5,
						},
					},
				},
				{
					Year: 2023,
					Gains: []calculator.Gain{
						{
							Asset:    "ETH",
							Amount:   0.25,
							Costs:    25,
							Proceeds: 75,
						},
						{
							Asset:    "TOTAL",
							Costs:    25,
							Proceeds: 75,
						},
					},
					Balances: []calculator.Balance{
						{
							Asset:  "ETH",
							Amount: 1.51,
						},
						{
							Asset:  "BTC",
							Amount: 0.5,
						},
					},
				},
				{
					Year: 2024,
					Gains: []calculator.Gain{
						{
							Asset:    "TOTAL",
							Costs:    579,
							Proceeds: 1300,
						},
						{
							Asset:    "BTC",
							Amount:   0.4,
							Costs:    360,
							Proceeds: 800,
						},
						{
							Asset:    "ETH",
							Amount:   1.25,
							Costs:    219,
							Proceeds: 500,
						},
					},
					Balances: []calculator.Balance{
						{
							Asset:  "ETH",
							Amount: 0.26,
						},
						{
							Asset:  "BTC",
							Amount: 0.1,
						},
					},
				},
			},
		},
		{
			name: "Rate not included, calls rates client",
			seed: []transactions.Transaction{
				{
					Currency:     "ETH",
					DetectedType: transactions.TypeBuy,
					Amount:       1,
					Timestamp:    time.Unix(1705835901, 0),
				},
				{
					Currency:     "ETH",
					DetectedType: transactions.TypeSell,
					Amount:       0.4,
					Timestamp:    time.Unix(1705835912, 0),
				},
			},
			rateMockCalls: []*mock.Call{
				new(mock.Mock).On("ValueAtTime", "ZAR", "ETH", int64(1705835901)).Return(float64(100), nil),
				new(mock.Mock).On("ValueAtTime", "ZAR", "ETH", int64(1705835912)).Return(float64(1800), nil),
			},
			expected: []calculator.YearEndTotal{
				{
					Year: 2024,
					Gains: []calculator.Gain{
						{
							Asset:    "ETH",
							Amount:   0.4,
							Costs:    40,
							Proceeds: 720,
						},
						{
							Asset:    "TOTAL",
							Costs:    40,
							Proceeds: 720,
						},
					},
					Balances: []calculator.Balance{
						{
							Asset:  "ETH",
							Amount: 0.6,
						},
					},
				},
			},
		},
		{
			name: "Rate fiat incorrect, calls rates client",
			seed: []transactions.Transaction{
				{
					Currency:     "ETH",
					DetectedType: transactions.TypeBuy,
					Amount:       1,
					Timestamp:    time.Unix(1705835901, 0),
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "USD",
						Price: 2000,
					},
				},
				{
					Currency:     "ETH",
					DetectedType: transactions.TypeSell,
					Amount:       0.4,
					Timestamp:    time.Unix(1705835912, 0),
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "USD",
						Price: 2000,
					},
				},
			},
			rateMockCalls: []*mock.Call{
				new(mock.Mock).On("ValueAtTime", "ZAR", "ETH", int64(1705835901)).Return(float64(100), nil),
				new(mock.Mock).On("ValueAtTime", "ZAR", "ETH", int64(1705835912)).Return(float64(1800), nil),
			},
			expected: []calculator.YearEndTotal{
				{
					Year: 2024,
					Gains: []calculator.Gain{
						{
							Asset:    "ETH",
							Amount:   0.4,
							Costs:    40,
							Proceeds: 720,
						},
						{
							Asset:    "TOTAL",
							Costs:    40,
							Proceeds: 720,
						},
					},
					Balances: []calculator.Balance{
						{
							Asset:  "ETH",
							Amount: 0.6,
						},
					},
				},
			},
		},
		{
			name: "Rate client returns unexpected error",
			seed: []transactions.Transaction{
				{
					Currency:     "ETH",
					DetectedType: transactions.TypeBuy,
					Amount:       1,
					Timestamp:    time.Unix(1519812502, 0),
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "USD",
						Price: 2000,
					},
				},
				{
					Currency:     "ETH",
					DetectedType: transactions.TypeSell,
					Amount:       0.4,
					Timestamp:    time.Unix(1705835912, 0),
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "USD",
						Price: 2000,
					},
				},
			},
			rateMockCalls: []*mock.Call{
				new(mock.Mock).On("ValueAtTime", "ZAR", "ETH", int64(1519812502)).Return(float64(0), errors.New("unexpected")),
			},
			expectedErr: errors.New("unexpected"),
		},
		{
			name: "Rate client returns expected error, skip currency",
			seed: []transactions.Transaction{
				{
					Currency:     "ETH",
					DetectedType: transactions.TypeBuy,
					Amount:       1,
					Timestamp:    time.Unix(1705835901, 0),
				},
				{
					Currency:     "BTC",
					DetectedType: transactions.TypeBuy,
					Amount:       1,
					Timestamp:    time.Unix(1705835902, 0),
				},
				{
					Currency:     "ETH",
					DetectedType: transactions.TypeSell,
					Amount:       0.4,
					Timestamp:    time.Unix(1705835912, 0),
				},
				{
					Currency:     "BTC",
					DetectedType: transactions.TypeSell,
					Amount:       0.4,
					Timestamp:    time.Unix(1705835913, 0),
				},
			},
			rateMockCalls: []*mock.Call{
				new(mock.Mock).On("ValueAtTime", "ZAR", "ETH", int64(1705835901)).Return(float64(100), nil),
				new(mock.Mock).On("ValueAtTime", "ZAR", "BTC", int64(1705835902)).Return(float64(0), conversionrate.ErrNoRatesFound),
				new(mock.Mock).On("ValueAtTime", "ZAR", "ETH", int64(1705835912)).Return(float64(1800), nil),
				new(mock.Mock).On("ValueAtTime", "ZAR", "BTC", int64(1705835913)).Return(float64(0), conversionrate.ErrNoRatesFound),
			},
			expected: []calculator.YearEndTotal{
				{
					Year: 2024,
					Gains: []calculator.Gain{
						{
							Asset:    "ETH",
							Amount:   0.4,
							Costs:    40,
							Proceeds: 720,
						},
						{
							Asset:    "TOTAL",
							Costs:    40,
							Proceeds: 720,
						},
					},
					Balances: []calculator.Balance{
						{
							Asset:  "ETH",
							Amount: 0.6,
						},
					},
				},
			},
		},
		{
			name: "Invalid Transaction Order",
			seed: []transactions.Transaction{
				{
					Currency:     "BTC",
					DetectedType: transactions.TypeBuy,
					Amount:       0.5,
					Timestamp:    time.Unix(1535450915, 0),
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 900,
					},
				},
				{
					Currency:     "BTC",
					DetectedType: transactions.TypeBuy,
					Amount:       0.56,
					Timestamp:    time.Unix(1519812503, 0),
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 100,
					},
				},
			},
			expectedErr: transactions.ErrInvalidTransactionOrder,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rc := new(rates_mock.Client)
			rc.ExpectedCalls = append(rc.ExpectedCalls, tc.rateMockCalls...)

			b := di.SetupDIForTesting(di.BackendsTest{
				RatesClient: rc,
			})

			actual, err := calculator.Calculate(b, "ZAR", tc.seed)
			jtest.Require(t, tc.expectedErr, err)

			require.Equal(t, len(tc.expected), len(actual))
			for i, e := range tc.expected {

				require.Equal(t, len(e.Balances), len(actual[i].Balances))
				for j, eb := range e.Balances {
					require.Equal(t, eb.Asset, actual[i].Balances[j].Asset)
					require.InDelta(t, eb.Amount, actual[i].Balances[j].Amount, 1e-10)
				}

				require.Equal(t, len(e.Gains), len(actual[i].Gains))
				for j, eb := range e.Gains {
					require.Equal(t, eb.Asset, actual[i].Gains[j].Asset)
					require.InDelta(t, eb.Amount, actual[i].Gains[j].Amount, 1e-10)
					require.InDelta(t, eb.Costs, actual[i].Gains[j].Costs, 1e-10)
					require.InDelta(t, eb.Proceeds, actual[i].Gains[j].Proceeds, 1e-10)
				}
			}
		})
	}
}

func TestPopulateOverriddenTypes(t *testing.T) {
	dbc := db.ConnectForTesting()
	b := di.SetupDIForTesting(di.BackendsTest{
		DB: dbc,
	})

	_, err := calculatordb.Upsert(dbc, "1234", transactions.TypeAirdrop)
	jtest.RequireNil(t, err)
	_, err = calculatordb.Upsert(dbc, "1235", transactions.TypeFee)
	jtest.RequireNil(t, err)
	_, err = calculatordb.Upsert(dbc, "1236", transactions.TypeSell)
	jtest.RequireNil(t, err)

	tests := []struct {
		name     string
		ts       []transactions.Transaction
		expected []transactions.Transaction
	}{
		{
			name: "golden path",
			ts: []transactions.Transaction{
				{
					UID: "1234",
				},
				{
					UID:           "1236",
					OverridedType: transactions.TypeInterest,
				},
				{
					UID:           "1237",
					OverridedType: transactions.TypeInterest,
				},
			},
			expected: []transactions.Transaction{
				{
					UID:           "1234",
					OverridedType: transactions.TypeAirdrop,
				},
				{
					UID:           "1236",
					OverridedType: transactions.TypeSell,
				},
				{
					UID:           "1237",
					OverridedType: transactions.TypeInterest,
				},
			},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := calculator.PopulateOverriddenTypes(b, tt.ts)
			jtest.RequireNil(t, err)

			require.Equal(t, tt.expected, actual)
		})
	}
}
