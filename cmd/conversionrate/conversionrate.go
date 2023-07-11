package conversionrate

import (
	"errors"
	"log"
)

func MarketValue(currencyFrom, currencyTo string, timestamp int) {

}

func SyncAll() error {
	failedSyncs := 0
	for p, syncers := range PairSyncers {
		successful := false
		for _, s := range syncers {
			err := s.sync()
			if err != nil {
				log.Println("Failed Sync for " + p.currency1 + p.currency2 + ", failing over...")
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
