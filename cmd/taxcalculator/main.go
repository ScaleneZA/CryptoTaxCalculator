package main

import (
	"fmt"
	"os"

	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/taxcalculator/calculator"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/taxcalculator/transformer"
)

func main() {
	pwd, _ := os.Getwd()

	ts, err := transformer.Transform(pwd+"/cmd/taxcalculator/transformer/testData/LUNO_XBT.csv", transformer.TransformTypeLuno)
	if err != nil {
		panic(err)
	}
	tax := calculator.Calculate(ts)
	fmt.Println(tax)
}
