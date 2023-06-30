package sources

import (
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/taxcalculator/sharedtypes"
	"math"
	"strconv"
)

type TestSource struct{}

func (s TestSource) TransformRow(row []string) (sharedtypes.Transaction, error) {
	amount, err := strconv.ParseFloat(row[2], 64)
	if err != nil {
		return sharedtypes.Transaction{}, err
	}

	timestamp, err := strconv.Atoi(row[3])
	if err != nil {
		return sharedtypes.Transaction{}, err
	}

	typ := sharedtypes.TypeBuy
	if amount < 0 {
		amount = math.Abs(amount)
		typ = sharedtypes.TypeSell
	}

	wholePrice, err := strconv.ParseFloat(row[4], 64)
	if err != nil {
		return sharedtypes.Transaction{}, err
	}

	return sharedtypes.Transaction{
		Typ:               typ,
		Amount:            amount,
		Timestamp:         int64(timestamp),
		WholePriceAtPoint: wholePrice,
	}, nil
}
