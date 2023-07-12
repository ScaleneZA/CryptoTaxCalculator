package sync

import (
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/conversionrate/sync/readtransformer"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/conversionrate/sync/writer"
)

// syncer syncs an entire batch at once
type syncer struct {
	readTransformer readtransformer.ReadTransformer
	writer          writer.Writer
}

func (s syncer) sync() error {
	mps, err := s.readTransformer.ReadAndTransform()

	err = s.writer.Write(mps)
	if err != nil {
		return err
	}

	return nil
}
