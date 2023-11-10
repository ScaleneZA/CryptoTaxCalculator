package webhandlers

import (
	"database/sql"
)

type Backends interface {
	DB() *sql.DB
}
