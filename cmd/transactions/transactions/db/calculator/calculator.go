package calculator

import (
	"database/sql"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/transactions/transactions"
	"sync"
)

var mtx sync.Mutex

func Upsert(db *sql.DB, uid string, overrideType transactions.TransactionType) (int64, error) {
	mtx.Lock()
	defer mtx.Unlock()

	stmt, err := db.Prepare(`
		INSERT INTO calculator_transaction_overrides (uid, overridden_type, created_at)
		VALUES (?, ?, strftime('%Y-%m-%d %H-%M-%S','now'))
		ON CONFLICT(uid) DO UPDATE SET overridden_type=?, updated_at=strftime('%Y-%m-%d %H-%M-%S','now');
`)

	result, err := stmt.Exec(uid, overrideType, overrideType)
	if err != nil {
		return 0, err
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastInsertID, nil
}

func LookupTypeByUid(db *sql.DB, uid string) (transactions.TransactionType, error) {
	stmt, err := db.Prepare("SELECT overridden_type FROM calculator_transaction_overrides WHERE uid = ? LIMIT 1")
	if err != nil {
		return transactions.TypeUnknown, err
	}
	defer stmt.Close()

	result, err := stmt.Query(uid)
	if err != nil {
		return transactions.TypeUnknown, err
	}
	defer result.Close()

	result.Next()

	var typ transactions.TransactionType
	if err := result.Scan(&typ); err != nil {
		return transactions.TypeUnknown, err
	}

	if err := result.Err(); err != nil {
		return transactions.TypeUnknown, err
	}

	return typ, nil
}
