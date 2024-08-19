package sqlite

import (
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sqlx.DB
}

func Dial(dsn string) (*DB, error) {
	// dsn example: file:test.db?cache=shared&mode=memory
	db, err := sqlx.Connect("sqlite", dsn)
	if err != nil {
		return nil, errors.New("error connecting to db")
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(1 * time.Minute)

	return &DB{DB: db}, nil
}
