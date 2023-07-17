package main

import (
	"fmt"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/conversionrate/ops/sync"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/di"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/taxcalculator/webhandlers"
	"log"
	"net/http"
	"time"
)

func main() {
	b := di.SetupDI()
	go syncCurrenciesForever(b)

	http.HandleFunc("/", webhandlers.Home)
	http.HandleFunc("/ajax/upload", webhandlers.UploadTransform)
	http.HandleFunc("/ajax/calculate", webhandlers.Calculate)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func syncCurrenciesForever(b di.Backends) {
	for {
		err := sync.SyncAll(b)
		if err != nil {
			log.Println("Pair Syncing Failed:")
			log.Println(err.Error())
		}

		log.Println("Synced all currencies")

		time.Sleep(time.Hour * 12)
	}
}
