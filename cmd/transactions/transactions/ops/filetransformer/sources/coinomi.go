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

type CoinomiSource struct{}

func (s CoinomiSource) TransformRow(row []string) ([]transactions.Transaction, error) {
	// Amount seems to include fees
	amount, err := strconv.ParseFloat(row[4], 64)
	if err != nil {
		return nil, err
	}

	tim, err := time.Parse("2006-01-02T15:04Z", row[10])
	if err != nil {
		return nil, err
	}
	ts := []transactions.Transaction{{
		UID:          row[8],
		Transformer:  transactions.TransformTypeCoinomi,
		Currency:     row[5],
		DetectedType: s.inferType(row, amount),
		Amount:       math.Abs(amount),
		Timestamp:    tim.Unix(),
	}}

	feeAmnt, err := strconv.ParseFloat(row[6], 64)
	if err != nil {
		feeAmnt = 0
	}
	if feeAmnt > 0 {
		hash := md5.Sum([]byte(strings.Join(row[:], ",")))
		hashString := hex.EncodeToString(hash[:])

		ts = append(ts, transactions.Transaction{
			UID:          hashString,
			Transformer:  transactions.TransformTypeCoinomi,
			Currency:     row[5],
			DetectedType: transactions.TypeFee,
			Amount:       math.Abs(feeAmnt),
			Timestamp:    tim.Unix(),
		})
	}

	return ts, nil
}

func (s CoinomiSource) inferType(row []string, amount float64) transactions.TransactionType {
	if amount < 0 {
		// Currently no way to infer external sends, safer to assume internal and let user override.
		return transactions.TypeSendInternal
	}

	// Currently no way to infer external receives, safer to assume internal and let user override.
	return transactions.TypeReceiveInternal
}
