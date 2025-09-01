package ldap

import (
	"bytes"
	"encoding/base64"
	"errors"

	"github.com/go-ldap/ldap/v3"
)

func DeleteAttributeValue(url string, port int64, ssl bool, bindDN, bindPass, dn, attribute, value string) error {
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

	if binaryAttrs[attribute] {
		target, err := base64.StdEncoding.DecodeString(value)
		if err != nil {
			return err
		}

		attrs, err := SearchAttributes(url, port, ssl, dn, bindDN, bindPass)
		if err != nil {
			return err
		}

		existingB64, ok := attrs[attribute]
		if !ok || len(existingB64) == 0 {
			return errors.New("attribute not present")
		}

		remaining := make([]string, 0, len(existingB64))
		removed := false
		for _, b64 := range existingB64 {
			b, decErr := base64.StdEncoding.DecodeString(b64)
			if decErr != nil {
				return decErr
			}
			if !bytes.Equal(b, target) {
				remaining = append(remaining, string(b))
			} else {
				removed = true
			}
		}

		if !removed {
			return errors.New("value not found")
		}

		req := ldap.NewModifyRequest(dn, nil)
		if len(remaining) == 0 {
			req.Delete(attribute, nil)
		} else {
			req.Replace(attribute, remaining)
		}
		return l.Modify(req)
	}

	modifyReq := ldap.NewModifyRequest(dn, nil)
	modifyReq.Delete(attribute, []string{value})

	if err = l.Modify(modifyReq); err != nil {
		return err
	}
	return nil
}

func DeleteDn(url string, port int64, ssl bool, bindDN, bindPass, dn string, smart bool) error {
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

	if smart {
		searchReq := ldap.NewSearchRequest(
			dn, ldap.ScopeSingleLevel, ldap.NeverDerefAliases, 0, 0, false,
			"(objectClass=*)",
			[]string{"dn"},
			nil,
		)

		searchRes, err := l.Search(searchReq)
		if err != nil {
			return err
		}

		for _, entry := range searchRes.Entries {
			if err := DeleteDn(url, port, ssl, bindDN, bindPass, entry.DN, smart); err != nil {
				return err
			}
		}
	}

	delReq := ldap.NewDelRequest(dn, nil)
	if err = l.Del(delReq); err != nil {
		return err
	}

	return nil
}
