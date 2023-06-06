package rates

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteDB struct {
	sqldb  *sql.DB
	get    *sql.Stmt
	insert *sql.Stmt
}

func OpenDB(fpath string) (*SQLiteDB, error) {
	sqldb, err := sql.Open("sqlite3", fpath+"?_busy_timeout=10000&_journal=WAL&_sync=NORMAL&cache=shared")
	if err != nil {
		return nil, fmt.Errorf("opening database %s: %v", fpath, err)
	}

	if _, err := sqldb.Exec(`
		create table if not exists rates (
			date     text not null,
			currency text not null,
			rate     real not null,
			primary key (date, currency)
		);
	`); err != nil {
		return nil, err
	}

	get, err := sqldb.Prepare("select currency, rate from rates where date = ?")
	if err != nil {
		return nil, err
	}
	insert, err := sqldb.Prepare("insert or ignore into rates (date, currency, rate) values (?, ?, ?)") // ignore in order to keep (date, currency) constant
	if err != nil {
		return nil, err
	}

	return &SQLiteDB{
		sqldb:  sqldb,
		get:    get,
		insert: insert,
	}, nil
}

func (db *SQLiteDB) Get(date string) (map[string]float64, error) {
	rows, err := db.get.Query(date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rates = make(map[string]float64)
	for rows.Next() {
		var currency string
		var rate float64
		if err := rows.Scan(&currency, &rate); err != nil {
			return nil, err
		}
		rates[currency] = rate
	}
	return rates, nil
}

func (db *SQLiteDB) Insert(date string, rates map[string]float64) error {
	tx, err := db.sqldb.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for currency, rate := range rates {
		if _, err := db.insert.Exec(date, currency, rate); err != nil {
			return err
		}
	}
	return tx.Commit()
}
