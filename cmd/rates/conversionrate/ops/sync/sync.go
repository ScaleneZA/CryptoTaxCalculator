package sync

import (
	conversionrate2 "github.com/ScaleneZA/CryptoTaxCalculator/cmd/rates/conversionrate"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/rates/conversionrate/ops/sync/readtransformer"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/rates/conversionrate/ops/sync/readtransformer/csvreader"
	syncer2 "github.com/ScaleneZA/CryptoTaxCalculator/cmd/rates/conversionrate/ops/sync/syncer"
	"github.com/luno/jettison/errors"
	"github.com/luno/jettison/j"
	"log"
)

var PairSyncers = map[conversionrate2.Pair][]syncer2.Syncer{
	conversionrate2.PairUSDBTC: {
		syncer2.HolisticSyncer{
			ReadTransformer: readtransformer.GeminiCSV{
				Reader: csvreader.HTTPCSVReader{
					Location: "https://www.cryptodatadownload.com/cdd/Gemini_BTCUSD_1h.csv",
					SkipRows: 2,
				},
			},
		},
	},
	conversionrate2.PairUSDETH: {
		syncer2.HolisticSyncer{
			ReadTransformer: readtransformer.GeminiCSV{
				Reader: csvreader.HTTPCSVReader{
					Location: "https://www.cryptodatadownload.com/cdd/Gemini_ETHUSD_1h.csv",
					SkipRows: 2,
				},
			},
		},
	},
	conversionrate2.PairUSDLTC: {
		syncer2.HolisticSyncer{
			ReadTransformer: readtransformer.GeminiCSV{
				Reader: csvreader.HTTPCSVReader{
					Location: "https://www.cryptodatadownload.com/cdd/Gemini_LTCUSD_1h.csv",
					SkipRows: 2,
				},
			},
		},
	},
	conversionrate2.PairUSDBCH: {
		syncer2.HolisticSyncer{
			ReadTransformer: readtransformer.GeminiCSV{
				Reader: csvreader.HTTPCSVReader{
					Location: "https://www.cryptodatadownload.com/cdd/Gemini_BCHUSD_d.csv",
					SkipRows: 2,
				},
			},
		},
		syncer2.HolisticSyncer{
			ReadTransformer: readtransformer.GeminiCSV{
				Reader: csvreader.HTTPCSVReader{
					Location: "https://www.cryptodatadownload.com/cdd/Gemini_BCHUSD_1h.csv",
					SkipRows: 2,
				},
			},
		},
	},
	conversionrate2.PairUSDBAT: {
		syncer2.HolisticSyncer{
			ReadTransformer: readtransformer.GeminiCSV{
				Reader: csvreader.HTTPCSVReader{
					Location: "https://www.cryptodatadownload.com/cdd/Gemini_BATUSD_d.csv",
					SkipRows: 2,
				},
			},
		},
		syncer2.HolisticSyncer{
			ReadTransformer: readtransformer.GeminiCSV{
				Reader: csvreader.HTTPCSVReader{
					Location: "https://www.cryptodatadownload.com/cdd/Gemini_BATUSD_1h.csv",
					SkipRows: 2,
				},
			},
		},
	},
	conversionrate2.PairUSDLINK: {
		syncer2.HolisticSyncer{
			ReadTransformer: readtransformer.GeminiCSV{
				Reader: csvreader.HTTPCSVReader{
					Location: "https://www.cryptodatadownload.com/cdd/Gemini_LINKUSD_d.csv",
					SkipRows: 2,
				},
			},
		},
		syncer2.HolisticSyncer{
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
	var failedSyncs []conversionrate2.Pair
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
			failedSyncs = append(failedSyncs, p)
		} else {
			log.Println("Synced:" + p.String())
		}
	}

	if len(failedSyncs) > 0 {
		return errors.Wrap(conversionrate2.ErrPairSyncFailed, "", j.MKV{
			"failed_syncs": failedSyncs,
		})
	}

	return nil
}
