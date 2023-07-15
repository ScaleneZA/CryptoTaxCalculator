package marketvalue

import (
	"database/sql"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/conversionrate/db/markets"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/conversionrate/sharedtypes"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/di"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFindClosest(t *testing.T) {
	b := di.SetupDIForTesting()
	seedData(t, b.DB())

	testCases := []struct {
		name      string
		pair      sharedtypes.Pair
		timestamp int64
		expected  *sharedtypes.MarketPair
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
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := FindClosest(b, tc.pair, tc.timestamp)
			require.Nil(t, err)
			require.Equal(t, tc.expected, actual)
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
	}

	for _, mp := range mps {
		_, err := markets.Create(dbc, mp)
		require.Nil(t, err)
	}
}
