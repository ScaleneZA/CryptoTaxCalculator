package syncer

import (
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/conversionrate/db/markets"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/conversionrate/ops/sync/readtransformer"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/db"
)

// HolisticSyncer deletes all data and syncs from the beginning of time.
type HolisticSyncer struct {
	ReadTransformer readtransformer.ReadTransformer
}

func (s HolisticSyncer) Sync(b Backends) error {
	mps, err := s.ReadTransformer.ReadAndTransform()
	if err != nil {
		return err
	}

	// TODO: Don't do this.
	db.ResetDB(b.DB())

	for _, mp := range mps {
		_, err := markets.Create(b.DB(), mp)
		if err != nil {
			return err
		}
	}

	return nil
}
