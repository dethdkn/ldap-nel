package models

import (
	"errors"

	"github.com/dethdkn/ldap-nel/backend/db"
	"github.com/dethdkn/ldap-nel/backend/passwords"
)

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Admin    bool   `json:"admin"`
}

func (u *User) Validate() error {
	row := db.SQL.QueryRow(`SELECT password, admin FROM users WHERE username = ?`, u.Username)

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
	stmt, err := db.SQL.Prepare(`INSERT INTO users (username, password, admin) VALUES (?, ?, ?)`)
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

func (u *User) Update() error {
	stmt, err := db.SQL.Prepare(`UPDATE users SET username = ?, password = ?, admin = ? WHERE id = ?`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	if u.Password != "" {
		u.Password, err = passwords.Sha512Crypt(u.Password, "")
		if err != nil {
			return err
		}
	} else {
		// If password is not set, we keep the existing one
		var existingPassword string
		row := db.SQL.QueryRow(`SELECT password FROM users WHERE id = ?`, u.ID)
		if err := row.Scan(&existingPassword); err != nil {
			return err
		}
		u.Password = existingPassword
	}

	_, err = stmt.Exec(u.Username, u.Password, u.Admin, u.ID)
	return err
}

func (u *User) UpdatePassword() error {
	stmt, err := db.SQL.Prepare(`UPDATE users SET password = ? WHERE username = ?`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	u.Password, err = passwords.Sha512Crypt(u.Password, "")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(u.Password, u.Username)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) Delete() error {
	stmt, err := db.SQL.Prepare(`DELETE FROM users WHERE id = ?`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(u.ID)
	return err
}

func GetUserByID(id int64) (*User, error) {
	row := db.SQL.QueryRow(`SELECT id, username, admin FROM users WHERE id = ?`, id)

	var user User
	if err := row.Scan(&user.ID, &user.Username, &user.Admin); err != nil {
		return nil, errors.New("user not found")
	}

	return &user, nil
}

func GetAllUsers() ([]*User, error) {
	rows, err := db.SQL.Query(`SELECT id, username, admin FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Username, &user.Admin); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func IsUsersEmpty() (bool, error) {
	row := db.SQL.QueryRow(`SELECT COUNT(*) FROM users`)
	var count int
	if err := row.Scan(&count); err != nil {
		return false, err
	}
	return count == 0, nil
}

func AlreadyExists(username string) (bool, error) {
	row := db.SQL.QueryRow(`SELECT COUNT(*) FROM users WHERE username = ?`, username)
	var count int
	if err := row.Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}
