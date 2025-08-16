package utils

import (
	"regexp"
	"strings"
)

func SanitizeUsername(input string) string {
	username := strings.TrimSpace(input)

	username = strings.ToLower(username)

	re := regexp.MustCompile(`[^a-z0-9._]`)
	username = re.ReplaceAllString(username, "")

	return username
}

func ValidatePassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	lower := regexp.MustCompile(`[a-z]`)
	upper := regexp.MustCompile(`[A-Z]`)
	digit := regexp.MustCompile(`[0-9]`)
	symbol := regexp.MustCompile(`[^a-zA-Z0-9]`)

	return lower.MatchString(password) && upper.MatchString(password) && digit.MatchString(password) && symbol.MatchString(password)
}
