package readtransformer

import (
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/conversionrate/sync/readtransformer/csvreader"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/conversionrate/sync/sharedtypes"
	"strconv"
)

type GeminiCSV struct {
	Reader csvreader.Reader
}

func (t GeminiCSV) ReadAndTransform() ([]sharedtypes.MarketSlice, error) {
	rows, err := t.Reader.Read()
	if err != nil {
		return nil, err
	}
	var mps []sharedtypes.MarketSlice
	for _, r := range rows {
		mp, err := t.transformRow(r)
		if err != nil {
			return nil, err
		}

		mps = append(mps, mp)
	}

	return mps, nil
}

func (t GeminiCSV) transformRow(row []string) (sharedtypes.MarketSlice, error) {
	tim, err := strconv.Atoi(row[0])
	if err != nil {
		return sharedtypes.MarketSlice{}, err
	}

	open, err := strconv.ParseFloat(row[3], 64)
	if err != nil {
		return sharedtypes.MarketSlice{}, err
	}
	high, err := strconv.ParseFloat(row[4], 64)
	if err != nil {
		return sharedtypes.MarketSlice{}, err
	}
	low, err := strconv.ParseFloat(row[5], 64)
	if err != nil {
		return sharedtypes.MarketSlice{}, err
	}
	clos, err := strconv.ParseFloat(row[6], 64)
	if err != nil {
		return sharedtypes.MarketSlice{}, err
	}

	return sharedtypes.MarketSlice{
		Timestamp: tim,
		Open:      open,
		High:      high,
		Low:       low,
		Close:     clos,
	}, nil
}
