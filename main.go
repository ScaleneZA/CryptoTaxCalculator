package main

import (
	"fmt"
	"math"
	"time"
)

type transactionType int
const (
	typeBuy = 0
	typeSell	= 1
)

type transaction struct {
	typ transactionType
	amount float64
	timestamp int64
	// TODO: Make this dynamic if possible
	wholePriceAtPoint float64
}

func main() {
	transactions := []transaction{
		{
			typ:    typeBuy,
			amount: 0.56,
			timestamp: 1519812503,
			wholePriceAtPoint: 100,
		},
		{
			typ:    typeBuy,
			amount: 1.2,
			timestamp: 1535450903,
			wholePriceAtPoint: 200,
		},
		{
			typ: typeSell,
			amount: 0.25,
			timestamp: 1656410903,
			wholePriceAtPoint: 300,
		},
		{
			typ: typeSell,
			amount: 1.25,
			timestamp: 1687946903,
			wholePriceAtPoint: 400,
		},
	}

	var tally []transaction
	taxableAmounts := make(map[int]float64)

	for j, t := range transactions {
		if t.typ == typeBuy {
			tally = append(tally, t)
			continue
		}

		toSubtract := t.amount
		for i, tt := range tally {
			// Skip tallys that have already been counted
			if tt.amount <= 0 {
				continue
			}

			amount := tt.amount - toSubtract

			absAmnt := math.Min(tt.amount, toSubtract)
			zarValueB := zarValue(tt.timestamp, absAmnt, tt.wholePriceAtPoint)
			zarValueS := zarValue(t.timestamp, absAmnt, t.wholePriceAtPoint)

			fmt.Println(j, "Event value when bought:", "(R", zarValueB, ")", "Event value when sold:", "(R", zarValueS, ") Taxable amount: R", zarValueS-zarValueB)

			taxableAmounts[taxableYear(t.timestamp)] += zarValueS-zarValueB

			// Amount is greater than the current tally item, fall over to the next item
			if amount <= 0 {
				toSubtract = math.Abs(amount)
				amount = 0
			} else {
				toSubtract = 0
			}

			tally[i].amount = amount

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
