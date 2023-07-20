package sync

import (
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/rates/conversionrate"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/rates/conversionrate/ops/sync/readtransformer"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/rates/conversionrate/ops/sync/readtransformer/csvreader"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/rates/conversionrate/ops/sync/syncer"
	"github.com/luno/jettison/errors"
	"github.com/luno/jettison/j"
	"log"
)

var PairSyncers = map[conversionrate.Pair][]syncer.Syncer{
	// TODO: Find a better way to seed this data. Perhaps an XML reader for this remote resource:
	// https://custom.resbank.co.za/SarbWebApi/WebIndicators/Shared/GetTimeseriesObservations//EXCX135D/2015-01-01/2023-07-20
	conversionrate.PairZARUSD: {
		syncer.HolisticSyncer{
			ReadTransformer: readtransformer.SARBCSV{
				Reader: csvreader.LocalCSVReader{
					Location: "cmd/rates/data/USD_ZAR_seed.csv",
					SkipRows: 4,
				},
				FromCurrency: conversionrate.PairZARUSD.FromCurrency,
				ToCurrency:   conversionrate.PairZARUSD.ToCurrency,
			},
		},
	},
	conversionrate.PairUSDBTC: {
		syncer.HolisticSyncer{
			ReadTransformer: readtransformer.GeminiCSV{
				Reader: csvreader.HTTPCSVReader{
					Location: "https://www.cryptodatadownload.com/cdd/Gemini_BTCUSD_1h.csv",
					SkipRows: 2,
				},
			},
		},
	},
	conversionrate.PairUSDETH: {
		syncer.HolisticSyncer{
			ReadTransformer: readtransformer.GeminiCSV{
				Reader: csvreader.HTTPCSVReader{
					Location: "https://www.cryptodatadownload.com/cdd/Gemini_ETHUSD_1h.csv",
					SkipRows: 2,
				},
			},
		},
	},
	conversionrate.PairUSDLTC: {
		syncer.HolisticSyncer{
			ReadTransformer: readtransformer.GeminiCSV{
				Reader: csvreader.HTTPCSVReader{
					Location: "https://www.cryptodatadownload.com/cdd/Gemini_LTCUSD_1h.csv",
					SkipRows: 2,
				},
			},
		},
	},
	conversionrate.PairUSDBCH: {
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
	conversionrate.PairUSDBAT: {
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
	conversionrate.PairUSDLINK: {
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
	var failedSyncs []conversionrate.Pair
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
		return errors.Wrap(conversionrate.ErrPairSyncFailed, "", j.MKV{
			"failed_syncs": failedSyncs,
		})
	}

	return nil
}
