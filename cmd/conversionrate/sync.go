package conversionrate

import (
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/conversionrate/reader"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/conversionrate/writer"
)

type syncer struct {
	reader reader.Reader
	writer writer.Writer
}

func (s syncer) sync() error {
	r, err := s.reader.Read()
	if err != nil {
		return err
	}
	defer r.Close()

	err = s.writer.Write(r)
	if err != nil {
		return err
	}

	return nil
}
