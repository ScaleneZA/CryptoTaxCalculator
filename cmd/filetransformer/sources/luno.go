package sources

// TODO(add const configuration here)

import (
	"crypto/md5"
	"encoding/hex"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/sharedtypes"
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

	hash := md5.Sum([]byte(strings.Join(row[:], ",")))
	hashString := hex.EncodeToString(hash[:])

	return sharedtypes.Transaction{
		UID:               hashString,
		Transformer:       sharedtypes.TransformTypeLuno,
		Currency:          mapCurrency(row[4]),
		DetectedType:      inferType(row, amount),
		Amount:            math.Abs(amount),
		Timestamp:         tim.Unix(),
		WholePriceAtPoint: wholePrice,
	}, nil
}

func inferType(row []string, amount float64) sharedtypes.TransactionType {
	if amount < 0 {
		if strings.Contains(row[3], "Sold") {
			return sharedtypes.TypeSell
		}

		if strings.Contains(row[3], "fee") {
			return sharedtypes.TypeFee
		}

		// Currently no way to infer external sends, safer to assume internal and let user override.
		return sharedtypes.TypeSendInternal
	}

	if strings.Contains(row[3], "Bought") {
		return sharedtypes.TypeBuy
	}

	// Currently no way to infer external receives, safer to assume internal and let user override.
	return sharedtypes.TypeReceiveInternal
}

func mapCurrency(s string) string {
	cur, ok := currencyMap[s]
	if ok {
		return cur
	}

	return s
}
