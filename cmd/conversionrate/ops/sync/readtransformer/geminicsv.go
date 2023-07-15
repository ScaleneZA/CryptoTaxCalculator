package readtransformer

import (
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/conversionrate/ops/sync/readtransformer/csvreader"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/conversionrate/sharedtypes"
	"strconv"
	"strings"
)

type GeminiCSV struct {
	Reader csvreader.Reader
}

func (t GeminiCSV) ReadAndTransform() ([]sharedtypes.MarketPair, error) {
	rows, err := t.Reader.Read()
	if err != nil {
		return nil, err
	}

	var mps []sharedtypes.MarketPair
	for _, r := range rows {
		mp, err := t.transformRow(r)
		if err != nil {
			return nil, err
		}

		mps = append(mps, mp)
	}

	return mps, nil
}

func (t GeminiCSV) transformRow(row []string) (sharedtypes.MarketPair, error) {
	tim, err := strconv.Atoi(row[0])
	if err != nil {
		return sharedtypes.MarketPair{}, err
	}

	open, err := strconv.ParseFloat(row[3], 64)
	if err != nil {
		return sharedtypes.MarketPair{}, err
	}
	high, err := strconv.ParseFloat(row[4], 64)
	if err != nil {
		return sharedtypes.MarketPair{}, err
	}
	low, err := strconv.ParseFloat(row[5], 64)
	if err != nil {
		return sharedtypes.MarketPair{}, err
	}
	clos, err := strconv.ParseFloat(row[6], 64)
	if err != nil {
		return sharedtypes.MarketPair{}, err
	}

	currencyParts := strings.Split(row[2], "/")

	return sharedtypes.MarketPair{
		Pair: sharedtypes.Pair{
			FromCurrency: currencyParts[1],
			ToCurrency:   currencyParts[0],
		},
		MarketSlice: sharedtypes.MarketSlice{
			Timestamp: int64(tim),
			Open:      open,
			High:      high,
			Low:       low,
			Close:     clos,
		},
	}, nil
}
