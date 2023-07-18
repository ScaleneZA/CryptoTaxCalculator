package sources

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/taxcalculator"
	"math"
	"strconv"
	"strings"
)

type BasicSource struct{}

func (s BasicSource) TransformRow(row []string) (taxcalculator.Transaction, error) {
	amount, err := strconv.ParseFloat(row[2], 64)
	if err != nil {
		return taxcalculator.Transaction{}, err
	}

	timestamp, err := strconv.Atoi(row[3])
	if err != nil {
		return taxcalculator.Transaction{}, err
	}

	typ := taxcalculator.TypeBuy
	if amount < 0 {
		amount = math.Abs(amount)
		typ = taxcalculator.TypeSell
	}

	wholePrice, err := strconv.ParseFloat(row[4], 64)
	if err != nil {
		return taxcalculator.Transaction{}, err
	}

	hash := md5.Sum([]byte(strings.Join(row[:], ",")))
	hashString := hex.EncodeToString(hash[:])

	return taxcalculator.Transaction{
		UID:               hashString,
		Transformer:       taxcalculator.TransformTypeBasic,
		Currency:          row[1],
		DetectedType:      typ,
		Amount:            amount,
		Timestamp:         int64(timestamp),
		WholePriceAtPoint: wholePrice,
	}, nil
}
