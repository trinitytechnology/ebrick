package utils

import "strings"

func IsBlank(s *string) bool {
	if s == nil {
		return true
	}
	return strings.TrimSpace(*s) == ""
}

func Default(s *string, defaultValue string) string {
	if IsBlank(s) {
		return defaultValue
	}
	return *s
}
