package db

import (
	"database/sql"
	_ "github.com/lib/pq"
)

func OpenDB(DSN string) (*sql.DB, error) {
	db, err := sql.Open("postgres", DSN)
	if err != nil {
		return nil, err
	}

	return db, nil
}
