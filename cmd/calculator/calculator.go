package calculator

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/sharedtypes"
)

/*
Next Steps:
* Multi-file
* Way to initialize at a point in time
* Fetch whole price if not supplied
* Separate Sends, Receives, Buys, Sells - Override detected types
*/

func Calculate(transactions []sharedtypes.Transaction) (YearEndTotals, error) {
	yearEndTotals := make(YearEndTotals)

	for currency := range uniqueCurrencies(transactions) {
		var tally []sharedtypes.Transaction
		var lastTimestamp int64

		for _, t := range transactions {
			if t.Currency != currency {
				continue
			}

			if t.Timestamp < lastTimestamp {
				return nil, errors.New("Transaction Slice not ordered correctly")
			}
			lastTimestamp = t.Timestamp

			if t.FinalType().ShouldIncreaseTally() {
				tally = append(tally, t)
			}

			if !t.FinalType().ShouldDecreaseTally() {
				continue
			}

			eatFromTallyUntilSatisfied(t, tally, yearEndTotals)
		}
	}

	for year, amounts := range yearEndTotals {
		var yearTotal float64
		for _, amount := range amounts {
			yearTotal += amount
		}
		yearEndTotals[year]["TOTAL"] = yearTotal
	}

	fmt.Println(yearEndTotals)
	return yearEndTotals, nil
}

func eatFromTallyUntilSatisfied(currentTransaction sharedtypes.Transaction, tally []sharedtypes.Transaction, yet YearEndTotals) YearEndTotals {
	toSubtract := currentTransaction.Amount
	for i, tt := range tally {
		// Skip tallys that have already been counted
		if tt.Amount <= 0 {
			continue
		}

		newAmount := tt.Amount - toSubtract
		actualSubtracted := toSubtract

		// Amount is greater than the current tally item, fall over to the next item
		if newAmount <= 0 {
			actualSubtracted = math.Min(tt.Amount, toSubtract)
			toSubtract = math.Abs(newAmount)
			newAmount = 0
		} else {
			// Nothing else to subtract after this round.
			toSubtract = 0
		}

		fiatValueWhenBought := zarValue(tt.Timestamp, actualSubtracted, tt.WholePriceAtPoint)
		fiatValueWhenSold := zarValue(currentTransaction.Timestamp, actualSubtracted, currentTransaction.WholePriceAtPoint)

		// Yuck
		if yet[taxableYear(currentTransaction.Timestamp)] == nil {
			yet[taxableYear(currentTransaction.Timestamp)] = make(map[string]float64)
		}
		yet[taxableYear(currentTransaction.Timestamp)][tt.Currency] += fiatValueWhenSold - fiatValueWhenBought

		tally[i].Amount = newAmount
		if toSubtract <= float64(0) {
			break
		} else if i+1 >= len(tally) {
			fmt.Println(fmt.Sprintf("WARNING: Trying to sell asset that we don't have. Amount over: %f", toSubtract))
		}
	}

	return yet
}

func uniqueCurrencies(transactions []sharedtypes.Transaction) map[string]bool {
	currencies := make(map[string]bool)

	for _, t := range transactions {
		currencies[t.Currency] = true
	}

	return currencies
}

func taxableYear(timestamp int64) int {
	t := time.Unix(timestamp, 0)
	if t.Month() > time.February {
		return t.Year() + 1
	}
	return t.Year()
}

func zarValue(timestamp int64, amount, wholeValue float64) float64 {
	//wholeValue := fetchZARValueAtTime(timestamp)
	return amount * wholeValue
}

func fetchZARValueAtTime(timestamp int64) float64 {
	return float64(timestamp)
}
