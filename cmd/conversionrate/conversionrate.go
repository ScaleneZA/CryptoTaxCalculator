package conversionrate

import (
	"fmt"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/conversionrate/sharedtypes"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/db/markets"
)

func MarketValueAtTime(b Backends, from, to string, timestamp int) (float64, error) {
	g, err := buildCurrencyGraph(b, timestamp)
	if err != nil {
		return 0, err
	}

	rate, found := g.findRate(from, to)

	if found {
		fmt.Printf("Exchange rate: %.2f USD to JPY\n", rate)
	} else {
		fmt.Println("No exchange rate found")
	}

	return 0.00, nil
}

type currencyGraph map[string]map[string]float64

func (g currencyGraph) findRate(from, to string) (float64, bool) {
	rates, visited := g.findExchangeRateHelper(from, to, make(map[string]bool))
	return rates[to], visited[to]
}

// Helper function for finding the exchange rate between two currencies using DFS
func (g currencyGraph) findExchangeRateHelper(from, to string, visited map[string]bool) (map[string]float64, map[string]bool) {
	visited[from] = true

	if from == to {
		return map[string]float64{from: 1}, visited
	}

	for currency, rate := range g[from] {
		if !visited[currency] {
			rates, visited := g.findExchangeRateHelper(currency, to, visited)
			if rates != nil {
				newRate := rate * rates[currency]
				rates[from] = 1 / newRate
				return rates, visited
			}
		}
	}

	return nil, visited
}

func buildCurrencyGraph(b Backends, timestamp int) (currencyGraph, error) {
	var allRatesAtPoint []sharedtypes.MarketPair

	for _, p := range sharedtypes.AllPairs() {
		closestBefore, err := markets.FindClosestToBefore(b.DB(), p.Currency1, p.Currency2, timestamp)
		if err != nil {
			return currencyGraph{}, nil
		}
		closestAfter, err := markets.FindClosestToAfter(b.DB(), p.Currency1, p.Currency2, timestamp)
		if err != nil {
			return currencyGraph{}, nil
		}

		if (timestamp - closestBefore.Timestamp) > (closestAfter.Timestamp - timestamp) {
			allRatesAtPoint = append(allRatesAtPoint, *closestAfter)
		} else {
			allRatesAtPoint = append(allRatesAtPoint, *closestBefore)
		}
	}

	graph := make(map[string]map[string]float64)
	for _, mp := range allRatesAtPoint {
		if _, ok := graph[mp.Currency1]; !ok {
			graph[mp.Currency1] = make(map[string]float64)
		}

		graph[mp.Currency1][mp.Currency2] = mp.Close
	}

	return graph, nil
}
