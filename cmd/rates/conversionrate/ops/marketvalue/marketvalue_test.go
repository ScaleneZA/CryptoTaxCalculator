package marketvalue

import (
	"database/sql"
	conversionrate2 "github.com/ScaleneZA/CryptoTaxCalculator/cmd/rates/conversionrate"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/rates/conversionrate/db/markets"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/rates/di"
	"github.com/luno/jettison/jtest"
	"github.com/stretchr/testify/require"
	"math"
	"testing"
)

func TestFindClosest(t *testing.T) {
	b := di.SetupDIForTesting()
	seedData(t, b.DB())

	testCases := []struct {
		name        string
		pair        conversionrate2.Pair
		timestamp   int64
		expected    *conversionrate2.MarketPair
		expectedErr error
	}{
		{
			name:      "find closest before",
			pair:      conversionrate2.PairUSDBTC,
			timestamp: 1689375600001,
			expected: &conversionrate2.MarketPair{
				Pair: conversionrate2.PairUSDBTC,
				MarketSlice: conversionrate2.MarketSlice{
					Timestamp: 1689375600000,
					Open:      30270.01,
					High:      30341.7,
					Low:       30250.01,
					Close:     30336.28,
				},
			},
		},
		{
			name:      "find closest after",
			pair:      conversionrate2.PairUSDBTC,
			timestamp: 1689368399999,
			expected: &conversionrate2.MarketPair{
				Pair: conversionrate2.PairUSDBTC,
				MarketSlice: conversionrate2.MarketSlice{
					Timestamp: 1689368400000,
					Open:      30205.98,
					High:      30295.32,
					Low:       30175.01,
					Close:     30263.39,
				},
			},
		},
		{
			name:        "closest rates are stale",
			pair:        conversionrate2.PairUSDBTC,
			timestamp:   1234,
			expected:    nil,
			expectedErr: conversionrate2.ErrStoredRateExceedsThreshold,
		},
		{
			name:        "no market for pair",
			pair:        conversionrate2.PairUSDLTC,
			timestamp:   1234,
			expected:    nil,
			expectedErr: conversionrate2.ErrNoMarket,
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
	b := di.SetupDIForTesting()
	seedData(t, b.DB())

	testCases := []struct {
		name      string
		from      string
		to        string
		timestamp int64
		expected  float64
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
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := ValueAtTime(b, tc.from, tc.to, tc.timestamp)
			require.Nil(t, err)
			require.Equal(t, tc.expected, math.Round(actual*100)/100)
		})
	}

}

func seedData(t *testing.T, dbc *sql.DB) {
	mps := []conversionrate2.MarketPair{
		{
			Pair: conversionrate2.PairUSDBTC,
			MarketSlice: conversionrate2.MarketSlice{
				Timestamp: 1689375600000,
				Open:      30270.01,
				High:      30341.7,
				Low:       30250.01,
				Close:     30336.28,
			},
		},
		{
			Pair: conversionrate2.PairUSDBTC,
			MarketSlice: conversionrate2.MarketSlice{
				Timestamp: 1689372000000,
				Open:      30263.39,
				High:      30280.28,
				Low:       30233.52,
				Close:     30270.01,
			},
		},
		{
			Pair: conversionrate2.PairUSDBTC,
			MarketSlice: conversionrate2.MarketSlice{
				Timestamp: 1689368400000,
				Open:      30205.98,
				High:      30295.32,
				Low:       30175.01,
				Close:     30263.39,
			},
		},
		{
			Pair: conversionrate2.PairZARUSD,
			MarketSlice: conversionrate2.MarketSlice{
				Timestamp: 1689375600000,
				Open:      19.10,
				High:      19.56,
				Low:       19.09,
				Close:     19.40,
			},
		},
		{
			Pair: conversionrate2.PairZARUSD,
			MarketSlice: conversionrate2.MarketSlice{
				Timestamp: 1689372000000,
				Open:      20.10,
				High:      20.12,
				Low:       19.05,
				Close:     19.10,
			},
		},
		{
			Pair: conversionrate2.PairZARUSD,
			MarketSlice: conversionrate2.MarketSlice{
				Timestamp: 1689368400000,
				Open:      21.50,
				High:      21.56,
				Low:       20.06,
				Close:     20.10,
			},
		},
	}

	for _, mp := range mps {
		_, err := markets.CreateIgnoreDuplicate(dbc, mp)
		require.Nil(t, err)
	}
}
