package sources

// TODO(add const configuration here)

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/transactions/transactions"
	"math"
	"strconv"
	"strings"
	"time"
)

type LunoSource struct{}

var currencyMap = map[string]string{
	"XBT": "BTC",
}

func (s LunoSource) TransformRow(row []string) ([]transactions.Transaction, error) {
	amount, err := strconv.ParseFloat(row[5], 64)
	if err != nil {
		return nil, err
	}

	tim, err := time.Parse("2006-01-02 15:04:05", row[2])
	if err != nil {
		return nil, err
	}

	// We need to be careful with this, it could cause mis-matches with the tax fiat selected.
	fiatValue, err := strconv.ParseFloat(row[12], 64)
	if err != nil {
		return nil, err
	}

	wholePrice := fiatValue / math.Abs(amount)

	hash := md5.Sum([]byte(strings.Join(row[:], ",")))
	hashString := hex.EncodeToString(hash[:])

	return []transactions.Transaction{{
		UID:          hashString,
		Transformer:  transactions.TransformTypeLuno,
		Currency:     mapCurrency(row[4]),
		DetectedType: s.inferType(row, amount),
		Amount:       math.Abs(amount),
		Timestamp:    tim.Unix(),
		WholePriceAtPoint: transactions.FiatPrice{
			Fiat:  row[11],
			Price: wholePrice,
		},
	}}, nil
}

func (s LunoSource) inferType(row []string, amount float64) transactions.TransactionType {
	if amount < 0 {
		if strings.Contains(row[3], "Sold") {
			return transactions.TypeSell
		}

		if strings.Contains(row[3], "fee") {
			return transactions.TypeFee
		}

		// Currently no way to infer external sends, safer to assume internal and let user override.
		return transactions.TypeSendInternal
	}

	if strings.Contains(row[3], "Bought") {
		return transactions.TypeBuy
	}

	// Currently no way to infer external receives, safer to assume internal and let user override.
	return transactions.TypeReceiveInternal
}

func mapCurrency(s string) string {
	cur, ok := currencyMap[s]
	if ok {
		return cur
	}

	return s
}
