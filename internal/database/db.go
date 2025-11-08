package database

import (
	"database/sql"
	"fmt"
	"todo-api/internal/config"

	_ "github.com/lib/pq"
)

type DB struct {
	conn *sql.DB
}

func NewConnection(cfg config.DatabaseConfig) (*DB, error) {

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName,
	)

	conn, err := sql.Open("postgres", dsn)

	if err != nil {
		return nil, err
	}

	if err := conn.Ping(); err != nil {
		return nil, err
	}

	return &DB{conn: conn}, nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func (db *DB) Ping() error {
	return db.conn.Ping()
}
