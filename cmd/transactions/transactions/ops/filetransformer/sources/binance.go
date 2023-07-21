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

type BinanceSource struct{}

func (s BinanceSource) TransformRow(row []string) ([]transactions.Transaction, error) {
	amount, err := strconv.ParseFloat(row[5], 64)
	if err != nil {
		return nil, err
	}

	tim, err := time.Parse("2006-01-02 15:04:05", row[1])
	if err != nil {
		return nil, err
	}

	hash := md5.Sum([]byte(strings.Join(row[:], ",")))
	hashString := hex.EncodeToString(hash[:])

	return []transactions.Transaction{{
		UID:          hashString,
		Transformer:  transactions.TransformTypeBinance,
		Currency:     row[4],
		DetectedType: s.inferType(row, amount),
		Amount:       math.Abs(amount),
		Timestamp:    tim.Unix(),
	}}, nil
}

func (s BinanceSource) inferType(row []string, amount float64) transactions.TransactionType {
	if amount < 0 {
		if strings.Contains(row[3], "trading") {
			return transactions.TypeSell
		}

		// Currently no way to infer external sends, safer to assume internal and let user override.
		return transactions.TypeSendInternal
	}

	if strings.Contains(row[3], "trading") {
		return transactions.TypeBuy
	}

	if strings.Contains(row[3], "interest") {
		return transactions.TypeInterest
	}

	// Currently no way to infer external receives, safer to assume internal and let user override.
	return transactions.TypeReceiveInternal
}
