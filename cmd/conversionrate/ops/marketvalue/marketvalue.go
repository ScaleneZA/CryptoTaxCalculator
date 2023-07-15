package marketvalue

import (
	"errors"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/conversionrate/db/markets"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/conversionrate/sharedtypes"
)

func ValueAtTime(b Backends, from, to string, timestamp int64) (float64, error) {
	g, err := buildCurrencyGraph(b, timestamp)
	if err != nil {
		return 0, err
	}

	rate, found := g.findRate(from, to)

	if !found {
		return 0, errors.New("no rate found")
	}

	return rate, nil
}

type currencyGraph map[string]map[string]float64

func (g currencyGraph) findRate(from, to string) (float64, bool) {
	return g[from][to], true
}

func buildCurrencyGraph(b Backends, timestamp int64) (currencyGraph, error) {
	var allRatesAtPoint []sharedtypes.MarketPair

	for _, p := range sharedtypes.AllPairs() {
		closest, err := FindClosest(b, p, timestamp)
		if err != nil {
			return nil, err
		}
		allRatesAtPoint = append(allRatesAtPoint, *closest)
	}

	graph := make(map[string]map[string]float64)
	for _, mp := range allRatesAtPoint {
		if _, ok := graph[mp.FromCurrency]; !ok {
			graph[mp.FromCurrency] = make(map[string]float64)
		}

		graph[mp.FromCurrency][mp.ToCurrency] = mp.Close
	}

	return graph, nil
}

func FindClosest(b Backends, p sharedtypes.Pair, timestamp int64) (*sharedtypes.MarketPair, error) {
	closestBefore, _ := markets.FindClosestToBefore(b.DB(), p.FromCurrency, p.ToCurrency, timestamp)
	closestAfter, _ := markets.FindClosestToAfter(b.DB(), p.FromCurrency, p.ToCurrency, timestamp)

	if closestAfter == nil && closestBefore == nil {
		return nil, errors.New("cannot find a market price for: " + p.FromCurrency + "/" + p.ToCurrency)
	} else if closestBefore == nil {
		return closestAfter, nil
	} else if closestAfter == nil {
		return closestBefore, nil
	}

	if timestamp-closestBefore.Timestamp < closestAfter.Timestamp-timestamp {
		return closestBefore, nil
	}

	return closestAfter, nil
}
