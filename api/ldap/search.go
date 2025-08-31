package ldap

import (
	"encoding/base64"
	"errors"
	"regexp"
	"sort"
	"strings"

	"github.com/dethdkn/ldap-nel/api/utils"
	"github.com/go-ldap/ldap/v3"
)

func SearchChilds(url string, port int64, ssl bool, dn, bindDN, bindPass string) ([]string, error) {
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
		dn,
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

func SearchAttributes(url string, port int64, ssl bool, dn, bindDN, bindPass string) (map[string][]string, error) {
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
		dn,
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
		if binaryAttrs[attr.Name] && len(attr.ByteValues) > 0 {
			b64Values := make([]string, len(attr.ByteValues))
			for i, b := range attr.ByteValues {
				b64Values[i] = base64.StdEncoding.EncodeToString(b)
			}
			attrs[attr.Name] = b64Values
		} else {
			attrs[attr.Name] = attr.Values
		}
	}

	return attrs, nil
}

func GetPossibleAttributes(url string, port int64, ssl bool, bindDN, bindPass, dn string) ([]string, error) {
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

	entryReq := ldap.NewSearchRequest(
		dn,
		ldap.ScopeBaseObject, ldap.NeverDerefAliases, 0, 0, false,
		"(objectClass=*)",
		[]string{"objectClass", "subschemaSubentry", "+", "*"},
		nil,
	)
	entryRes, err := l.Search(entryReq)
	if err != nil {
		return nil, errors.New("failed to read entry")
	}
	if len(entryRes.Entries) == 0 {
		return nil, errors.New("entry not found")
	}
	entry := entryRes.Entries[0]

	objClasses := entry.GetAttributeValues("objectClass")
	for i := range objClasses {
		objClasses[i] = strings.ToLower(objClasses[i])
	}

	subschemaDN := entry.GetAttributeValue("subschemaSubentry")
	if subschemaDN == "" {
		subschemaDN = "cn=subschema"
	}

	schemaReq := ldap.NewSearchRequest(
		subschemaDN,
		ldap.ScopeBaseObject, ldap.NeverDerefAliases, 0, 0, false,
		"(objectClass=subschema)",
		[]string{"objectClasses"},
		nil,
	)
	schemaRes, err := l.Search(schemaReq)
	if err != nil || len(schemaRes.Entries) == 0 {
		return nil, errors.New("subschema entry not found or unreadable")
	}
	defs := schemaRes.Entries[0].GetAttributeValues("objectClasses")

	type oc struct{ must, may, sup []string }
	ocMap := map[string]oc{}

	reQuoted := regexp.MustCompile(`'([^']*)'`)
	nextKW := func(s string, pos int) int {
		keywords := []string{
			" NAME ", " DESC ", " OBSOLETE ", " SUP ", " ABSTRACT ",
			" STRUCTURAL ", " AUXILIARY ", " MUST ", " MAY ", " EQUALITY ",
			" ORDERING ", " SUBSTR ", " SYNTAX ", " SINGLE-VALUE ",
			" NO-USER-MODIFICATION ", " USAGE ", " X-",
			" )",
		}
		min := -1
		for _, k := range keywords {
			i := strings.Index(s[pos:], k)
			if i >= 0 {
				i = pos + i
				if min == -1 || i < min {
					min = i
				}
			}
		}
		return min
	}

	attrCanonical := map[string]string{}

	extractSection := func(s, kw string) string {
		ls := strings.ToUpper(s)
		kw = " " + strings.ToUpper(kw) + " "
		i := strings.Index(ls, kw)
		if i < 0 {
			kwAlt := " " + strings.ToUpper(kw[:len(kw)-1])
			i = strings.Index(ls, kwAlt)
			if i < 0 {
				return ""
			}
			i += len(kwAlt)
		} else {
			i += len(kw)
		}
		j := nextKW(ls, i)
		if j < 0 {
			j = len(s)
		}
		return strings.TrimSpace(s[i:j])
	}

	parseList := func(sec string) []string {
		if sec == "" {
			return nil
		}
		sec = strings.TrimSpace(sec)
		if strings.Contains(sec, "'") {
			ms := reQuoted.FindAllStringSubmatch(sec, -1)
			out := make([]string, 0, len(ms))
			for _, m := range ms {
				if v := strings.TrimSpace(m[1]); v != "" {
					low := strings.ToLower(v)
					out = append(out, low)
					if _, ok := attrCanonical[low]; !ok {
						attrCanonical[low] = v
					}
				}
			}
			return out
		}
		sec = strings.Trim(sec, "() ")
		parts := strings.FieldsFunc(sec, func(r rune) bool { return r == '$' || r == ' ' || r == '\n' || r == '\t' })
		out := make([]string, 0, len(parts))
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p == "" {
				continue
			}
			low := strings.ToLower(p)
			out = append(out, low)
			if _, ok := attrCanonical[low]; !ok {
				attrCanonical[low] = p
			}
		}
		return out
	}

	for _, raw := range defs {
		text := " " + strings.ReplaceAll(strings.ReplaceAll(raw, "\n", " "), "\r", " ") + " "
		names := parseList(extractSection(text, "NAME"))
		if len(names) == 0 {
			continue
		}
		must := parseList(extractSection(text, "MUST"))
		may := parseList(extractSection(text, "MAY"))
		sup := parseList(extractSection(text, "SUP"))

		o := oc{must: must, may: may, sup: sup}
		for _, n := range names {
			ocMap[strings.ToLower(n)] = o
		}
	}

	seen := map[string]bool{}
	acc := map[string]struct{}{}

	var visit func(string)
	visit = func(n string) {
		n = strings.ToLower(n)
		if seen[n] {
			return
		}
		seen[n] = true
		o, ok := ocMap[n]
		if !ok {
			return
		}
		for _, a := range o.must {
			acc[a] = struct{}{}
		}
		for _, a := range o.may {
			acc[a] = struct{}{}
		}
		for _, s := range o.sup {
			visit(s)
		}
	}
	for _, ocn := range objClasses {
		visit(ocn)
	}

	out := make([]string, 0, len(acc))
	for a := range acc {
		if canon, ok := attrCanonical[a]; ok {
			out = append(out, canon)
		} else {
			out = append(out, a)
		}
	}
	sort.Strings(out)
	return out, nil
}

func ExportLdap(url string, port int64, ssl bool, bindDN, bindPass, dn string) (string, error) {
	l, err := Connect(url, port, ssl)
	if err != nil {
		return "", err
	}
	defer l.Unbind()

	if bindDN != "" && bindPass != "" {
		if err = l.Bind(bindDN, bindPass); err != nil {
			return "", errors.New("failed to bind with provided credentials")
		}
	}

	searchReq := ldap.NewSearchRequest(
		dn,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(objectClass=*)",
		nil,
		nil,
	)

	result, err := l.Search(searchReq)
	if err != nil {
		return "", errors.New("failed to search for entries")
	}

	if len(result.Entries) == 0 {
		return "", nil
	}

	var builder strings.Builder
	for _, entry := range result.Entries {
		builder.WriteString("dn: " + entry.DN + "\n")
		for _, attr := range entry.Attributes {
			for _, value := range attr.Values {
				shouldEncode := binaryAttrs[attr.Name] || utils.ContainsNonASCII(value)

				if shouldEncode {
					encoded := base64.StdEncoding.EncodeToString([]byte(value))
					builder.WriteString(attr.Name + ":: " + encoded + "\n")
				} else {
					builder.WriteString(attr.Name + ": " + value + "\n")
				}
			}
		}
		builder.WriteString("\n")
	}

	return builder.String(), nil
}
