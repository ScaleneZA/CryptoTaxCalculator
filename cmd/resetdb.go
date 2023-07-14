package main

import (
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/db"
)

func main() {
	dbc := db.Connect()
	defer dbc.Close()

	db.ResetDB(dbc)
}
