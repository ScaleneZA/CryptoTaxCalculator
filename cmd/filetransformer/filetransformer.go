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

func Transform(files []io.Reader, typ TransformType) ([]sharedtypes.Transaction, error) {
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

	sort.Slice(ts, func(i, j int) bool {
		return ts[i].Timestamp < ts[j].Timestamp
	})

	return ts, nil
}

func sourceFromType(typ TransformType) (sources.Source, error) {
	var src sources.Source
	switch typ {
	case TransformTypeBasic:
		src = sources.BasicSource{}
	case TransformTypeLuno:
		src = sources.LunoSource{}
	default:
		return nil, errors.New("invalid source")
	}

	return src, nil
}
