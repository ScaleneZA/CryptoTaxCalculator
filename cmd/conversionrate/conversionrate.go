package conversionrate

import (
	"fmt"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/conversionrate/sharedtypes"
)

func MarketValueAtTime(from, to string, timestamp int) (float64, error) {
	g := buildCurrencyGraph(timestamp)

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

func buildCurrencyGraph(timestamp int) currencyGraph {
	// Find all currency rates closest to timestamp
	allRatesAtPoint := []sharedtypes.MarketPair{}

	graph := make(map[string]map[string]float64)
	for _, mp := range allRatesAtPoint {
		if _, ok := graph[mp.Currency1]; !ok {
			graph[mp.Currency1] = make(map[string]float64)
		}

		graph[mp.Currency1][mp.Currency2] = mp.Close
	}

	return graph
}
