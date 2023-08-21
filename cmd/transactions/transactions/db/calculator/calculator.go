package calculator

import (
	"database/sql"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/transactions/transactions"
	"strings"
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
	stmt, err := db.Prepare("SELECT uid, overridden_type	 FROM calculator_transaction_overrides WHERE uid = ? LIMIT 1")
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

	ot, err := scanRow(result)
	if err != nil {
		return transactions.TypeUnknown, err
	}

	if err := result.Err(); err != nil {
		return transactions.TypeUnknown, err
	}

	return ot.OverriddenType, nil
}

func ListByTypeByUid(db *sql.DB, uids []string) ([]transactions.OverrideType, error) {
	if len(uids) == 0 {
		return nil, nil
	}

	// Create placeholders for the IN clause
	placeholders := make([]string, len(uids))
	args := make([]interface{}, len(uids))
	for i := range uids {
		placeholders[i] = "?"
		args[i] = uids[i]
	}

	// Construct the SQL query
	query := "SELECT uid, overridden_type FROM calculator_transaction_overrides WHERE uid IN (" + strings.Join(placeholders, ", ") + ")"

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ots []transactions.OverrideType
	for rows.Next() {
		ot, err := scanRow(rows)
		if err != nil {
			return nil, err
		}
		ots = append(ots, *ot)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ots, nil
}

func scanRow(rows *sql.Rows) (*transactions.OverrideType, error) {
	var ot transactions.OverrideType
	if err := rows.Scan(&ot.UID, &ot.OverriddenType); err != nil {
		return nil, err
	}
	return &ot, nil
}
