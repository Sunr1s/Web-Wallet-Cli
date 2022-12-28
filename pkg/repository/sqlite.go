package repository

import (
	"database/sql"
	"log"
)

const (
	usersTable = "users"
	txTable    = "transactions"
)

type Config struct {
	Host   string
	DBName string
}

func NewSqlite3DB(cfg Config) (*sql.DB, error) {
	db, err := sql.Open(cfg.Host, cfg.DBName)
	if err != nil {
		log.Fatalf("error while open db")
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("db dont answer")
		return nil, err
	}

	return db, nil
}
