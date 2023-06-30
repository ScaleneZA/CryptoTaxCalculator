package transformer

import (
	"errors"
	"fmt"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/taxcalculator/sharedtypes"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/taxcalculator/transformer/sources"
	"github.com/xuri/excelize/v2"
	"sort"
)

func Transform(filename string, typ TransformType) ([]sharedtypes.Transaction, error) {
	rows, err := readFile(filename)
	if err != nil {
		return nil, err
	}

	var src sources.Source
	switch typ {
	case TransformTypeTest:
		src = sources.TestSource{}
	default:
		return nil, errors.New("invalid source")
	}

	var ts []sharedtypes.Transaction
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

	sort.Slice(ts, func(i, j int) bool {
		return ts[i].Timestamp < ts[j].Timestamp
	})

	return ts, nil
}

// TODO: Pull this out to its own package perhaps?
func readFile(filename string) ([][]string, error) {
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return nil, err
	}

	return rows, nil
}
