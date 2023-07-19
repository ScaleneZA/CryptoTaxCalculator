package filetransformer

import (
	"encoding/csv"
	"fmt"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/transactions/transactions"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/transactions/transactions/ops/filetransformer/sources"
	"github.com/luno/jettison/errors"
	"github.com/luno/jettison/j"
	"io"
	"sort"
)

func TransformAll(typeFiles map[transactions.TransformType][]io.Reader) ([]transactions.Transaction, error) {
	var ts []transactions.Transaction

	for typ, files := range typeFiles {
		tts, err := Transform(files, typ)
		if err != nil {
			return nil, err
		}

		ts = append(ts, tts...)
	}

	return sortTransactions(ts), nil
}

func Transform(files []io.Reader, typ transactions.TransformType) ([]transactions.Transaction, error) {
	var ts []transactions.Transaction
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

func sortTransactions(ts []transactions.Transaction) []transactions.Transaction {
	sort.Slice(ts, func(i, j int) bool {
		if ts[i].Timestamp == ts[j].Timestamp {
			return ts[i].DetectedType < ts[j].DetectedType
		}
		return ts[i].Timestamp < ts[j].Timestamp
	})

	return ts
}

func sourceFromType(typ transactions.TransformType) (sources.Source, error) {
	var src sources.Source
	switch typ {
	case transactions.TransformTypeBasic:
		src = sources.BasicSource{}
	case transactions.TransformTypeLuno:
		src = sources.LunoSource{}
	default:
		return nil, errors.Wrap(transactions.ErrUnsupportedTranformType, "", j.MKV{
			"type": typ,
		})
	}

	return src, nil
}
