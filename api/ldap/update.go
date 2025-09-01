package ldap

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/go-ldap/ldap/v3"
)

type Attribute struct {
	Attribute string `json:"attribute" binding:"required"`
	Value     string `json:"value" binding:"required"`
}

func AddAttributeValue(url string, port int64, ssl bool, bindDN, bindPass, dn, attribute, value string) error {
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
		data, err := base64.StdEncoding.DecodeString(value)
		if err != nil {
			return err
		}

		addReq := ldap.NewModifyRequest(dn, nil)
		addReq.Add(attribute, []string{string(data)})
		if err := l.Modify(addReq); err != nil {
			if le, ok := err.(*ldap.Error); ok && le.ResultCode == ldap.LDAPResultInappropriateMatching {
				attrs, err := SearchAttributes(url, port, ssl, dn, bindDN, bindPass)
				if err != nil {
					return err
				}

				var all []string
				if existingB64, ok := attrs[attribute]; ok {
					for _, b64 := range existingB64 {
						b, decErr := base64.StdEncoding.DecodeString(b64)
						if decErr != nil {
							return decErr
						}
						all = append(all, string(b))
					}
				}
				all = append(all, string(data))

				repReq := ldap.NewModifyRequest(dn, nil)
				repReq.Replace(attribute, all)
				return l.Modify(repReq)
			}
			return err
		}
		return nil
	}

	modifyReq := ldap.NewModifyRequest(dn, nil)
	modifyReq.Add(attribute, []string{value})
	if err = l.Modify(modifyReq); err != nil {
		return err
	}

	return nil
}

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

	if binaryAttrs[attribute] {
		oldBytes, err := base64.StdEncoding.DecodeString(value)
		if err != nil {
			return err
		}
		newBytes, err := base64.StdEncoding.DecodeString(newValue)
		if err != nil {
			return err
		}

		attrs, err := SearchAttributes(url, port, ssl, dn, bindDN, bindPass)
		if err != nil {
			return err
		}

		var combined []string
		if existingB64, ok := attrs[attribute]; ok {
			for _, b64 := range existingB64 {
				b, decErr := base64.StdEncoding.DecodeString(b64)
				if decErr != nil {
					return decErr
				}
				if string(b) != string(oldBytes) {
					combined = append(combined, string(b))
				}
			}
		}
		combined = append(combined, string(newBytes))

		req := ldap.NewModifyRequest(dn, nil)
		req.Replace(attribute, combined)
		return l.Modify(req)
	}

	modifyReq := ldap.NewModifyRequest(dn, nil)
	modifyReq.Delete(attribute, []string{value})
	modifyReq.Add(attribute, []string{newValue})

	if err = l.Modify(modifyReq); err != nil {
		return err
	}
	return nil
}

func ImportLdap(url string, port int64, ssl bool, bindDN, bindPass string, fileData []byte) error {
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

	lines := unfold(fileData)
	blocks := splitBlocks(lines)

	for i, block := range blocks {
		dn, attrs, err := parseAddBlock(block)
		if err != nil {
			return fmt.Errorf("ldif entry #%d: %w", i+1, err)
		}
		if dn == "" {
			return fmt.Errorf("ldif entry #%d: missing dn", i+1)
		}

		req := ldap.NewAddRequest(dn, nil)
		for name, values := range attrs {
			if strings.EqualFold(name, "dn") {
				continue
			}
			req.Attribute(name, values)
		}
		if err := l.Add(req); err != nil {
			return fmt.Errorf("add %q failed: %w", dn, err)
		}
	}
	return nil
}

func AddDn(url string, port int64, ssl bool, bindDN, bindPass, dn string, attributes []Attribute) error {
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

	req := ldap.NewAddRequest(dn, nil)

	grouped := make(map[string][]string)

	for _, attr := range attributes {
		if binaryAttrs[attr.Attribute] {
			data, err := base64.StdEncoding.DecodeString(attr.Value)
			if err != nil {
				return err
			}
			grouped[attr.Attribute] = append(grouped[attr.Attribute], string(data))
		} else {
			grouped[attr.Attribute] = append(grouped[attr.Attribute], attr.Value)
		}
	}

	for name, values := range grouped {
		req.Attribute(name, values)
	}

	if err := l.Add(req); err != nil {
		return err
	}

	return nil
}
