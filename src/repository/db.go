package repository

import (
	"github.com/jmoiron/sqlx"
	_"github.com/lib/pq"
)

func NewDB(connectionUrl string) (*sqlx.DB, error) {
	conn, err := sqlx.Open("postgres", connectionUrl)
	if err != nil {
		return nil, err
	}
	if err := conn.Ping(); err != nil {
		return nil, err
	}
	return conn, err
}
