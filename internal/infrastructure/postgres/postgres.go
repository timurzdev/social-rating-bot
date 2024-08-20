package postgres

import (
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DB struct {
	*sqlx.DB
}

func Dial(dsn string) (*DB, error) {
	fmt.Println(dsn)
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, errors.New("error connecting to db")
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(1 * time.Minute)

	return &DB{DB: db}, nil
}
