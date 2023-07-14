package markets

import (
	"database/sql"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/conversionrate/sharedtypes"
)

func ListAll(db *sql.DB) ([]sharedtypes.MarketPair, error) {
	rows, err := db.Query("SELECT timestamp, `from`, `to`, open, high, low, close FROM markets")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mps []sharedtypes.MarketPair
	for rows.Next() {
		var mp sharedtypes.MarketPair
		if err := rows.Scan(&mp.Timestamp, &mp.Currency1, &mp.Currency2, &mp.Open, &mp.High, &mp.Low, &mp.Close); err != nil {
			return nil, err
		}
		mps = append(mps, mp)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return mps, nil
}

func Create(db *sql.DB, pair sharedtypes.MarketPair) (int64, error) {
	stmt, err := db.Prepare("INSERT INTO markets(timestamp, `from`, `to`, open, high, low, close) VALUES(?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(pair.Timestamp, pair.Currency1, pair.Currency2, pair.Open, pair.High, pair.Low, pair.Close)
	if err != nil {
		return 0, err
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastInsertID, nil
}

func Truncate(db *sql.DB) error {
	_, err := db.Query("DELETE FROM markets;")
	if err != nil {
		return err
	}

	_, err = db.Query("VACUUM;")
	if err != nil {
		return err
	}

	_, err = db.Query("delete from sqlite_sequence where name='markets';")
	if err != nil {
		return err
	}

	return nil
}
