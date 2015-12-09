package main

import (
	"database/sql"
	_ "github.com/lib/pq"
)

var dbInstance *sql.DB = nil

func GetDB() (*sql.DB, error) {
	var err error
	if dbInstance == nil {
		dbInstance, err = sql.Open("postgres", "postgres://localhost/nanocloud?sslmode=disable")
	}
	return dbInstance, err
}
