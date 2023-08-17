package marketvalue

import (
	"database/sql"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/rates/conversionrate"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/rates/conversionrate/db/markets"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/rates/di"
	"github.com/luno/jettison/jtest"
	"github.com/stretchr/testify/require"
	"math"
	"testing"
)

var b = di.SetupDIForTesting()

func TestFindClosest(t *testing.T) {
	seedData(t, b.DB())

	testCases := []struct {
		name        string
		pair        conversionrate.Pair
		timestamp   int64
		expected    *conversionrate.MarketPair
		expectedErr error
	}{
		{
			name:      "find closest before",
			pair:      conversionrate.PairUSDBTC,
			timestamp: 1689375600001,
			expected: &conversionrate.MarketPair{
				Pair: conversionrate.PairUSDBTC,
				MarketSlice: conversionrate.MarketSlice{
					Timestamp: 1689375600000,
					Open:      30270.01,
					High:      30341.7,
					Low:       30250.01,
					Close:     30336.28,
				},
			},
		},
		{
			name:      "find closest before again",
			pair:      conversionrate.PairUSDBTC,
			timestamp: 1689372000001,
			expected: &conversionrate.MarketPair{
				Pair: conversionrate.PairUSDBTC,
				MarketSlice: conversionrate.MarketSlice{
					Timestamp: 1689372000000,
					Open:      30263.39,
					High:      30280.28,
					Low:       30233.52,
					Close:     30270.01,
				},
			},
		},
		{
			name:      "find closest after",
			pair:      conversionrate.PairUSDBTC,
			timestamp: 1689368399999,
			expected: &conversionrate.MarketPair{
				Pair: conversionrate.PairUSDBTC,
				MarketSlice: conversionrate.MarketSlice{
					Timestamp: 1689368400000,
					Open:      30205.98,
					High:      30295.32,
					Low:       30175.01,
					Close:     30263.39,
				},
			},
		},
		{
			name:      "find closest after again",
			pair:      conversionrate.PairUSDBTC,
			timestamp: 1689371999999,
			expected: &conversionrate.MarketPair{
				Pair: conversionrate.PairUSDBTC,
				MarketSlice: conversionrate.MarketSlice{
					Timestamp: 1689372000000,
					Open:      30263.39,
					High:      30280.28,
					Low:       30233.52,
					Close:     30270.01,
				},
			},
		},
		{
			name:        "no market for pair",
			pair:        conversionrate.PairUSDLTC,
			timestamp:   1234,
			expected:    nil,
			expectedErr: conversionrate.ErrNoMarket,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := FindClosest(b, tc.pair, tc.timestamp)
			jtest.Require(t, tc.expectedErr, err)
			require.Equal(t, tc.expected, actual)
		})
	}

}

func TestValueAtTime(t *testing.T) {
	seedData(t, b.DB())

	testCases := []struct {
		name        string
		from        string
		to          string
		timestamp   int64
		expected    float64
		expectedErr error
	}{
		{
			name:      "One hop",
			from:      "USD",
			to:        "BTC",
			timestamp: 1689375600001,
			expected:  30336.28,
		},
		{
			name:      "Two hops",
			from:      "ZAR",
			to:        "BTC",
			timestamp: 1689375600001,
			expected:  588523.83,
		},
		{
			name:        "A rate is stale",
			from:        "ZAR",
			to:          "BTC",
			timestamp:   1234,
			expectedErr: conversionrate.ErrNoRatesFound,
		},
		{
			name:        "Rate not found",
			from:        "USD",
			to:          "MOON",
			timestamp:   1689375600001,
			expectedErr: conversionrate.ErrNoRatesFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := ValueAtTime(b, tc.from, tc.to, tc.timestamp)
			jtest.Require(t, tc.expectedErr, err)
			require.Equal(t, tc.expected, math.Round(actual*100)/100)
		})
	}

}

func seedData(t *testing.T, dbc *sql.DB) {
	mps := []conversionrate.MarketPair{
		{
			Pair: conversionrate.PairUSDBTC,
			MarketSlice: conversionrate.MarketSlice{
				Timestamp: 1689375600000,
				Open:      30270.01,
				High:      30341.7,
				Low:       30250.01,
				Close:     30336.28,
			},
		},
		{
			Pair: conversionrate.PairUSDBTC,
			MarketSlice: conversionrate.MarketSlice{
				Timestamp: 1689372000000,
				Open:      30263.39,
				High:      30280.28,
				Low:       30233.52,
				Close:     30270.01,
			},
		},
		{
			Pair: conversionrate.PairUSDBTC,
			MarketSlice: conversionrate.MarketSlice{
				Timestamp: 1689368400000,
				Open:      30205.98,
				High:      30295.32,
				Low:       30175.01,
				Close:     30263.39,
			},
		},
		{
			Pair: conversionrate.PairZARUSD,
			MarketSlice: conversionrate.MarketSlice{
				Timestamp: 1689375600000,
				Open:      19.10,
				High:      19.56,
				Low:       19.09,
				Close:     19.40,
			},
		},
		{
			Pair: conversionrate.PairZARUSD,
			MarketSlice: conversionrate.MarketSlice{
				Timestamp: 1689372000000,
				Open:      20.10,
				High:      20.12,
				Low:       19.05,
				Close:     19.10,
			},
		},
		{
			Pair: conversionrate.PairZARUSD,
			MarketSlice: conversionrate.MarketSlice{
				Timestamp: 1689368400000,
				Open:      21.50,
				High:      21.56,
				Low:       20.06,
				Close:     20.10,
			},
		},
	}

	for _, mp := range mps {
		markets.CreateIgnoreDuplicate(dbc, mp)
	}
}
