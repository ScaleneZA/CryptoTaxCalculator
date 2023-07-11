package sync

import (
	"errors"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/conversionrate/sync/reader"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/conversionrate/sync/transformer"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/conversionrate/sync/writer"
	"log"
)

type syncer struct {
	reader      reader.Reader
	transformer transformer.Transformer
	writer      writer.Writer
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

func (s syncer) sync() error {
	r, err := s.reader.Read()
	if err != nil {
		return err
	}
	defer r.Close()

	// TODO(Transform)

	err = s.writer.Write(r)
	if err != nil {
		return err
	}

	return nil
}
