package db

import (
	"database/sql"

	"github.com/dethdkn/ldap-nel/backend/passwords"
	_ "github.com/mattn/go-sqlite3"
)

var SQL *sql.DB

func InitDB() (*sql.DB, error) {
	SQL, err := sql.Open("sqlite3", "db.sqlite")
	if err != nil {
		return nil, err
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL,
		password TEXT NOT NULL
	);`

	_, err = SQL.Exec(createTableQuery)
	if err != nil {
		return nil, err
	}

	// Create a default user if the table is empty
	row := SQL.QueryRow("SELECT COUNT(*) FROM users")
	var count int
	if err := row.Scan(&count); err != nil {
		return nil, err
	}
	if count == 0 {
		_, err := SQL.Exec("INSERT INTO users (username, password) VALUES (?, ?)", "nel", passwords.Sha512Crypt("test123", ""))
		if err != nil {
			return nil, err
		}
	}

	return SQL, nil
}
