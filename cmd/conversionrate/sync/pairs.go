package sync

import (
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/conversionrate/sharedtypes"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/conversionrate/sync/readtransformer"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/conversionrate/sync/readtransformer/csvreader"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/conversionrate/sync/writer"
)

var (
	PairUSDBTC = sharedtypes.Pair{
		Currency1: "USD",
		Currency2: "BTC",
	}
)

var PairSyncers = map[sharedtypes.Pair][]syncer{
	PairUSDBTC: {
		syncer{
			readTransformer: readtransformer.GeminiCSV{
				Reader: csvreader.HTTPCSVReader{
					Location: "https://www.cryptodatadownload.com/cdd/Gemini_BTCUSD_1h.csv",
					SkipRows: 2,
				},
			},
			writer: writer.SQLLiteWriter{
				FromCurrency: "USD",
				ToCurrency:   "BTC",
			},
		},
	},
}
