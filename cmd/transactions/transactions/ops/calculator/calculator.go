package calculator

import (
	"context"
	"fmt"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/rates/conversionrate"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/transactions/transactions"
	"github.com/luno/jettison/errors"
	"github.com/luno/jettison/j"
	"github.com/luno/jettison/log"
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

func Calculate(b Backends, fiat string, ts []transactions.Transaction) (YearEndTotals, error) {
	yearEndTotals := make(YearEndTotals)

	for currency := range uniqueCurrencies(ts) {
		var tally []transactions.Transaction
		var prevTimestamp int64

		for _, t := range ts {
			if t.Currency != currency {
				continue
			}

			if t.Timestamp < prevTimestamp {
				return nil, errors.Wrap(transactions.ErrInvalidTransactionOrder, "", j.MKV{
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

			err := eatFromTallyUntilSatisfied(b, fiat, t, tally, yearEndTotals)
			if errors.IsAny(err, conversionrate.ErrNoRatesFound, conversionrate.ErrStoredRateExceedsThreshold) {
				// We don't have rates for this, skip it.
				log.Error(context.TODO(), err)
				continue
			} else if err != nil {
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

func eatFromTallyUntilSatisfied(b Backends, fiat string, currentTransaction transactions.Transaction, tally []transactions.Transaction, yet YearEndTotals) error {
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

		fiatValueWhenBought, err := fiatValue(b, tt.Timestamp, fiat, currentTransaction.Currency, actualSubtracted, tt.WholePriceAtPoint)
		if err != nil {
			return err
		}

		fiatValueWhenSold, err := fiatValue(b, currentTransaction.Timestamp, fiat, currentTransaction.Currency, actualSubtracted, currentTransaction.WholePriceAtPoint)
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
		}
		if i+1 >= len(tally) {
			fmt.Println(fmt.Sprintf("WARNING: Trying to sell asset that we don't have. Amount over: %f", toSubtract))
		}
	}

	return nil
}

func uniqueCurrencies(transactions []transactions.Transaction) map[string]bool {
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

func fiatValue(b Backends, timestamp int64, fiat, coin string, amount float64, wholeValue transactions.FiatPrice) (float64, error) {
	var rate float64
	if wholeValue.Price > 0 && fiat == wholeValue.Fiat {
		rate = wholeValue.Price
	} else {
		r, err := b.RatesClient().ValueAtTime(fiat, coin, timestamp)
		if err != nil {
			return 0, err
		}

		rate = r
	}

	return amount * rate, nil
}
