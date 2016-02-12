package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"os"
)

var _db *sql.DB

func getInstance() (*sql.DB, error) {
	if _db == nil {

		databaseURI := os.Getenv("DATABASE_URI")
		if len(databaseURI) == 0 {
			databaseURI = "postgres://localhost/nanocloud?sslmode=disable"
		}

		var err error
		_db, err = sql.Open("postgres", databaseURI)
		return _db, err
	}
	return _db, nil
}

func Query(query string, args ...interface{}) (*sql.Rows, error) {
	db, err := getInstance()
	if err != nil {
		return nil, err
	}
	return db.Query(query, args...)
}

func Exec(query string, args ...interface{}) (sql.Result, error) {
	db, err := getInstance()
	if err != nil {
		return nil, err
	}
	return db.Exec(query, args...)
}
