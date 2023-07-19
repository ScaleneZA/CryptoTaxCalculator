package syncer_test

import (
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/rates/conversionrate"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/rates/conversionrate/db/markets"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/rates/conversionrate/ops/sync/readtransformer"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/rates/conversionrate/ops/sync/readtransformer/csvreader"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/rates/conversionrate/ops/sync/syncer"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/rates/di"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHolisticSyncer_Sync(t *testing.T) {
	b := di.SetupDIForTesting()

	mps, err := markets.ListAll(b.DB())
	require.Nil(t, err)
	require.Equal(t, []conversionrate.MarketPair(nil), mps)

	s := syncer.HolisticSyncer{
		ReadTransformer: readtransformer.GeminiCSV{
			Reader: csvreader.LocalCSVReader{
				Location: "../readtransformer/csvreader/test_data/Gemini_BTCUSD_1h.csv",
				SkipRows: 2,
			},
		},
	}

	err = s.Sync(b)
	require.Nil(t, err)

	actual, err := markets.ListAll(b.DB())
	require.Nil(t, err)

	require.Equal(t, []conversionrate.MarketPair{
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
	}, actual)
}
