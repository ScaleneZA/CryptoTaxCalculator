package calculator

import(
	"fmt"
	"math"
	"time"
)

func Calculate(transactions []Transaction) map[int]float64{
	var tally []Transaction
	taxableAmounts := make(map[int]float64)

	for _, t := range transactions {
		if t.Typ == TypeBuy {
			tally = append(tally, t)
			continue
		}

		toSubtract := t.Amount
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
			fiatValueWhenSold := zarValue(t.Timestamp, actualSubtracted, t.WholePriceAtPoint)
			taxableAmounts[taxableYear(t.Timestamp)] += fiatValueWhenSold-fiatValueWhenBought

			tally[i].Amount = newAmount
			if toSubtract <= 0 {
				break
			} else if i+2 > len(tally)  {
				panic(fmt.Sprintf("FATAL: Trying to sell asset that we don't have. Amount over: %f",toSubtract))
			}
		}
	}

	fmt.Println(taxableAmounts)
	return taxableAmounts
}

func taxableYear(timestamp int64) int{
	t := time.Unix(timestamp, 0)
	if t.Month() > time.February {
		return t.Year()+1
	}
	return t.Year()
}

func zarValue(timestamp int64, amount, wholeValue float64) float64{
	//wholeValue := fetchZARValueAtTime(timestamp)
	return amount * wholeValue
}

func fetchZARValueAtTime(timestamp int64) float64{
	return float64(timestamp)
}
