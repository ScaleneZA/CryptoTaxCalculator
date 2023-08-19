package marketvalue

import (
	"fmt"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/rates/conversionrate"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/rates/conversionrate/db/markets"
	"github.com/luno/jettison/errors"
	"github.com/luno/jettison/j"
)

// depthLimit is used for the recursive path finding to make sure we don't end up in an
// endless loop if there are currencies that can cause loops. E.G. USD -> ETH -> BTC -> USD
const depthLimit = 3

func ValueAtTime(b Backends, from, to string, timestamp int64) (float64, error) {
	mps, err := closestMarketPairsAtPoint(b, timestamp)
	if err != nil {
		return 0, err
	}

	val, err := findRate(mps, from, to, timestamp, 0)
	if errors.Is(err, conversionrate.ErrNoRatesFound) {
		return 0, errors.Wrap(err, "", j.MKV{
			"from":      from,
			"to":        to,
			"timestamp": timestamp,
			"mps_len":   len(mps),
		})
	}

	return val, err
}

func closestMarketPairsAtPoint(b Backends, timestamp int64) ([]conversionrate.MarketPair, error) {
	var allRatesAtPoint []conversionrate.MarketPair

	for _, p := range conversionrate.AllPairs() {
		closest, err := FindClosest(b, p, timestamp)
		if errors.Is(err, conversionrate.ErrNoMarket) {
			// NoReturnErr: Skip markets we don't have data for.
			continue
		} else if err != nil {
			return nil, err
		}
		allRatesAtPoint = append(allRatesAtPoint, *closest)
	}

	return allRatesAtPoint, nil
}

func FindClosest(b Backends, p conversionrate.Pair, timestamp int64) (*conversionrate.MarketPair, error) {
	closestBefore, _ := markets.FindClosestToBefore(b.DB(), p.FromCurrency, p.ToCurrency, timestamp)
	closestAfter, _ := markets.FindClosestToAfter(b.DB(), p.FromCurrency, p.ToCurrency, timestamp)

	var closest *conversionrate.MarketPair
	if closestAfter == nil && closestBefore == nil {
		return nil, errors.Wrap(conversionrate.ErrNoMarket, "", j.MKV{
			"pair": p.String(),
		})
	} else if closestBefore == nil {
		closest = closestAfter
	} else if closestAfter == nil {
		closest = closestBefore
	} else if timestamp-closestBefore.Timestamp < closestAfter.Timestamp-timestamp {
		closest = closestBefore
	} else {
		closest = closestAfter
	}

	if rateTimeExceedsThreshold(timestamp, closest.Timestamp) {
		return nil, errors.Wrap(conversionrate.ErrNoMarket, "", j.MKV{
			"pair":      p.String(),
			"timestamp": timestamp,
		})
	}

	return closest, nil
}

// findRate currently only works for increasing value pairs. For example ZAR -> USD -> BTC. It
// would not work in reverse, for example USD -> BTC -> ETH unless the values imported are negative
// and already reversed.
func findRate(mps []conversionrate.MarketPair, from, to string, timestamp int64, depth int) (float64, error) {
	depth++
	for _, mp := range mps {
		if mp.FromCurrency != from {
			continue
		}
		if mp.ToCurrency == to {
			return mp.Close, nil
		}

		if depth <= depthLimit {
			rate, err := findRate(mps, mp.ToCurrency, to, timestamp, depth)
			if errors.Is(err, conversionrate.ErrNoRatesFound) {
				fmt.Print("X")
				// Depth first search cannot bubble up error until last path is done.
				continue
			} else if err != nil {
				return 0, err
			}

			return mp.Close * rate, nil
		}
	}

	return 0, errors.Wrap(conversionrate.ErrNoRatesFound, "")
}

func rateTimeExceedsThreshold(timestamp, closestTimestamp int64) bool {
	const week = 604800
	return (timestamp-closestTimestamp) > week || (closestTimestamp-timestamp) > week
}
