package main

import (
	"fmt"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/db"
	"io/ioutil"
	"log"
)

func main() {
	dbc := db.Connect()
	defer dbc.Close()

	// Read SQL file contents
	sqlBytes, err := ioutil.ReadFile("cmd/db/schema.sql")
	if err != nil {
		log.Fatal(err)
	}
	sqlQuery := string(sqlBytes)

	// Execute the SQL query
	_, err = dbc.Exec(sqlQuery)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Schema reset and all data deleted.")
}
