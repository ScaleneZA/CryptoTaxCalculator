package sync

import (
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/conversionrate/sync/reader"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/conversionrate/sync/writer"
)

var (
	PairUSDBTC = Pair{
		currency1: "USD",
		currency2: "BTC",
	}
)

var PairSyncers = map[Pair][]syncer{
	PairUSDBTC: {
		syncer{
			reader: reader.HttpReader{
				Location: "https://www.cryptodatadownload.com/cdd/Gemini_BTCUSD_1h.csv",
			},
			writer: writer.FileWriter{
				Filename: "USD_BTC.csv",
			},
		},
	},
}
