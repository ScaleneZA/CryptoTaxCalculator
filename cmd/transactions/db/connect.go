package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func Connect() *sql.DB {
	dbc, err := sql.Open("sqlite3", "cryptotax.db")
	if err != nil {
		log.Fatal(err)
	}

	return dbc
}

func ConnectForTesting() *sql.DB {
	dbc, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		log.Fatal(err)
	}

	ResetDB(dbc)

	return dbc
}
