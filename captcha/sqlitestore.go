package captcha

import (
	"database/sql"
	"log"
	"time"

	"github.com/dchest/captcha"
)

func Initialize(sqlitedb string) error {
	db, err := sql.Open("sqlite3", sqlitedb+"?_busy_timeout=10000&_journal=WAL&_sync=NORMAL&cache=shared")
	if err != nil {
		return err
	}
	store, err := newSQLiteStore(db)
	if err != nil {
		return err
	}
	captcha.SetCustomStore(store)

	go func() {
		var ticker = time.NewTicker(1 * time.Hour)
		for range ticker.C {
			store.Cleanup()
		}
	}()

	return nil
}

type sqliteStore struct {
	sqlDB *sql.DB
	clean *sql.Stmt
	del   *sql.Stmt
	get   *sql.Stmt
	set   *sql.Stmt
}

func newSQLiteStore(db *sql.DB) (*sqliteStore, error) {
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS captcha (
			id     TEXT    PRIMARY KEY,
			time   INTEGER NOT NULL,
			digits BLOB    NOT NULL
		)`); err != nil {
		return nil, err
	}

	clean, err := db.Prepare("DELETE FROM captcha WHERE time < ?")
	if err != nil {
		return nil, err
	}
	del, err := db.Prepare("DELETE FROM captcha WHERE id = ?")
	if err != nil {
		return nil, err
	}
	get, err := db.Prepare("SELECT digits FROM captcha WHERE id = ? LIMIT 1")
	if err != nil {
		return nil, err
	}
	set, err := db.Prepare("INSERT OR REPLACE INTO captcha (id, time, digits) VALUES (?, ?, ?)")
	if err != nil {
		return nil, err
	}

	return &sqliteStore{
		sqlDB: db,
		clean: clean,
		del:   del,
		get:   get,
		set:   set,
	}, nil
}

// Cleanup removes captchas older than two days.
func (s *sqliteStore) Cleanup() {
	if _, err := s.clean.Exec(time.Now().AddDate(0, 0, -2).Unix()); err != nil {
		log.Printf("error cleaning up captcha sqlite store: %v", err)
	}
}

// Get returns stored digits for the captcha id. Clear indicates
// whether the captcha must be deleted from the store.
func (s *sqliteStore) Get(id string, clear bool) (digits []byte) {
	if err := s.get.QueryRow(id).Scan(&digits); err != nil {
		if err == sql.ErrNoRows {
			// no problem, happens if a POST request is repeated
		} else {
			log.Printf("error getting from captcha sqlite store: %v", err)
		}
	}
	if clear {
		if _, err := s.del.Exec(id); err != nil {
			log.Printf("error deleting from captcha sqlite store: %v", err)
		}
	}
	return
}

// Set sets the digits for the captcha id.
func (s *sqliteStore) Set(id string, digits []byte) {
	if _, err := s.set.Exec(id, time.Now().Unix(), digits); err != nil {
		log.Printf("error setting captcha sqlite store: %v", err)
	}
}
