package calculator

import (
	"fmt"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/taxcalculator/sharedtypes"
	"math"
	"time"
)

/*
Next Steps:
* Multi-file
* Fetch whole price if not supplied
* Separate Sends, Receives, Buys, Sells - Can maybe use a "Known Addresses" feature.
*   - Any Send/Receive to/from a known address does not affect the tally?
*/

func Calculate(transactions []sharedtypes.Transaction) map[int]map[string]float64 {
	var tally []sharedtypes.Transaction
	// TODO(Export this as a type)
	yearEndTotals := make(map[int]map[string]float64)

	for j, t := range transactions {
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

			if yearEndTotals[taxableYear(t.Timestamp)] == nil {
				yearEndTotals[taxableYear(t.Timestamp)] = make(map[string]float64)
			}
			yearEndTotals[taxableYear(t.Timestamp)][tt.Currency] += fiatValueWhenSold - fiatValueWhenBought

			tally[i].Amount = newAmount
			if toSubtract <= float64(0) {
				break
			} else if i+1 >= len(tally) {
				fmt.Println(fmt.Sprintf("WARNING: Trying to sell asset that we don't have (row-index: %d). Amount over: %f", j, toSubtract))
			}
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
	return yearEndTotals
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
