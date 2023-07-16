package marketvalue

import (
	"database/sql"
	"errors"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/conversionrate/db/markets"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/conversionrate/sharedtypes"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/di"
	"github.com/stretchr/testify/require"
	"math"
	"testing"
)

func TestFindClosest(t *testing.T) {
	b := di.SetupDIForTesting()
	seedData(t, b.DB())

	testCases := []struct {
		name        string
		pair        sharedtypes.Pair
		timestamp   int64
		expected    *sharedtypes.MarketPair
		expectedErr error
	}{
		{
			name:      "find closest before",
			pair:      sharedtypes.PairUSDBTC,
			timestamp: 1689375600001,
			expected: &sharedtypes.MarketPair{
				Pair: sharedtypes.PairUSDBTC,
				MarketSlice: sharedtypes.MarketSlice{
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
			pair:      sharedtypes.PairUSDBTC,
			timestamp: 1689368399999,
			expected: &sharedtypes.MarketPair{
				Pair: sharedtypes.PairUSDBTC,
				MarketSlice: sharedtypes.MarketSlice{
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
			pair:        sharedtypes.PairUSDBTC,
			timestamp:   1234,
			expected:    nil,
			expectedErr: errors.New("closest timestamps of stored rates exceed threshold of 1 week"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := FindClosest(b, tc.pair, tc.timestamp)
			require.Equal(t, tc.expectedErr, err)
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
	mps := []sharedtypes.MarketPair{
		{
			Pair: sharedtypes.PairUSDBTC,
			MarketSlice: sharedtypes.MarketSlice{
				Timestamp: 1689375600000,
				Open:      30270.01,
				High:      30341.7,
				Low:       30250.01,
				Close:     30336.28,
			},
		},
		{
			Pair: sharedtypes.PairUSDBTC,
			MarketSlice: sharedtypes.MarketSlice{
				Timestamp: 1689372000000,
				Open:      30263.39,
				High:      30280.28,
				Low:       30233.52,
				Close:     30270.01,
			},
		},
		{
			Pair: sharedtypes.PairUSDBTC,
			MarketSlice: sharedtypes.MarketSlice{
				Timestamp: 1689368400000,
				Open:      30205.98,
				High:      30295.32,
				Low:       30175.01,
				Close:     30263.39,
			},
		},
		{
			Pair: sharedtypes.PairZARUSD,
			MarketSlice: sharedtypes.MarketSlice{
				Timestamp: 1689375600000,
				Open:      19.10,
				High:      19.56,
				Low:       19.09,
				Close:     19.40,
			},
		},
		{
			Pair: sharedtypes.PairZARUSD,
			MarketSlice: sharedtypes.MarketSlice{
				Timestamp: 1689372000000,
				Open:      20.10,
				High:      20.12,
				Low:       19.05,
				Close:     19.10,
			},
		},
		{
			Pair: sharedtypes.PairZARUSD,
			MarketSlice: sharedtypes.MarketSlice{
				Timestamp: 1689368400000,
				Open:      21.50,
				High:      21.56,
				Low:       20.06,
				Close:     20.10,
			},
		},
	}

	for _, mp := range mps {
		_, err := markets.Create(dbc, mp)
		require.Nil(t, err)
	}
}
