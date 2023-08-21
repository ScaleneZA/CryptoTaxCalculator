package main

import (
	"fmt"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/transactions/di"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/transactions/webserver/webhandlers"
	"net/http"
)

func main() {
	httpServers()
}

func httpServers() {
	b := di.SetupDI()

	http.HandleFunc("/", webhandlers.Home)
	http.HandleFunc("/overrides", webhandlers.UploadTransform)

	o := webhandlers.OverrideHandler{B: b}
	http.HandleFunc("/overrides/override", o.Override)

	c := webhandlers.CalculateHandler{B: b}
	http.HandleFunc("/taxpacks", c.Calculate)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
