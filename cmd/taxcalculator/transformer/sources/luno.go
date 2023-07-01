package sources

// TODO(add const configuration here)

import (
	"math"
	"strconv"
	"strings"
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

	fiatValue, err := strconv.ParseFloat(row[12], 64)
	if err != nil {
		return sharedtypes.Transaction{}, err
	}

	wholePrice := fiatValue / math.Abs(amount)

	return sharedtypes.Transaction{
		Currency:          mapCurrency(row[4]),
		Typ:               inferType(row, amount, nil),
		Amount:            math.Abs(amount),
		Timestamp:         tim.Unix(),
		WholePriceAtPoint: wholePrice,
	}, nil
}

func inferType(row []string, amount float64, internalAddresses map[string]bool) sharedtypes.TransactionType {
	if amount < 0 {
		if internalAddresses[row[8]] {
			return sharedtypes.TypeSendInternal
		}
		if strings.Contains(row[3], "Sold") {
			return sharedtypes.TypeSell
		}

		if strings.Contains(row[3], "fee") {
			return sharedtypes.TypeFee
		}

		return sharedtypes.TypeSendExternal
	}

	if strings.Contains(row[3], "Bought") {
		return sharedtypes.TypeBuy
	}

	// Currently no way to infer external receives, assume all internal for now.
	return sharedtypes.TypeReceiveInternal
}

func mapCurrency(s string) string {
	cur, ok := currencyMap[s]
	if ok {
		return cur
	}

	return s
}
