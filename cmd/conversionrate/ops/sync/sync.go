package sync

import (
	"errors"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/conversionrate/ops/sync/readtransformer"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/conversionrate/ops/sync/readtransformer/csvreader"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/conversionrate/ops/sync/syncer"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/conversionrate/sharedtypes"
	"log"
)

var PairSyncers = map[sharedtypes.Pair][]syncer.Syncer{
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
	sharedtypes.PairUSDLTC: {
		syncer.HolisticSyncer{
			ReadTransformer: readtransformer.GeminiCSV{
				Reader: csvreader.HTTPCSVReader{
					Location: "https://www.cryptodatadownload.com/cdd/Gemini_LTCUSD_1h.csv",
					SkipRows: 2,
				},
			},
		},
	},
	sharedtypes.PairUSDBCH: {
		syncer.HolisticSyncer{
			ReadTransformer: readtransformer.GeminiCSV{
				Reader: csvreader.HTTPCSVReader{
					Location: "https://www.cryptodatadownload.com/cdd/Gemini_BCHUSD_d.csv",
					SkipRows: 2,
				},
			},
		},
		syncer.HolisticSyncer{
			ReadTransformer: readtransformer.GeminiCSV{
				Reader: csvreader.HTTPCSVReader{
					Location: "https://www.cryptodatadownload.com/cdd/Gemini_BCHUSD_1h.csv",
					SkipRows: 2,
				},
			},
		},
	},
	sharedtypes.PairUSDBAT: {
		syncer.HolisticSyncer{
			ReadTransformer: readtransformer.GeminiCSV{
				Reader: csvreader.HTTPCSVReader{
					Location: "https://www.cryptodatadownload.com/cdd/Gemini_BATUSD_d.csv",
					SkipRows: 2,
				},
			},
		},
		syncer.HolisticSyncer{
			ReadTransformer: readtransformer.GeminiCSV{
				Reader: csvreader.HTTPCSVReader{
					Location: "https://www.cryptodatadownload.com/cdd/Gemini_BATUSD_1h.csv",
					SkipRows: 2,
				},
			},
		},
	},
	sharedtypes.PairUSDLINK: {
		syncer.HolisticSyncer{
			ReadTransformer: readtransformer.GeminiCSV{
				Reader: csvreader.HTTPCSVReader{
					Location: "https://www.cryptodatadownload.com/cdd/Gemini_LINKUSD_d.csv",
					SkipRows: 2,
				},
			},
		},
		syncer.HolisticSyncer{
			ReadTransformer: readtransformer.GeminiCSV{
				Reader: csvreader.HTTPCSVReader{
					Location: "https://www.cryptodatadownload.com/cdd/Gemini_LINKUSD_1h.csv",
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
				log.Println("Failed Sync for " + p.String() + ", failing over...")
				log.Println(err.Error())
			} else {
				// Successful sync, no need to fail over.
				successful = true
				break
			}
		}
		if !successful {
			failedSyncs++
		} else {
			log.Println("Synced:" + p.String())
		}
	}

	if failedSyncs > 0 {
		return errors.New("one or more pairs failed to sync")
	}

	return nil
}
