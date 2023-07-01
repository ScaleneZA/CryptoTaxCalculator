package sources

import (
	"math"
	"strconv"
	"time"

	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/taxcalculator/sharedtypes"
)

type LunoSource struct{}

var currencyMap = map[string]string{
	"XBT": "BTC",
}

func (s LunoSource) TransformRow(row []string) (sharedtypes.Transaction, error) {
	amount, err := strconv.ParseFloat(row[5], 64)
	if err != nil {
		return sharedtypes.Transaction{}, err
	}

	tim, err := time.Parse("2006-01-02 15:04:05", row[2])
	if err != nil {
		return sharedtypes.Transaction{}, err
	}

	typ := sharedtypes.TypeBuy
	if amount < 0 {
		amount = math.Abs(amount)
		typ = sharedtypes.TypeSell
	}

	fiatValue, err := strconv.ParseFloat(row[12], 64)
	if err != nil {
		return sharedtypes.Transaction{}, err
	}

	wholePrice := fiatValue / amount

	return sharedtypes.Transaction{
		Currency:          mapCurrency(row[4]),
		Typ:               typ,
		Amount:            amount,
		Timestamp:         tim.Unix(),
		WholePriceAtPoint: wholePrice,
	}, nil
}

func mapCurrency(s string) string {
	cur, ok := currencyMap[s]
	if ok {
		return cur
	}

	return s
}
