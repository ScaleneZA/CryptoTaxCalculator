package sync

import (
	"errors"
	"log"
)

func SyncAll(b Backends) error {
	failedSyncs := 0
	for p, syncers := range PairSyncers {
		successful := false
		for _, s := range syncers {
			err := s.sync(b)
			if err != nil {
				log.Println("Failed Sync for " + p.Currency1 + p.Currency2 + ", failing over...")
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
