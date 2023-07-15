package syncer

import (
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/conversionrate/db/markets"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/conversionrate/ops/sync/readtransformer"
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

	err = markets.Truncate(b.DB())
	if err != nil {
		return err
	}

	for _, mp := range mps {
		_, err := markets.Create(b.DB(), mp)
		if err != nil {
			return err
		}
	}

	return nil
}
