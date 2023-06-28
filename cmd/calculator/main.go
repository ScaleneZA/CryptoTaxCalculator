package main

import (
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/calculator/calculator"
)

func main() {
	ts := []calculator.Transaction{
		{
			Typ:    calculator.TypeBuy,
			Amount: 0.56,
			Timestamp: 1519812503,
			WholePriceAtPoint: 100,
		},
		{
			Typ:    calculator.TypeBuy,
			Amount: 1.2,
			Timestamp: 1535450903,
			WholePriceAtPoint: 200,
		},
		{
			Typ: calculator.TypeSell,
			Amount: 0.25,
			Timestamp: 1656410903,
			WholePriceAtPoint: 300,
		},
		{
			Typ: calculator.TypeSell,
			Amount: 1.25,
			Timestamp: 1687946903,
			WholePriceAtPoint: 400,
		},
	}

	calculator.Calculate(ts)
}
