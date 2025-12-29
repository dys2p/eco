package rates

import (
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteDB struct {
	sqldb  *sql.DB
	get    *sql.Stmt
	insert *sql.Stmt
	latest *sql.Stmt
}

func OpenDB(fpath string) (*SQLiteDB, error) {
	sqldb, err := sql.Open("sqlite3", fpath+"?_busy_timeout=10000&_journal=WAL&_sync=NORMAL&cache=shared")
	if err != nil {
		return nil, fmt.Errorf("opening database %s: %v", fpath, err)
	}

	if _, err := sqldb.Exec(`
		create table if not exists rates_history (
			date  text primary key,
			rates text not null -- json map
		);
		create index if not exists date_index on rates_history (date);
	`); err != nil {
		return nil, err
	}

	get, err := sqldb.Prepare("select rates from rates_history where date = ?")
	if err != nil {
		return nil, err
	}
	insert, err := sqldb.Prepare("insert or ignore into rates_history (date, rates) values (?, ?)") // ignore existing, don't modify them
	if err != nil {
		return nil, err
	}
	latest, err := sqldb.Prepare("select ifnull(max(date), '0000-00-00') from rates_history where date <= ?")
	if err != nil {
		return nil, err
	}

	return &SQLiteDB{
		sqldb:  sqldb,
		get:    get,
		insert: insert,
		latest: latest,
	}, nil
}

// Get returns ErrNoRows if no data is found.
func (db *SQLiteDB) Get(date string) (map[string]float64, error) {
	var encoded []byte
	if err := db.get.QueryRow(date).Scan(&encoded); err != nil {
		return nil, err
	}
	var rs = make(map[string]float64)
	if err := json.Unmarshal(encoded, &rs); err != nil {
		return nil, err
	}
	return rs, nil
}

func (db *SQLiteDB) Insert(date string, m map[string]float64) error {
	encoded, err := json.Marshal(m)
	if err != nil {
		return err
	}
	_, err = db.insert.Exec(date, encoded)
	return err
}

func (db *SQLiteDB) LatestDate(maxDate string) (string, error) {
	var latest string
	return latest, db.latest.QueryRow(maxDate).Scan(&latest)
}
