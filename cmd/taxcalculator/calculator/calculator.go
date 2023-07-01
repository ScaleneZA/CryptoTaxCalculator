package calculator

import (
	"fmt"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/taxcalculator/sharedtypes"
	"math"
	"time"
)

func Calculate(transactions []sharedtypes.Transaction) map[int]map[string]float64 {
	var tally []sharedtypes.Transaction
	// TODO(Export this as a type)
	taxableAmounts := make(map[int]map[string]float64)

	for _, t := range transactions {
		// TODO: Throw error if current transaction has a date before the previous one.

		if t.Typ == sharedtypes.TypeBuy {
			tally = append(tally, t)
			continue
		}

		toSubtract := t.Amount
		for i, tt := range tally {
			// Skip tallys that have already been counted
			if tt.Amount <= 0 {
				continue
			}

			if tt.Currency != t.Currency {
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
			fiatValueWhenSold := zarValue(t.Timestamp, actualSubtracted, t.WholePriceAtPoint)

			if taxableAmounts[taxableYear(t.Timestamp)] == nil {
				taxableAmounts[taxableYear(t.Timestamp)] = make(map[string]float64)
			}
			taxableAmounts[taxableYear(t.Timestamp)][tt.Currency] += fiatValueWhenSold - fiatValueWhenBought

			tally[i].Amount = newAmount
			if toSubtract <= 0 {
				break
			} else if i+2 > len(tally) {
				panic(fmt.Sprintf("FATAL: Trying to sell asset that we don't have. Amount over: %f", toSubtract))
			}
		}
	}

	fmt.Println(taxableAmounts)
	return taxableAmounts
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
