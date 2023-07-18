package filetransformer

import (
	"encoding/csv"
	"fmt"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/taxcalculator"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/taxcalculator/ops/filetransformer/sources"
	"github.com/luno/jettison/errors"
	"github.com/luno/jettison/j"
	"io"
	"sort"
)

func TransformAll(typeFiles map[taxcalculator.TransformType][]io.Reader) ([]taxcalculator.Transaction, error) {
	var ts []taxcalculator.Transaction

	for typ, files := range typeFiles {
		tts, err := Transform(files, typ)
		if err != nil {
			return nil, err
		}

		ts = append(ts, tts...)
	}

	return sortTransactions(ts), nil
}

func Transform(files []io.Reader, typ taxcalculator.TransformType) ([]taxcalculator.Transaction, error) {
	var ts []taxcalculator.Transaction
	for _, file := range files {
		reader := csv.NewReader(file)
		rows, err := reader.ReadAll()
		if err != nil {
			return nil, err
		}

		src, err := sourceFromType(typ)
		if err != nil {
			return nil, err
		}

		headerCount := 0
		for i, r := range rows {
			t, err := src.TransformRow(r)
			if err != nil {
				if headerCount < 1 {
					fmt.Println(fmt.Sprintf("Skipping row %d may be header", i))
					headerCount++
					continue
				}
				return nil, err
			}

			ts = append(ts, t)
		}

	}

	return sortTransactions(ts), nil
}

func sortTransactions(ts []taxcalculator.Transaction) []taxcalculator.Transaction {
	sort.Slice(ts, func(i, j int) bool {
		if ts[i].Timestamp == ts[j].Timestamp {
			return ts[i].DetectedType < ts[j].DetectedType
		}
		return ts[i].Timestamp < ts[j].Timestamp
	})

	return ts
}

func sourceFromType(typ taxcalculator.TransformType) (sources.Source, error) {
	var src sources.Source
	switch typ {
	case taxcalculator.TransformTypeBasic:
		src = sources.BasicSource{}
	case taxcalculator.TransformTypeLuno:
		src = sources.LunoSource{}
	default:
		return nil, errors.Wrap(taxcalculator.ErrUnsupportedTranformType, "", j.MKV{
			"type": typ,
		})
	}

	return src, nil
}
