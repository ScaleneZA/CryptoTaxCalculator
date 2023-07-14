package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func Connect() *sql.DB {
	db, err := sql.Open("sqlite3", "cryptotax.db")
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func ConnectForTesting() *sql.DB {
	dbc, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal(err)
	}

	ResetDB(dbc)

	return dbc
}
