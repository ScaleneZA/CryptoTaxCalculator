package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/calculator"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/filetransformer"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/webhandlers"
)

func main() {
	http.HandleFunc("/", webhandlers.Home)
	http.HandleFunc("/ajax/upload_transform", webhandlers.UploadTransform)
	http.HandleFunc("/ajax/calculate", webhandlers.Calculate)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func main_old() {
	pwd, _ := os.Getwd()

	ts, err := filetransformer.Transform(pwd+"/cmd/filetransformer/testData/LUNO_XBT.csv", filetransformer.TransformTypeLuno)
	if err != nil {
		panic(err)
	}
	tax, err := calculator.Calculate(ts)
	if err != nil {
		panic(err)
	}
	fmt.Println(tax)
}
