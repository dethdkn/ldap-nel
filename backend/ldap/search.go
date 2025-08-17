package ldap

import (
	"errors"

	"github.com/go-ldap/ldap/v3"
)

func SearchChilds(url string, port int64, ssl bool, DN, bindDN, bindPass string) ([]string, error) {
	l, err := Connect(url, port, ssl)
	if err != nil {
		return nil, err
	}

	defer l.Unbind()

	if bindDN != "" && bindPass != "" {
		if err = l.Bind(bindDN, bindPass); err != nil {
			return nil, errors.New("failed to bind with provided credentials")
		}
	}

	searchReq := ldap.NewSearchRequest(
		DN,
		ldap.ScopeChildren, ldap.NeverDerefAliases, 0, 0, false,
		"(objectClass=*)",
		[]string{"dn"},
		nil,
	)

	result, err := l.Search(searchReq)
	if err != nil {
		return nil, errors.New("failed to search for child entries")
	}

	if len(result.Entries) == 0 {
		return nil, nil
	}

	var childDNs []string
	for _, entry := range result.Entries {
		childDNs = append(childDNs, entry.DN)
	}

	return childDNs, nil
}

func SearchAttributes(url string, port int64, ssl bool, DN, bindDN, bindPass string) (map[string][]string, error) {
	l, err := Connect(url, port, ssl)
	if err != nil {
		return nil, err
	}

	defer l.Unbind()

	if bindDN != "" && bindPass != "" {
		if err = l.Bind(bindDN, bindPass); err != nil {
			return nil, errors.New("failed to bind with provided credentials")
		}
	}

	searchReq := ldap.NewSearchRequest(
		DN,
		ldap.ScopeBaseObject, ldap.NeverDerefAliases, 0, 0, false,
		"(objectClass=*)",
		nil,
		nil,
	)

	result, err := l.Search(searchReq)
	if err != nil {
		return nil, errors.New("failed to search for attributes")
	}

	if len(result.Entries) == 0 {
		return nil, nil
	}

	attrs := make(map[string][]string)
	for _, attr := range result.Entries[0].Attributes {
		attrs[attr.Name] = attr.Values
	}

	return attrs, nil
}
