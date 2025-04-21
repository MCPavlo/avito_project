// internal/db/db.go
package db

import (
	"database/sql"
	"fmt"
	_ "log"
	"tz/internal/config"

	_ "github.com/lib/pq" // драйвер для PostgreSQL
)

// DB представляет собой структуру для работы с базой данных
type DB struct {
	*sql.DB
}

// NewDB создает новое подключение к базе данных
func NewDB(cfg *config.Config) (*DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.DBName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

// Get выполняет запрос и возвращает одну запись
func (db *DB) Get(dest interface{}, query string, args ...interface{}) error {
	return db.QueryRow(query, args...).Scan(dest)
}
