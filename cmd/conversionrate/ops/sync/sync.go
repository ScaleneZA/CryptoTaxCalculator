package sync

import (
	"errors"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/conversionrate/ops/sync/readtransformer"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/conversionrate/ops/sync/readtransformer/csvreader"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/conversionrate/ops/sync/syncer"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/conversionrate/sharedtypes"
	"log"
)

var PairSyncers = map[sharedtypes.Pair][]syncer.HolisticSyncer{
	sharedtypes.PairUSDBTC: {
		syncer.HolisticSyncer{
			ReadTransformer: readtransformer.GeminiCSV{
				Reader: csvreader.HTTPCSVReader{
					Location: "https://www.cryptodatadownload.com/cdd/Gemini_BTCUSD_1h.csv",
					SkipRows: 2,
				},
			},
		},
	},
	sharedtypes.PairUSDETH: {
		syncer.HolisticSyncer{
			ReadTransformer: readtransformer.GeminiCSV{
				Reader: csvreader.HTTPCSVReader{
					Location: "https://www.cryptodatadownload.com/cdd/Gemini_ETHUSD_1h.csv",
					SkipRows: 2,
				},
			},
		},
	},
}

func SyncAll(b Backends) error {
	failedSyncs := 0
	for p, syncers := range PairSyncers {
		successful := false
		for _, s := range syncers {
			err := s.Sync(b)
			if err != nil {
				log.Println("Failed Sync for " + p.FromCurrency + p.ToCurrency + ", failing over...")
				log.Println(err.Error())
			} else {
				// Successful sync, no need to fail over.
				successful = true
				break
			}
		}
		if !successful {
			failedSyncs++
		}
	}

	if failedSyncs > 0 {
		return errors.New("one or more pairs failed to sync")
	}

	return nil
}
