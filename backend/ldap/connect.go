package ldap

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/go-ldap/ldap/v3"
)

func Connect(url string, port int64, ssl bool) (*ldap.Conn, error) {
	scheme := "ldap"
	if ssl {
		scheme = "ldaps"
	}
	addr := fmt.Sprintf("%s://%s:%d", scheme, url, port)

	dialer := &net.Dialer{Timeout: 10 * time.Second}

	l, err := ldap.DialURL(addr, ldap.DialWithDialer(dialer))
	if err != nil {
		return nil, errors.New("failed to connect to LDAP server")
	}

	return l, nil
}

func TestLdap(url string, port int64, ssl bool, baseDN, bindDN, bindPass string) error {
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

	searchReq := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeBaseObject, ldap.NeverDerefAliases, 0, 0, false,
		"(objectClass=*)",
		[]string{"dn"},
		nil,
	)

	if _, err = l.Search(searchReq); err != nil {
		return errors.New("failed to search baseDN")
	}

	return nil
}
