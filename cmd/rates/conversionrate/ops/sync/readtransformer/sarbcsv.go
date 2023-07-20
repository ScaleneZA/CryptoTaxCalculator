package readtransformer

import (
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/rates/conversionrate"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/rates/conversionrate/ops/sync/readtransformer/csvreader"
	"strconv"
	"strings"
	"time"
)

type SARBCSV struct {
	Reader       csvreader.Reader
	FromCurrency string
	ToCurrency   string
}

func (t SARBCSV) ReadAndTransform() ([]conversionrate.MarketPair, error) {
	rows, err := t.Reader.Read()
	if err != nil {
		return nil, err
	}

	var mps []conversionrate.MarketPair
	for _, r := range rows {
		mp, err := t.transformRow(r)
		if err != nil {
			return nil, err
		}

		mps = append(mps, mp)
	}

	return mps, nil
}

func (t SARBCSV) transformRow(row []string) (conversionrate.MarketPair, error) {
	date := strings.TrimSpace(row[0])

	tim, err := time.Parse(time.DateOnly, date)
	if err != nil {
		return conversionrate.MarketPair{}, err
	}

	price, err := strconv.ParseFloat(row[1], 64)
	if err != nil {
		return conversionrate.MarketPair{}, err
	}

	return conversionrate.MarketPair{
		Pair: conversionrate.Pair{
			FromCurrency: t.FromCurrency,
			ToCurrency:   t.ToCurrency,
		},
		MarketSlice: conversionrate.MarketSlice{
			Timestamp: tim.Unix(),
			Open:      price,
			High:      price,
			Low:       price,
			Close:     price,
		},
	}, nil
}
