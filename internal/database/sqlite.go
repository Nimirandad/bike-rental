package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type SQLiteClient struct {
	DB   *sql.DB
	path string
}

func Connect(dbPath string) (*SQLiteClient, error) {
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create directory for SQLite database: %v", err)
	}

	database, err := sql.Open("sqlite", dbPath+"?parseTime=true")
	if err != nil {
		return nil, fmt.Errorf("error while establishing new connection to DB: %v", err)
	}

	if err := database.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging database: %v", err)
	}

	return &SQLiteClient{
		DB:   database,
		path: dbPath,
	}, nil
}

func (conn *SQLiteClient) Close() error {
	if err := conn.DB.Close(); err != nil {
		return fmt.Errorf("error while closing the DB connection: %v", err)
	}
	return nil
}