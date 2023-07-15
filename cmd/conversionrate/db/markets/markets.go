package markets

import (
	"database/sql"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/conversionrate/sharedtypes"
	"log"
	"sync"
)

// Unfortunately this is necessary for the in-memory sqlite, as it is not thread-safe.
var mtx sync.Mutex

func FindClosestToBefore(db *sql.DB, from, to string, timestamp int64) (*sharedtypes.MarketPair, error) {
	return findWhere(db, "`from` = ? AND `to` = ? AND timestamp <= ? ORDER BY timestamp DESC", from, to, timestamp)
}

func FindClosestToAfter(db *sql.DB, from, to string, timestamp int64) (*sharedtypes.MarketPair, error) {
	return findWhere(db, "`from` = ? AND `to` = ? AND timestamp >= ? ORDER BY timestamp", from, to, timestamp)
}

func findWhere(db *sql.DB, where string, vars ...any) (*sharedtypes.MarketPair, error) {
	mtx.Lock()
	defer mtx.Unlock()
	stmt, err := db.Prepare("SELECT timestamp, `from`, `to`, open, high, low, close FROM markets WHERE " + where + " LIMIT 1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	result, err := stmt.Query(vars...)
	if err != nil {
		return nil, err
	}

	result.Next()
	mp, err := scanRow(result)
	if err != nil {
		return nil, err
	}

	if err := result.Err(); err != nil {
		return nil, err
	}

	return mp, nil
}

func ListAll(db *sql.DB) ([]sharedtypes.MarketPair, error) {
	rows, err := db.Query("SELECT timestamp, `from`, `to`, open, high, low, close FROM markets")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mps []sharedtypes.MarketPair
	for rows.Next() {
		mp, err := scanRow(rows)
		if err != nil {
			return nil, err
		}
		mps = append(mps, *mp)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return mps, nil
}

func scanRow(rows *sql.Rows) (*sharedtypes.MarketPair, error) {
	var mp sharedtypes.MarketPair
	if err := rows.Scan(&mp.Timestamp, &mp.FromCurrency, &mp.ToCurrency, &mp.Open, &mp.High, &mp.Low, &mp.Close); err != nil {
		return nil, err
	}
	return &mp, nil
}

func Create(db *sql.DB, pair sharedtypes.MarketPair) (int64, error) {
	stmt, err := db.Prepare("INSERT INTO markets(timestamp, `from`, `to`, open, high, low, close) VALUES(?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(pair.Timestamp, pair.FromCurrency, pair.ToCurrency, pair.Open, pair.High, pair.Low, pair.Close)
	if err != nil {
		return 0, err
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastInsertID, nil
}

// TODO: Look into why this seems to be dropping the table
func Truncate(db *sql.DB) error {
	_, err := db.Query("DELETE FROM markets;")
	if err != nil {
		return err
	}

	_, err = db.Query("VACUUM;")
	if err != nil {
		return err
	}

	_, err = db.Query("DELETE FROM sqlite_sequence WHERE name='markets';")
	if err != nil {
		log.Println("Failed to delete from sqlite_sequence DB. Could be that there is no table of this name?")
		// Not critical error
	}

	return nil
}
