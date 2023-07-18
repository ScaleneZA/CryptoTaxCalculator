package sources

// TODO(add const configuration here)

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/taxcalculator"
	"math"
	"strconv"
	"strings"
	"time"
)

type LunoSource struct{}

var currencyMap = map[string]string{
	"XBT": "BTC",
}

func (s LunoSource) TransformRow(row []string) (taxcalculator.Transaction, error) {
	amount, err := strconv.ParseFloat(row[5], 64)
	if err != nil {
		return taxcalculator.Transaction{}, err
	}

	tim, err := time.Parse("2006-01-02 15:04:05", row[2])
	if err != nil {
		return taxcalculator.Transaction{}, err
	}

	fiatValue, err := strconv.ParseFloat(row[12], 64)
	if err != nil {
		return taxcalculator.Transaction{}, err
	}

	wholePrice := fiatValue / math.Abs(amount)

	hash := md5.Sum([]byte(strings.Join(row[:], ",")))
	hashString := hex.EncodeToString(hash[:])

	return taxcalculator.Transaction{
		UID:               hashString,
		Transformer:       taxcalculator.TransformTypeLuno,
		Currency:          mapCurrency(row[4]),
		DetectedType:      inferType(row, amount),
		Amount:            math.Abs(amount),
		Timestamp:         tim.Unix(),
		WholePriceAtPoint: wholePrice,
	}, nil
}

func inferType(row []string, amount float64) taxcalculator.TransactionType {
	if amount < 0 {
		if strings.Contains(row[3], "Sold") {
			return taxcalculator.TypeSell
		}

		if strings.Contains(row[3], "fee") {
			return taxcalculator.TypeFee
		}

		// Currently no way to infer external sends, safer to assume internal and let user override.
		return taxcalculator.TypeSendInternal
	}

	if strings.Contains(row[3], "Bought") {
		return taxcalculator.TypeBuy
	}

	// Currently no way to infer external receives, safer to assume internal and let user override.
	return taxcalculator.TypeReceiveInternal
}

func mapCurrency(s string) string {
	cur, ok := currencyMap[s]
	if ok {
		return cur
	}

	return s
}
