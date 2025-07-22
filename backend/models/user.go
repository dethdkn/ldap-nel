package models

import (
	"errors"

	"github.com/dethdkn/ldap-nel/backend/db"
	"github.com/dethdkn/ldap-nel/backend/passwords"
)

type User struct {
	ID       int64
	Username string `binding:"required"`
	Password string `binding:"required"`
	Admin    bool
}

func (u *User) Validate() error {
	query := `SELECT password, admin FROM users WHERE username = ?`
	row := db.SQL.QueryRow(query, u.Username)

	var storedPassword string
	var admin bool
	if err := row.Scan(&storedPassword, &admin); err != nil {
		return errors.New("invalid username or password")
	}

	u.Admin = admin

	passwordIsValid := passwords.VerifySha512Crypt(u.Password, storedPassword)

	if !passwordIsValid {
		return errors.New("invalid username or password")
	}

	return nil
}

func (u *User) Save() error {
	query := `INSERT INTO users (username, password, admin) VALUES (?, ?, ?)`
	stmt, err := db.SQL.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	u.Password, err = passwords.Sha512Crypt(u.Password, "")
	if err != nil {
		return err
	}

	result, err := stmt.Exec(u.Username, u.Password, u.Admin)
	if err != nil {
		return err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	u.ID = userID
	return nil
}
