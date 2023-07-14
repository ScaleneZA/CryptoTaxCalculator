package markets

import (
	"database/sql"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/conversionrate/sharedtypes"
)

func Create(db *sql.DB, pair sharedtypes.MarketPair) (int64, error) {
	stmt, err := db.Prepare("INSERT INTO markets(timestamp, `from`, `to`, open, high, low, close) VALUES(strftime('%s', 'now'), ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(pair.Currency1, pair.Currency2, pair.Open, pair.High, pair.Low, pair.Close)
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
	stmt, err := db.Prepare("DELETE FROM markets;")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	stmt, err = db.Prepare("VACUUM;")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	stmt, err = db.Prepare("delete from sqlite_sequence where name='markets';")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	return nil
}
