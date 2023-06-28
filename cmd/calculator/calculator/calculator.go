package calculator

import(
	"fmt"
	"math"
	"time"
)

func Calculate(transactions []Transaction) {
	var tally []Transaction
	taxableAmounts := make(map[int]float64)

	for j, t := range transactions {
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

			amount := tt.Amount - toSubtract

			absAmnt := math.Min(tt.Amount, toSubtract)
			zarValueB := zarValue(tt.Timestamp, absAmnt, tt.WholePriceAtPoint)
			zarValueS := zarValue(t.Timestamp, absAmnt, t.WholePriceAtPoint)

			fmt.Println(j, "Event value when bought:", "(R", zarValueB, ")", "Event value when sold:", "(R", zarValueS, ") Taxable amount: R", zarValueS-zarValueB)

			taxableAmounts[taxableYear(t.Timestamp)] += zarValueS-zarValueB

			// Amount is greater than the current tally item, fall over to the next item
			if amount <= 0 {
				toSubtract = math.Abs(amount)
				amount = 0
			} else {
				toSubtract = 0
			}

			tally[i].Amount = amount

			if toSubtract <= 0 {
				break
			} else if i+2 > len(tally)  {
				panic(fmt.Sprintf("FATAL: Trying to sell asset that we don't have. Amount over: %f",toSubtract))
			}
		}
	}

	fmt.Println(taxableAmounts)
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
