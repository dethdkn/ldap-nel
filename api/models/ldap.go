package models

import (
	"errors"

	"github.com/dethdkn/ldap-nel/api/db"
	"github.com/dethdkn/ldap-nel/api/ldap"
	"github.com/dethdkn/ldap-nel/api/passwords"
	"github.com/dethdkn/ldap-nel/api/utils"
)

type Ldap struct {
	ID       int64  `json:"id"`
	Name     string `json:"name" binding:"required"`
	URL      string `json:"url" binding:"required"`
	Port     int64  `json:"port" binding:"required"`
	SSL      bool   `json:"ssl"`
	BaseDN   string `json:"base_dn" binding:"required"`
	BindDN   string `json:"bind_dn"`
	BindPass string `json:"bind_pass"`
}

type LdapName struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (l *Ldap) Save() error {
	err := ldap.TestLdap(l.URL, l.Port, l.SSL, l.BaseDN, l.BindDN, l.BindPass)
	if err != nil {
		return err
	}

	stmt, err := db.SQL.Prepare(`INSERT INTO ldaps (name, url, port, ssl, base_dn, bind_dn, bind_pass) VALUES (?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	if l.BindPass != "" {
		l.BindPass, err = utils.Encrypt(l.BindPass)
		if err != nil {
			return err
		}
	}

	result, err := stmt.Exec(l.Name, l.URL, l.Port, l.SSL, l.BaseDN, l.BindDN, l.BindPass)
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

	decryptedPass, err := utils.Decrypt(l.BindPass)
	if err != nil {
		return err
	}

	err = ldap.TestLdap(l.URL, l.Port, l.SSL, l.BaseDN, l.BindDN, decryptedPass)
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

func GetAllLdapsNames() ([]LdapName, error) {
	row, err := db.SQL.Query(`SELECT id, name FROM ldaps`)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	var ldaps []LdapName
	for row.Next() {
		var l LdapName
		if err := row.Scan(&l.ID, &l.Name); err != nil {
			return nil, err
		}
		ldaps = append(ldaps, l)
	}
	if err := row.Err(); err != nil {
		return nil, err
	}
	return ldaps, nil
}

func GetLdapChilds(id int64, dn string) (string, []string, error) {
	l, err := GetLdapByID(id, true)
	if err != nil {
		return "", nil, err
	}

	if dn == "" {
		dn = l.BaseDN
	}

	childDNs, err := ldap.SearchChilds(l.URL, l.Port, l.SSL, dn, l.BindDN, l.BindPass)
	if err != nil {
		return "", nil, err
	}

	return dn, childDNs, nil
}

func GetLdapAttributes(id int64, dn string) (map[string][]string, error) {
	l, err := GetLdapByID(id, true)
	if err != nil {
		return nil, err
	}

	if dn == "" {
		dn = l.BaseDN
	}

	attrs, err := ldap.SearchAttributes(l.URL, l.Port, l.SSL, dn, l.BindDN, l.BindPass)
	if err != nil {
		return nil, err
	}

	return attrs, nil
}

func GetLdapPossibleAttributes(id int64, dn string) ([]string, error) {
	l, err := GetLdapByID(id, true)
	if err != nil {
		return nil, err
	}

	if dn == "" {
		dn = l.BaseDN
	}

	attrs, err := ldap.GetPossibleAttributes(l.URL, l.Port, l.SSL, l.BindDN, l.BindPass, dn)
	if err != nil {
		return nil, err
	}

	return attrs, nil
}

func AddLdapAttributeValue(id int64, dn, attribute, value string) error {
	l, err := GetLdapByID(id, true)
	if err != nil {
		return err
	}

	if dn == "" || attribute == "" || value == "" {
		return errors.New("dn, attribute, and value are required")
	}

	if attribute == "userPassword" {
		value, err = passwords.GetEncryptedPass(value)
		if err != nil {
			return err
		}
	}

	err = ldap.AddAttributeValue(l.URL, l.Port, l.SSL, l.BindDN, l.BindPass, dn, attribute, value)
	if err != nil {
		return err
	}

	return nil
}

func UpdateLdapAttributeValue(id int64, dn, attribute, value, newValue string) error {
	l, err := GetLdapByID(id, true)
	if err != nil {
		return err
	}

	if dn == "" || attribute == "" || value == "" || newValue == "" {
		return errors.New("dn, attribute, value, and newValue are required")
	}

	if attribute == "userPassword" {
		newValue, err = passwords.GetEncryptedPass(newValue)
		if err != nil {
			return err
		}
	}

	err = ldap.UpdateAttributeValue(l.URL, l.Port, l.SSL, l.BindDN, l.BindPass, dn, attribute, value, newValue)
	if err != nil {
		return err
	}

	return nil
}

func DeleteLdapAttributeValue(id int64, dn, attribute, value string) error {
	l, err := GetLdapByID(id, true)
	if err != nil {
		return err
	}

	if dn == "" || attribute == "" || value == "" {
		return errors.New("dn, attribute, and value are required")
	}

	err = ldap.DeleteAttributeValue(l.URL, l.Port, l.SSL, l.BindDN, l.BindPass, dn, attribute, value)
	if err != nil {
		return err
	}

	return nil
}

func ExportLdap(id int64, dn string) (string, error) {
	l, err := GetLdapByID(id, true)
	if err != nil {
		return "", err
	}

	if dn == "" {
		dn = l.BaseDN
	}

	ldif, err := ldap.ExportLdap(l.URL, l.Port, l.SSL, l.BindDN, l.BindPass, dn)
	if err != nil {
		return "", err
	}

	return ldif, nil
}

func ImportLdap(id int64, fileData []byte) error {
	l, err := GetLdapByID(id, true)
	if err != nil {
		return err
	}

	err = ldap.ImportLdap(l.URL, l.Port, l.SSL, l.BindDN, l.BindPass, fileData)
	if err != nil {
		return err
	}

	return nil
}

func AddLdapDn(id int64, dn string, attributes []ldap.Attribute) error {
	l, err := GetLdapByID(id, true)
	if err != nil {
		return err
	}

	for val, attr := range attributes {
		if attr.Attribute == "userPassword" {
			attr.Value, err = passwords.GetEncryptedPass(attr.Value)
			if err != nil {
				return err
			}
			attributes[val] = attr
		}
	}

	err = ldap.AddDn(l.URL, l.Port, l.SSL, l.BindDN, l.BindPass, dn, attributes)
	if err != nil {
		return err
	}

	return nil
}

func DeleteLdapDn(id int64, dn string, smart bool) error {
	l, err := GetLdapByID(id, true)
	if err != nil {
		return err
	}

	if dn == "" {
		return errors.New("dn is required")
	}

	err = ldap.DeleteDn(l.URL, l.Port, l.SSL, l.BindDN, l.BindPass, dn, smart)
	if err != nil {
		return err
	}

	return nil
}

func CopyLdapDn(id int64, dn, targetDn string) error {
	l, err := GetLdapByID(id, true)
	if err != nil {
		return err
	}

	if dn == "" || targetDn == "" {
		return errors.New("dn and targetDn are required")
	}

	err = ldap.CopyDn(l.URL, l.Port, l.SSL, l.BindDN, l.BindPass, dn, targetDn)
	if err != nil {
		return err
	}

	return nil
}

func MoveLdapDn(id int64, dn, targetDn string) error {
	l, err := GetLdapByID(id, true)
	if err != nil {
		return err
	}

	if dn == "" || targetDn == "" {
		return errors.New("dn and targetDn are required")
	}

	err = ldap.MoveDn(l.URL, l.Port, l.SSL, l.BindDN, l.BindPass, dn, targetDn)
	if err != nil {
		return err
	}

	return nil
}
