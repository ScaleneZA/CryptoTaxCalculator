package db

import (
	"database/sql"
	"fmt"
	"github.com/ScaleneZA/CryptoTaxCalculator/config"
	"log"
	"os"
	"path/filepath"
	"sync"
)

// Unfortunately this is necessary for the in-memory sqlite, as it is not thread-safe.
var mtx sync.Mutex

func ResetDB(dbc *sql.DB) {
	mtx.Lock()
	defer mtx.Unlock()
	schemaPath := filepath.Join(filepath.Dir(config.WorkingDirectory), "cmd", "rates", "db", "schema.sql")
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
