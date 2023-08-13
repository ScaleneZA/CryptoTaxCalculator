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

type KrakenSource struct{}

func (s KrakenSource) TransformRow(row []string) ([]transactions.Transaction, error) {
	amount, err := strconv.ParseFloat(row[7], 64)
	if err != nil {
		return nil, err
	}

	tim, err := time.Parse("2006-01-02 15:04:05", row[2])
	if err != nil {
		return nil, err
	}

	hash := md5.Sum([]byte(strings.Join(row[:], ",")))
	hashString := hex.EncodeToString(hash[:])

	ts := []transactions.Transaction{{
		UID:          hashString,
		Transformer:  transactions.TransformTypeKraken,
		Currency:     mapCurrency(row[6]),
		DetectedType: s.inferType(row, amount),
		Amount:       math.Abs(amount),
		Timestamp:    tim.Unix(),
	}}

	feeAmnt, err := strconv.ParseFloat(row[8], 64)
	if err != nil {
		// NoReturnErr: Default to 0.
		feeAmnt = 0
	}
	if feeAmnt > 0 {
		hash = md5.Sum([]byte(strings.Join(row[:], ",") + "-fee"))
		hashString = hex.EncodeToString(hash[:])

		ts = append(ts, transactions.Transaction{
			UID:          hashString,
			Transformer:  transactions.TransformTypeKraken,
			Currency:     mapCurrency(row[6]),
			DetectedType: transactions.TypeFee,
			Amount:       math.Abs(feeAmnt),
			Timestamp:    tim.Unix(),
		})
	}

	return ts, nil
}

func (s KrakenSource) inferType(row []string, amount float64) transactions.TransactionType {
	if amount < 0 {
		if strings.Contains(row[3], "trade") {
			return transactions.TypeSell
		}

		// Currently no way to infer external sends, safer to assume internal and let user override.
		return transactions.TypeSendInternal
	}

	if strings.Contains(row[3], "trade") {
		return transactions.TypeBuy
	}

	if strings.Contains(row[3], "transfer") {
		return transactions.TypeAirdrop
	}

	// Currently no way to infer external receives, safer to assume internal and let user override.
	return transactions.TypeReceiveInternal
}
