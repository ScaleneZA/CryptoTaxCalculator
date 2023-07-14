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

func (s syncer) sync(b Backends) error {
	mps, err := s.readTransformer.ReadAndTransform()

	err = s.writer.DeleteAll(b)
	if err != nil {
		return err
	}
	err = s.writer.WriteAll(b, mps)
	if err != nil {
		return err
	}

	return nil
}
