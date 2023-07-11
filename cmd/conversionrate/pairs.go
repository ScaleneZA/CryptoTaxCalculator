package conversionrate

import (
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/conversionrate/reader"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/conversionrate/writer"
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
				Location: "https://www.cryptodatadownload.com/cdd/Gemini_BTCUSD_d.csv",
			},
			writer: writer.FileWriter{
				Filename: "USD_BTC.csv",
			},
		},
	},
}
