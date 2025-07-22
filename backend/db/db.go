package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var SQL *sql.DB

func InitDB() (*sql.DB, error) {
	SQL, err := sql.Open("sqlite3", "/database/db.sqlite")
	if err != nil {
		return nil, err
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		admin BOOLEAN DEFAULT FALSE
	);`

	_, err = SQL.Exec(createTableQuery)
	if err != nil {
		return nil, err
	}

	return SQL, nil
}
