package main

import (
	db2 "github.com/ScaleneZA/CryptoTaxCalculator/cmd/rates/db"
)

func resetDB() {
	dbc := db2.Connect()
	defer dbc.Close()

	db2.ResetDB(dbc)
}
