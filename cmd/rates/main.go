package main

import (
	"fmt"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/rates/conversionrate/ops/sync"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/rates/conversionrate/server"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/rates/db"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/rates/di"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/rates/webserver/webhandlers"
	"log"
	"net/http"
	"time"
)

func main() {
	b := di.SetupDI()

	//resetDB()
	//go syncCurrenciesForever(b)

	// HTTP
	hh := webhandlers.HomeHandler{B: b}
	http.HandleFunc("/", hh.Handle)
	http.HandleFunc("/closest", hh.ClosestForTime)

	err := http.ListenAndServe(":8181", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}

	// GRPC
	err = server.Serve(b)
	if err != nil {
		panic(err)
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

func resetDB() {
	dbc := db.Connect()
	defer dbc.Close()

	db.ResetDB(dbc)
}
