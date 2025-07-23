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

	_, err = SQL.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			admin BOOLEAN DEFAULT FALSE
		);
	`)
	if err != nil {
		return nil, err
	}

	_, err = SQL.Exec(`
		CREATE TABLE IF NOT EXISTS ldaps (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			url TEXT NOT NULL,
			port INTEGER NOT NULL,
			ssl BOOLEAN DEFAULT FALSE,
			base_dn TEXT NOT NULL,
			bind_dn TEXT NOT NULL,
			bind_pass TEXT NOT NULL
		);
	`)
	if err != nil {
		return nil, err
	}

	return SQL, nil
}
