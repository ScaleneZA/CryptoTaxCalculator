package sources

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/transactions/transactions"
	"math"
	"strconv"
	"strings"
)

type BasicSource struct{}

func (s BasicSource) TransformRow(row []string) ([]transactions.Transaction, error) {
	amount, err := strconv.ParseFloat(row[2], 64)
	if err != nil {
		return nil, err
	}

	timestamp, err := strconv.Atoi(row[3])
	if err != nil {
		return nil, err
	}

	typ := transactions.TypeBuy
	if amount < 0 {
		amount = math.Abs(amount)
		typ = transactions.TypeSell
	}

	wholePrice, err := strconv.ParseFloat(row[4], 64)
	if err != nil {
		return nil, err
	}

	hash := md5.Sum([]byte(strings.Join(row[:], ",")))
	hashString := hex.EncodeToString(hash[:])

	return []transactions.Transaction{{
		UID:          hashString,
		Transformer:  transactions.TransformTypeBasic,
		Currency:     row[1],
		DetectedType: typ,
		Amount:       amount,
		Timestamp:    int64(timestamp),
		WholePriceAtPoint: transactions.FiatPrice{
			Fiat:  row[5],
			Price: wholePrice,
		},
	}}, nil
}
