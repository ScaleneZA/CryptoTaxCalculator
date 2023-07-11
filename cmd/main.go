package main

import (
	"fmt"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/conversionrate/sync"
	webhandlers2 "github.com/ScaleneZA/CryptoTaxCalculator/cmd/taxcalculator/webhandlers"
	"log"
	"net/http"
	"time"
)

func main() {
	go syncCurrenciesForever()

	http.HandleFunc("/", webhandlers2.Home)
	http.HandleFunc("/ajax/upload", webhandlers2.UploadTransform)
	http.HandleFunc("/ajax/calculate", webhandlers2.Calculate)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func syncCurrenciesForever() {
	for {
		err := sync.SyncAll()
		if err != nil {
			log.Println("Pair Syncing Failed:")
			log.Println(err.Error())
		}

		time.Sleep(time.Hour * 12)
	}
}
