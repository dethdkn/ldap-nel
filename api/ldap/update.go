package ldap

import (
	"errors"

	"github.com/go-ldap/ldap/v3"
)

func UpdateAttributeValue(url string, port int64, ssl bool, bindDN, bindPass, dn, attribute, value, newValue string) error {
	l, err := Connect(url, port, ssl)
	if err != nil {
		return err
	}

	defer l.Unbind()

	if bindDN != "" && bindPass != "" {
		if err = l.Bind(bindDN, bindPass); err != nil {
			return errors.New("failed to bind with provided credentials")
		}
	}

	modifyReq := ldap.NewModifyRequest(dn, nil)
	modifyReq.Delete(attribute, []string{value})
	modifyReq.Add(attribute, []string{newValue})

	if err = l.Modify(modifyReq); err != nil {
		return errors.New("failed to update attribute value")
	}

	return nil
}
