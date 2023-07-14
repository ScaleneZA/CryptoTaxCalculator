package db

import (
	"database/sql"
	"fmt"
	"github.com/ScaleneZA/CryptoTaxCalculator/config"
	"log"
	"os"
	"path/filepath"
)

func ResetDB(dbc *sql.DB) {
	schemaPath := filepath.Join(filepath.Dir(config.WorkingDirectory), "cmd", "db", "schema.sql")
	log.Println(schemaPath)

	sqlBytes, err := os.ReadFile(schemaPath)
	if err != nil {
		log.Fatal(err)
	}
	sqlQuery := string(sqlBytes)

	_, err = dbc.Exec(sqlQuery)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Schema reset and all data deleted.")
}
