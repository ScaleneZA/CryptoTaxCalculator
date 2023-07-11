package main

import (
	"fmt"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/conversionrate"
	"log"
	"net/http"

	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/webhandlers"
)

func main() {
	err := conversionrate.SyncAll()
	if err != nil {
		log.Println("Pair Syncing Failed:")
		log.Println(err.Error())
	}

	http.HandleFunc("/", webhandlers.Home)
	http.HandleFunc("/ajax/upload", webhandlers.UploadTransform)
	http.HandleFunc("/ajax/calculate", webhandlers.Calculate)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
