package filetransformer

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"sort"

	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/filetransformer/sources"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/sharedtypes"
)

func TransformAll(typeFiles map[sharedtypes.TransformType][]io.Reader) ([]sharedtypes.Transaction, error) {
	var ts []sharedtypes.Transaction

	for typ, files := range typeFiles {
		tts, err := Transform(files, typ)
		if err != nil {
			return nil, err
		}

		ts = append(ts, tts...)
	}

	return sortTransactions(ts), nil
}

func Transform(files []io.Reader, typ sharedtypes.TransformType) ([]sharedtypes.Transaction, error) {
	var ts []sharedtypes.Transaction
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

func sortTransactions(ts []sharedtypes.Transaction) []sharedtypes.Transaction {
	sort.Slice(ts, func(i, j int) bool {
		if ts[i].Timestamp == ts[j].Timestamp {
			return ts[i].DetectedType < ts[j].DetectedType
		}
		return ts[i].Timestamp < ts[j].Timestamp
	})

	return ts
}

func sourceFromType(typ sharedtypes.TransformType) (sources.Source, error) {
	var src sources.Source
	switch typ {
	case sharedtypes.TransformTypeBasic:
		src = sources.BasicSource{}
	case sharedtypes.TransformTypeLuno:
		src = sources.LunoSource{}
	default:
		return nil, errors.New("invalid source")
	}

	return src, nil
}
