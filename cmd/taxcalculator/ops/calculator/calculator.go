package calculator

import (
	"fmt"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/conversionrate/ops/marketvalue"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/di"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/taxcalculator"
	"github.com/luno/jettison/errors"
	"github.com/luno/jettison/j"
	"math"
	"time"
)

/*
Next Steps:
* Multi-file
* Way to initialize at a point in time
* Fetch whole price if not supplied
* Separate Sends, Receives, Buys, Sells - Override detected types
*/

func Calculate(fiat string, transactions []taxcalculator.Transaction) (YearEndTotals, error) {
	yearEndTotals := make(YearEndTotals)

	for currency := range uniqueCurrencies(transactions) {
		var tally []taxcalculator.Transaction
		var prevTimestamp int64

		for _, t := range transactions {
			if t.Currency != currency {
				continue
			}

			if t.Timestamp < prevTimestamp {
				return nil, errors.Wrap(taxcalculator.ErrInvalidTransactionOrder, "", j.MKV{
					"current_timestamp":  t.Timestamp,
					"previous_timestamp": prevTimestamp,
				})
			}
			prevTimestamp = t.Timestamp

			if t.FinalType().ShouldIncreaseTally() {
				tally = append(tally, t)
			}

			if !t.FinalType().ShouldDecreaseTally() {
				continue
			}

			err := eatFromTallyUntilSatisfied(fiat, t, tally, yearEndTotals)
			if err != nil {
				return nil, err
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
	return yearEndTotals, nil
}

func eatFromTallyUntilSatisfied(fiat string, currentTransaction taxcalculator.Transaction, tally []taxcalculator.Transaction, yet YearEndTotals) error {
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

		fiatValueWhenBought, err := fiatValue(tt.Timestamp, fiat, currentTransaction.Currency, actualSubtracted, tt.WholePriceAtPoint)
		if err != nil {
			return err
		}

		fiatValueWhenSold, err := fiatValue(currentTransaction.Timestamp, fiat, currentTransaction.Currency, actualSubtracted, currentTransaction.WholePriceAtPoint)
		if err != nil {
			return err
		}

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

	return nil
}

func uniqueCurrencies(transactions []taxcalculator.Transaction) map[string]bool {
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

func fiatValue(timestamp int64, fiat, coin string, amount, wholeValue float64) (float64, error) {
	if wholeValue > 0 {
		return amount * wholeValue, nil
	}

	// TODO: Make this a GRPC call
	rate, err := marketvalue.ValueAtTime(di.SetupDI(), fiat, coin, timestamp)
	if err != nil {
		return 0, err
	}

	return amount * rate, nil
}
