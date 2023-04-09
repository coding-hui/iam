package utils

import "strings"

// Sanitize the inputs by removing line endings
func Sanitize(s string) string {
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\r", "")

	return s
}
