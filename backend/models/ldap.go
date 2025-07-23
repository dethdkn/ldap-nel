package models

import (
	"errors"

	"github.com/dethdkn/ldap-nel/backend/db"
	"github.com/dethdkn/ldap-nel/backend/utils"
)

type Ldap struct {
	ID       int64  `json:"id"`
	Name     string `json:"name" binding:"required"`
	URL      string `json:"url" binding:"required"`
	Port     int64  `json:"port" binding:"required"`
	SSL      bool   `json:"ssl"`
	BaseDN   string `json:"base_dn" binding:"required"`
	BindDN   string `json:"bind_dn" binding:"required"`
	BindPass string `json:"bind_pass" binding:"required"`
}

func (l *Ldap) Save() error {
	stmt, err := db.SQL.Prepare(`INSERT INTO ldaps (name, url, port, ssl, base_dn, bind_dn, bind_pass) VALUES (?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	l.BindPass, err = utils.Encrypt(l.BindPass)
	if err != nil {
		return err
	}

	result, err := stmt.Exec(l.URL, l.Port, l.SSL, l.BaseDN, l.BindDN, l.BindPass)
	if err != nil {
		return err
	}

	ldapID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	l.ID = ldapID
	return nil
}

func (l *Ldap) Update() error {
	stmt, err := db.SQL.Prepare(`UPDATE ldaps SET name = ?, url = ?, port = ?, ssl = ?, base_dn = ?, bind_dn = ?, bind_pass = ? WHERE id = ?`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	if l.BindPass != "" {
		l.BindPass, err = utils.Encrypt(l.BindPass)
		if err != nil {
			return err
		}
	} else {
		// If bind password is not set, we keep the existing one
		var existingBindPassword string
		row := db.SQL.QueryRow(`SELECT bind_pass FROM ldaps WHERE id = ?`, l.ID)
		if err := row.Scan(&existingBindPassword); err != nil {
			return err
		}
		l.BindPass = existingBindPassword
	}

	l.BindPass, err = utils.Encrypt(l.BindPass)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(l.Name, l.URL, l.Port, l.SSL, l.BaseDN, l.BindDN, l.BindPass, l.ID)
	return err
}

func (l *Ldap) Delete() error {
	stmt, err := db.SQL.Prepare(`DELETE FROM ldaps WHERE id = ?`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(l.ID)
	return err
}

func GetLdapByID(id int64, getBindPassword bool) (*Ldap, error) {
	row := db.SQL.QueryRow(`SELECT id, name, url, port, ssl, base_dn, bind_dn, bind_pass FROM ldaps WHERE id = ?`, id)

	var l Ldap
	if err := row.Scan(&l.ID, &l.Name, &l.URL, &l.Port, &l.SSL, &l.BaseDN, &l.BindDN, &l.BindPass); err != nil {
		return nil, errors.New("ldap not found")
	}

	if getBindPassword {
		var err error
		l.BindPass, err = utils.Decrypt(l.BindPass)
		if err != nil {
			return nil, err
		}
	} else {
		l.BindPass = ""
	}

	return &l, nil
}

func GetAllLdaps() ([]*Ldap, error) {
	row, err := db.SQL.Query(`SELECT id, name, url, port, ssl, base_dn, bind_dn FROM ldaps`)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	var ldaps []*Ldap
	for row.Next() {
		var l Ldap
		if err := row.Scan(&l.ID, &l.Name, &l.URL, &l.Port, &l.SSL, &l.BaseDN, &l.BindDN); err != nil {
			return nil, err
		}
		ldaps = append(ldaps, &l)
	}
	if err := row.Err(); err != nil {
		return nil, err
	}
	return ldaps, nil
}
