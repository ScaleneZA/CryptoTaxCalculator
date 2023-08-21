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
	"sort"
	"time"
)

/*
Next Steps:
* Way to initialize at a point in time
*/

type yearGainsMap map[int]map[string]Gain
type yearBalancesMap map[int]map[string]Balance
type totalBalancesMap map[string]Balance

func Calculate(b Backends, fiat string, ts []transactions.Transaction) ([]YearEndTotal, error) {
	firstYear := taxableYear(ts[0].Timestamp)
	lastYear := taxableYear(ts[len(ts)-1].Timestamp)

	if firstYear > lastYear {
		return nil, errors.Wrap(transactions.ErrInvalidTransactionOrder, "", j.MKV{
			"first_year": firstYear,
			"last_year":  lastYear,
		})
	}

	balanceTotals := make(map[string]Balance)
	yearBalanceTotals := initYearBalanceMap(firstYear, lastYear)
	yearGainTotals := make(yearGainsMap)
	excludeCurrencies := make(map[string]bool)

	for currency := range uniqueCurrencies(ts) {
		var tally []transactions.Transaction
		var prevTimestamp int64

		// TODO: Look at ways to simplify the code because now we are tracking this year.
		for year := firstYear; year <= lastYear; year++ {
			fmt.Println(currency)
			for _, t := range ts {
				if t.Currency != currency {
					continue
				}

				if taxableYear(t.Timestamp) != year {
					continue
				}

				if t.Timestamp < prevTimestamp {
					return nil, errors.Wrap(transactions.ErrInvalidTransactionOrder, "", j.MKV{
						"current_timestamp":  t.Timestamp,
						"previous_timestamp": prevTimestamp,
					})
				}

				if t.FinalType().ShouldIncreaseTally() {
					tally = append(tally, t)
					balanceTotals[t.Currency] = Balance{
						Asset:  t.Currency,
						Amount: balanceTotals[t.Currency].Amount + t.Amount,
					}
				}

				if !t.FinalType().ShouldDecreaseTally() {
					prevTimestamp = t.Timestamp
					yearBalanceTotals[taxableYear(t.Timestamp)][t.Currency] = balanceTotals[t.Currency]

					continue
				}

				fmt.Print(".")

				err := eatFromTallyUntilSatisfied(b, fiat, t, tally, yearGainTotals, balanceTotals)
				if errors.IsAny(err, conversionrate.ErrNoRatesFound, conversionrate.ErrStoredRateExceedsThreshold) {
					log.Error(context.TODO(), err)
					// We don't have a particular rate for this currency, skip it.
					excludeCurrencies[currency] = true
					break
				} else if err != nil {
					return nil, err
				}

				prevTimestamp = t.Timestamp
			}

			if balanceTotals[currency].Asset != "" {
				yearBalanceTotals[year][currency] = balanceTotals[currency]
			}
		}
	}

	return sumUpTotals(yearGainTotals, yearBalanceTotals, excludeCurrencies), nil
}

func initYearBalanceMap(firstYear, lastYear int) yearBalancesMap {
	yearBalanceTotals := make(yearBalancesMap)
	for y := firstYear; y <= lastYear; y++ {
		yearBalanceTotals[y] = make(map[string]Balance)
	}

	return yearBalanceTotals
}

func eatFromTallyUntilSatisfied(b Backends, fiat string, currentTransaction transactions.Transaction, tally []transactions.Transaction, ygt yearGainsMap, bt totalBalancesMap) error {
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
		fmt.Print("\\")

		// We only need to lookup the value when we actually bought something.
		var fiatValueWhenBought float64
		if tt.FinalType() == transactions.TypeBuy {
			paid, err := fiatValue(b, tt.Timestamp, fiat, tt.Currency, actualSubtracted, tt.WholePriceAtPoint)
			if err != nil {
				return err
			}
			fiatValueWhenBought = paid
		}

		fmt.Print("_")

		fiatValueWhenSold, err := fiatValue(b, currentTransaction.Timestamp, fiat, tt.Currency, actualSubtracted, currentTransaction.WholePriceAtPoint)
		if err != nil {
			return err
		}

		fmt.Print("/")

		year := taxableYear(currentTransaction.Timestamp)

		if ygt[year] == nil {
			ygt[year] = make(map[string]Gain)
		}

		// Reassign asset with updated amounts
		ygt[year][tt.Currency] = Gain{
			Asset:    tt.Currency,
			Amount:   ygt[year][tt.Currency].Amount + actualSubtracted,
			Costs:    ygt[year][tt.Currency].Costs + fiatValueWhenBought,
			Proceeds: ygt[year][tt.Currency].Proceeds + fiatValueWhenSold,
		}

		bt[tt.Currency] = Balance{
			Asset:  tt.Currency,
			Amount: bt[tt.Currency].Amount - actualSubtracted,
		}

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
		if t.Currency == "" {
			continue
		}
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

func sumUpTotals(yearGainTotals yearGainsMap, yearBalanceTotals yearBalancesMap, excludeCurrencies map[string]bool) []YearEndTotal {
	var ygts []YearEndTotal
	for year, balances := range yearBalanceTotals {
		yt := YearEndTotal{
			Year: year,
		}

		yearGainTotal := Gain{
			Asset: "TOTAL",
		}

		for currency, bs := range balances {
			if excludeCurrencies[currency] {
				delete(yearGainTotals[year], currency)
				continue
			}

			yt.Balances = append(yt.Balances, bs)

			g, ok := yearGainTotals[year][currency]
			if ok {
				yt.Gains = append(yt.Gains, g)

				yearGainTotal.Costs += g.Costs
				yearGainTotal.Proceeds += g.Proceeds
			}
		}

		yt.Gains = append(yt.Gains, yearGainTotal)

		sort.Slice(yt.Gains, func(i, j int) bool {
			return (yt.Gains[i].Proceeds - yt.Gains[i].Costs) > (yt.Gains[j].Proceeds - yt.Gains[j].Costs)
		})

		sort.Slice(yt.Balances, func(i, j int) bool {
			return yt.Balances[i].Amount > yt.Balances[j].Amount
		})

		ygts = append(ygts, yt)
	}

	sort.Slice(ygts, func(i, j int) bool {
		return ygts[i].Year < ygts[j].Year
	})

	return ygts
}
