package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func Open(driver string, dsn string) (*sql.DB, error) {
	db, err := sql.Open(driver, dsn)

	return db, err
}
