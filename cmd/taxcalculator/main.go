package main

import (
	"fmt"
	"os"

	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/taxcalculator/calculator"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/taxcalculator/filetransformer"
)

func main() {
	pwd, _ := os.Getwd()

	ts, err := filetransformer.Transform(pwd+"/cmd/taxcalculator/filetransformer/testData/LUNO_XBT.csv", filetransformer.TransformTypeLuno)
	if err != nil {
		panic(err)
	}
	tax, err := calculator.Calculate(ts)
	if err != nil {
		panic(err)
	}
	fmt.Println(tax)
}
