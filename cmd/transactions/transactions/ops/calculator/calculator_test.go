package calculator_test

import (
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/rates/conversionrate"
	rates_mock "github.com/ScaleneZA/CryptoTaxCalculator/cmd/rates/conversionrate/client/mockery"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/transactions/di"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/transactions/transactions"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/transactions/transactions/ops/calculator"
	"github.com/luno/jettison/jtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestCalculate(t *testing.T) {
	testCases := []struct {
		name          string
		seed          []transactions.Transaction
		rateMockCalls []*mock.Call
		expected      calculator.YearEndTotals
		expectedErr   error
	}{
		{
			name: "Happy Path",
			seed: []transactions.Transaction{
				{
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
					Currency:     "BTC",
					DetectedType: transactions.TypeBuy,
					Amount:       0.5,
					Timestamp:    1535450915,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 900,
					},
				},
				{
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
					Currency:      "ETH",
					DetectedType:  transactions.TypeBuy,
					OverridedType: transactions.TypeSell,
					Amount:        0.25,
					Timestamp:     1656410903,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 300,
					},
				},
				{
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
					Currency:     "BTC",
					DetectedType: transactions.TypeSell,
					Amount:       0.4,
					Timestamp:    1705835912,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 2000,
					},
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
		{
			name: "Rate not included, calls rates client",
			seed: []transactions.Transaction{
				{
					Currency:     "ETH",
					DetectedType: transactions.TypeBuy,
					Amount:       1,
					Timestamp:    1519812502,
				},
				{
					Currency:     "ETH",
					DetectedType: transactions.TypeSell,
					Amount:       0.4,
					Timestamp:    1705835912,
				},
			},
			rateMockCalls: []*mock.Call{
				new(mock.Mock).On("ValueAtTime", "ZAR", "ETH", int64(1519812502)).Return(float64(100), nil),
				new(mock.Mock).On("ValueAtTime", "ZAR", "ETH", int64(1705835912)).Return(float64(1800), nil),
			},
			expected: calculator.YearEndTotals{
				2024: {
					"ETH":   680,
					"TOTAL": 680,
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
					Timestamp:    1519812502,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "USD",
						Price: 2000,
					},
				},
				{
					Currency:     "ETH",
					DetectedType: transactions.TypeSell,
					Amount:       0.4,
					Timestamp:    1705835912,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "USD",
						Price: 2000,
					},
				},
			},
			rateMockCalls: []*mock.Call{
				new(mock.Mock).On("ValueAtTime", "ZAR", "ETH", int64(1519812502)).Return(float64(100), nil),
				new(mock.Mock).On("ValueAtTime", "ZAR", "ETH", int64(1705835912)).Return(float64(1800), nil),
			},
			expected: calculator.YearEndTotals{
				2024: {
					"ETH":   680,
					"TOTAL": 680,
				},
			},
		},
		{
			name: "Rate client returns error",
			seed: []transactions.Transaction{
				{
					Currency:     "ETH",
					DetectedType: transactions.TypeBuy,
					Amount:       1,
					Timestamp:    1519812502,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "USD",
						Price: 2000,
					},
				},
				{
					Currency:     "ETH",
					DetectedType: transactions.TypeSell,
					Amount:       0.4,
					Timestamp:    1705835912,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "USD",
						Price: 2000,
					},
				},
			},
			rateMockCalls: []*mock.Call{
				new(mock.Mock).On("ValueAtTime", "ZAR", "ETH", int64(1519812502)).Return(float64(0), conversionrate.ErrNoRatesFound),
			},
			expectedErr: conversionrate.ErrNoRatesFound,
		},
		{
			name: "Invalid Transaction Order",
			seed: []transactions.Transaction{
				{
					Currency:     "BTC",
					DetectedType: transactions.TypeBuy,
					Amount:       0.5,
					Timestamp:    1535450915,
					WholePriceAtPoint: transactions.FiatPrice{
						Fiat:  "ZAR",
						Price: 900,
					},
				},
				{
					Currency:     "BTC",
					DetectedType: transactions.TypeBuy,
					Amount:       0.56,
					Timestamp:    1519812503,
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
			assert.Equal(t, tc.expected, actual)
		})
	}
}
