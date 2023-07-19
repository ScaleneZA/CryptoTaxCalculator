package server

import "database/sql"

type Backends interface {
	DB() *sql.DB
}
