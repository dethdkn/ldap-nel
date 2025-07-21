package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() (*sql.DB, error) {
	DB, err := sql.Open("sqlite3", "db.sqlite")
	if err != nil {
		return nil, err
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL,
		password TEXT NOT NULL
	);`

	_, err = DB.Exec(createTableQuery)
	if err != nil {
		return nil, err
	}

	// Create a default user if the table is empty
	row := DB.QueryRow("SELECT COUNT(*) FROM users")
	var count int
	if err := row.Scan(&count); err != nil {
		return nil, err
	}
	if count == 0 {
		_, err := DB.Exec("INSERT INTO users (username, password) VALUES (?, ?)", "nel", "123")
		if err != nil {
			return nil, err
		}
	}

	return DB, nil
}
