package writer_test

import (
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/conversionrate/sharedtypes"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/conversionrate/sync/writer"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/db/markets"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/di"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSQLLiteWriter_WriteAll(t *testing.T) {
	w := writer.SQLLiteWriter{
		Pair: sharedtypes.Pair{
			Currency2: "BTC",
			Currency1: "USD",
		},
	}

	b := di.SetupDIForTesting()

	data := []sharedtypes.MarketSlice{
		{
			Timestamp: 12345,
			Open:      1.23,
			High:      2.34,
			Low:       3.45,
			Close:     4.56,
		},
		{
			Timestamp: 234567,
			Open:      5.67,
			High:      6.78,
			Low:       7.89,
			Close:     8.90,
		},
	}

	err := w.WriteAll(b, data)
	require.Nil(t, err)

	mps, err := markets.ListAll(b.DB())
	require.Nil(t, err)

	require.Equal(t, data, mps)
}

func TestSQLLiteWriter_DeleteAll(t *testing.T) {

}
