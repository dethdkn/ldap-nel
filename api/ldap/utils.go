package ldap

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

var binaryAttrs = map[string]bool{
	"jpegPhoto": true,
}

func unfold(b []byte) []string {
	sc := bufio.NewScanner(bytes.NewReader(b))
	var out []string
	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		if len(out) > 0 && strings.HasPrefix(line, " ") {
			out[len(out)-1] += strings.TrimPrefix(line, " ")
		} else {
			out = append(out, line)
		}
	}
	return out
}

func splitBlocks(lines []string) [][]string {
	var blocks [][]string
	var cur []string
	for _, ln := range lines {
		if strings.TrimSpace(ln) == "" {
			if len(cur) > 0 {
				blocks = append(blocks, cur)
				cur = nil
			}
			continue
		}
		cur = append(cur, ln)
	}
	if len(cur) > 0 {
		blocks = append(blocks, cur)
	}
	return blocks
}

func parseAddBlock(lines []string) (dn string, attrs map[string][]string, err error) {
	attrs = make(map[string][]string)
	for _, ln := range lines {
		k, v, isB64, e := parseKV(ln)
		if e != nil {
			return "", nil, e
		}
		if strings.EqualFold(k, "changetype") {
			return "", nil, errors.New("changetype not supported in this importer")
		}

		if strings.EqualFold(k, "dn") {
			if isB64 {
				raw, err := base64.StdEncoding.DecodeString(v)
				if err != nil {
					return "", nil, fmt.Errorf("invalid base64 DN: %w", err)
				}
				dn = string(raw)
			} else {
				dn = v
			}
			continue
		}

		if isBinaryAttr(k) {
			raw, err := base64.StdEncoding.DecodeString(v)
			if err != nil {
				return "", nil, fmt.Errorf("%s must be base64: %w", k, err)
			}
			attrs[k] = append(attrs[k], string(raw))
			continue
		}

		if isB64 {
			raw, err := base64.StdEncoding.DecodeString(v)
			if err != nil {
				return "", nil, fmt.Errorf("invalid base64 for %s: %w", k, err)
			}
			attrs[k] = append(attrs[k], string(raw))
		} else {
			attrs[k] = append(attrs[k], v)
		}
	}
	return dn, attrs, nil
}

func parseKV(line string) (key, val string, isB64 bool, err error) {
	i := strings.IndexByte(line, ':')
	if i <= 0 {
		return "", "", false, fmt.Errorf("invalid line: %q", line)
	}
	key = strings.TrimSpace(line[:i])

	if i+1 >= len(line) {
		return key, "", false, nil
	}

	switch line[i+1] {
	case ':':
		isB64 = true
		val = strings.TrimSpace(line[i+2:])
		return
	case '<':
		val = strings.TrimSpace(line[i+2:])
		return
	default:
		val = strings.TrimSpace(line[i+1:])
		return
	}
}

func isBinaryAttr(attr string) bool {
	for k := range binaryAttrs {
		if strings.EqualFold(k, attr) {
			return true
		}
	}
	return false
}
